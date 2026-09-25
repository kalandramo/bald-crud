package gorm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	gormcrud "github.com/kalandramo/bald-crud/gorm"
	gormmixin "github.com/kalandramo/bald-crud/gorm/mixin"
	"github.com/kalandramo/bald-crud/viewer"
)

// tenantTestEntity 嵌入 TenantID mixin，触发租户隔离强制。
type tenantTestEntity struct {
	gorm.Model
	gormmixin.TenantID
	Name string
}

func openTenantTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&tenantTestEntity{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	// 必须注册租户回调（NewClient 路径会自动注册，但这里直接用 gorm.Open
	// 绕过了 NewClient，故手动注册以测试回调逻辑本身）
	if err := gormcrud.RegisterTenantCallbacks(db); err != nil {
		t.Fatalf("register tenant callbacks: %v", err)
	}
	return db
}

// viewerCtx 构造一个指定租户的 viewer Context。
func viewerCtx(tid string) context.Context {
	return viewer.WithContext(context.Background(), &stubViewer{tid: tid})
}

// platformViewerCtx 构造一个平台视图的 viewer Context（2026-09-25：平台身份
// 必须显式声明，不再由空租户推断）。
func platformViewerCtx() context.Context {
	return viewer.WithContext(context.Background(), &stubViewer{platform: true})
}

type stubViewer struct {
	tid      string
	platform bool
}

func (s *stubViewer) UserID() uint64                 { return 0 }
func (s *stubViewer) TenantID() string               { return s.tid }
func (s *stubViewer) OrgUnitID() uint64              { return 0 }
func (s *stubViewer) Permissions() []string          { return nil }
func (s *stubViewer) Roles() []string                { return nil }
func (s *stubViewer) DataScope() []viewer.DataScope  { return nil }
func (s *stubViewer) TraceID() string                { return "" }
func (s *stubViewer) HasPermission(_, _ string) bool { return false }
func (s *stubViewer) IsPlatformContext() bool        { return s.platform }
func (s *stubViewer) IsTenantContext() bool          { return s.tid != "" && !s.platform }
func (s *stubViewer) IsSystemContext() bool          { return false }
func (s *stubViewer) ShouldAudit() bool              { return false }

// TestTenantEnforce_QueryInjectsPredicate 租户业务视图下，查询必须带上
// tenant_id = ? 谓词。
func TestTenantEnforce_QueryInjectsPredicate(t *testing.T) {
	db := openTenantTestDB(t)
	ctx := viewerCtx("t-7")

	var out tenantTestEntity
	tx := db.WithContext(ctx).Session(&gorm.Session{DryRun: true}).First(&out, 1)
	if tx.Error != nil {
		t.Fatalf("query: %v", tx.Error)
	}
	sql := tx.Statement.SQL.String()
	if !strings.Contains(sql, "tenant_id") {
		t.Errorf("tenant context query must inject tenant_id predicate, got SQL: %q", sql)
	}
	if !strings.Contains(sql, "?") {
		t.Errorf("tenant predicate must be parameterized, got SQL: %q", sql)
	}
}

// TestTenantEnforce_PlatformContextPassThrough 平台视图（tid==""）不注入谓词。
func TestTenantEnforce_PlatformContextPassThrough(t *testing.T) {
	db := openTenantTestDB(t)
	ctx := platformViewerCtx()

	var out tenantTestEntity
	tx := db.WithContext(ctx).Session(&gorm.Session{DryRun: true}).First(&out, 1)
	if tx.Error != nil {
		t.Fatalf("query: %v", tx.Error)
	}
	sql := tx.Statement.SQL.String()
	if strings.Contains(strings.ToLower(sql), "tenant_id") {
		t.Errorf("platform context must NOT inject tenant predicate, got SQL: %q", sql)
	}
}

// TestTenantEnforce_MissingViewerFailClosed 缺 ViewerContext 必须报错。
func TestTenantEnforce_MissingViewerFailClosed(t *testing.T) {
	db := openTenantTestDB(t)
	// 无 viewer context
	var out tenantTestEntity
	tx := db.WithContext(context.Background()).First(&out, 1)
	if tx.Error == nil {
		t.Errorf("missing viewer context must fail-closed, got nil error")
	}
}

// TestTenantEnforce_EmptyTenantFailsClosed 端到端锁定方案 D 的核心闸门
// （2026-09-25，见《待处理事项》#2）：非平台、非系统、租户为空的身份
// 经真实 gorm 查询必须 **fail-closed**，而非注入 `tenant_id = ''`。
//
// 旧语义下该身份被 IsPlatformContext 由空租户**推断**为平台视图 → 放行
// （fail-open）；新语义要求平台身份显式声明，空租户落入「身份不完整」被拒。
// 这是真实 DB + 真实 callback 的完整路径（仅身份为注入，非 mock 中间层）。
func TestTenantEnforce_EmptyTenantFailsClosed(t *testing.T) {
	db := openTenantTestDB(t)
	ctx := viewer.WithContext(context.Background(), &stubViewer{tid: ""}) // 空租户、非平台

	var out tenantTestEntity
	tx := db.WithContext(ctx).First(&out, 1)
	if tx.Error == nil {
		t.Fatalf("空租户（非显式平台、非系统）必须 fail-closed，got nil error")
	}
}

// TestTenantEnforce_CreateForcesTenantID 租户业务视图下 Create 强制覆盖
// tenant_id，即使实体设置了其他租户值也会被改写为当前 viewer 的租户。
func TestTenantEnforce_CreateForcesTenantID(t *testing.T) {
	db := openTenantTestDB(t)
	ctx := viewerCtx("t-7")

	bad := &tenantTestEntity{Name: "x"}
	other := "t-99"
	// 通过 mixin 字段显式赋值他租户，验证强制覆盖
	tenantField := &bad.TenantID // mixin.TenantID 嵌入实例
	tenantField.TenantID = &other
	if err := db.WithContext(ctx).Create(bad).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	// 重新读出验证 tenant_id 被强制为 7
	var got tenantTestEntity
	if err := db.WithContext(ctx).First(&got, bad.ID).Error; err != nil {
		t.Fatalf("re-read: %v", err)
	}
	stored := got.TenantID.TenantID
	if stored == nil || *stored != "t-7" {
		t.Errorf("tenant_id must be force-set to t-7, got %v", stored)
	}
}
