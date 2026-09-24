package viewer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kalandramo/bald-crud/viewer"
)

// testViewer 是 viewer.Context 的可配置实现：三个 Is* 判定**独立可设**，
// 便于精确覆盖 EnforceTenant 的分支，而非依赖 TenantID 反推（后者会把
// 「平台视图」与「解析失败」混为一谈——正是 2026-09-24 修掉的那类缺陷）。
type testViewer struct {
	userID   uint64
	tenantID string
	orgUnit  uint64
	perms    []string
	roles    []string
	scopes   []viewer.DataScope
	traceID  string

	platform bool
	tenant   bool
	system   bool
	audit    bool
}

var _ viewer.Context = (*testViewer)(nil)

func (v *testViewer) UserID() uint64                { return v.userID }
func (v *testViewer) TenantID() string              { return v.tenantID }
func (v *testViewer) OrgUnitID() uint64             { return v.orgUnit }
func (v *testViewer) Permissions() []string         { return v.perms }
func (v *testViewer) Roles() []string               { return v.roles }
func (v *testViewer) DataScope() []viewer.DataScope { return v.scopes }
func (v *testViewer) TraceID() string               { return v.traceID }

func (v *testViewer) HasPermission(action, resource string) bool {
	want := action + ":" + resource
	for _, p := range v.perms {
		if p == want {
			return true
		}
	}
	return false
}

func (v *testViewer) IsPlatformContext() bool { return v.platform }
func (v *testViewer) IsTenantContext() bool   { return v.tenant }
func (v *testViewer) IsSystemContext() bool   { return v.system }
func (v *testViewer) ShouldAudit() bool       { return v.audit }

// ctxWith 把 viewer 注入 context。
func ctxWith(vc viewer.Context) context.Context {
	return viewer.WithContext(context.Background(), vc)
}

// ---------------------------------------------------------------------------
// EnforceTenant：三态决策（全链路租户隔离的唯一闸门）
// ---------------------------------------------------------------------------

// TestEnforceTenant_MissingViewerFailsClosed 缺身份必须**拒绝**而非静默放行
// ——否则租户过滤会悄悄消失（与 entgo TenantPrivacy 的 fail-closed 一致）。
func TestEnforceTenant_MissingViewerFailsClosed(t *testing.T) {
	dec, err := viewer.EnforceTenant(context.Background())
	if !errors.Is(err, viewer.ErrMissingViewer) {
		t.Fatalf("缺 viewer 应返回 ErrMissingViewer，got %v", err)
	}
	if dec.Enforce {
		t.Error("缺身份不得 Enforce（fail-closed）")
	}
}

// TestEnforceTenant_TenantViewEnforces 租户业务视图必须注入租户谓词。
func TestEnforceTenant_TenantViewEnforces(t *testing.T) {
	vc := &testViewer{tenantID: "t-42", tenant: true}
	dec, err := viewer.EnforceTenant(ctxWith(vc))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !dec.Enforce {
		t.Fatal("租户业务视图必须 Enforce=true")
	}
	if dec.TenantID != "t-42" {
		t.Errorf("TenantID=%q，want t-42", dec.TenantID)
	}
}

// TestEnforceTenant_PlatformAndSystemPassThrough 平台视图与系统任务均放行
// （pass-through：不注入租户谓词）。
func TestEnforceTenant_PlatformAndSystemPassThrough(t *testing.T) {
	cases := map[string]*testViewer{
		"platform": {platform: true},
		"system":   {system: true, tenantID: "t-1"},
	}
	for name, vc := range cases {
		t.Run(name, func(t *testing.T) {
			dec, err := viewer.EnforceTenant(ctxWith(vc))
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if dec.Enforce {
				t.Error("平台/系统视图不得 Enforce")
			}
		})
	}
}

// TestEnforceTenant_PlatformPrecedence 平台与租户标志同时为真时走平台分支
// ——决策由 IsPlatformContext 优先，不被其他标志干扰。
func TestEnforceTenant_PlatformPrecedence(t *testing.T) {
	vc := &testViewer{platform: true, tenant: true, tenantID: "t-9"}
	dec, err := viewer.EnforceTenant(ctxWith(vc))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dec.Enforce {
		t.Error("平台视图应优先，不得 Enforce")
	}
}

// TestEnforceTenant_NoopContextCurrentBehavior 钉住 noop（匿名）上下文经
// EnforceTenant 的**当前**决策。
//
// ⚠️ 这是《待处理事项》#2 记录的语义争议点，本测试记录**现状**、不背书其正确性：
// noopContext.IsPlatformContext() 返回 false（noop.go:14），故它不被认作平台视图，
// 落入租户业务视图分支 → Enforce=true 但 TenantID==""。
// 下游据此会注入 `tenant_id = ''` 谓词，与「匿名请求不应看到任何行」的
// fail-closed 预期存在张力。
//
// 若 #2 决策变更（如 noop 改为显式平台视图、或 fail-closed 拒绝），
// 本测试将变红——那正是提醒同步更新《待处理事项》#2 的信号。
func TestEnforceTenant_NoopContextCurrentBehavior(t *testing.T) {
	dec, err := viewer.EnforceTenant(ctxWith(viewer.NewNoopContext()))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !dec.Enforce {
		t.Fatalf("现状：noop 落入租户分支应 Enforce=true。" +
			"若此处失败，说明 #2 已决策变更——请同步更新本测试与《待处理事项》#2")
	}
	if dec.TenantID != "" {
		t.Errorf("现状：noop 的 TenantID 为空串，got %q", dec.TenantID)
	}
}

// ---------------------------------------------------------------------------
// ScopedModel 类型门控
// ---------------------------------------------------------------------------

// scopedEntity 实现 viewer.ScopedModel（指针接收者）——模拟嵌入 TenantID
// mixin 的租户隔离实体。
type scopedEntity struct{ tenantID *string }

func (e *scopedEntity) GetTenantID() *string  { return e.tenantID }
func (e *scopedEntity) SetTenantID(id string) { e.tenantID = &id }

// unscopedEntity 未实现 ScopedModel——平台级实体（表无 tenant_id 列）。
type unscopedEntity struct{ name string }

func TestIsTenantScopedType(t *testing.T) {
	if !viewer.IsTenantScopedType[scopedEntity]() {
		t.Error("scopedEntity 应被识别为 tenant-scoped")
	}
	if viewer.IsTenantScopedType[unscopedEntity]() {
		t.Error("unscopedEntity 不应被识别为 tenant-scoped")
	}
}

// TestIsTenantScopedType_PointerTypeArg 记录一个易误用点：本函数检查的是
// `*T`，故传指针类型 T=*scopedEntity 时实际检查 **scopedEntity → 不实现 → false。
// 调用方须传**值类型**。
func TestIsTenantScopedType_PointerTypeArg(t *testing.T) {
	if viewer.IsTenantScopedType[*scopedEntity]() {
		t.Error("传指针类型时应为 false（内部检查的是 *T = **scopedEntity）——调用方须传值类型")
	}
}

// ---------------------------------------------------------------------------
// EnforceOnScopedInstance / EnforceOnScopedInstanceAny（Create 路径实例级强制）
// ---------------------------------------------------------------------------

func TestEnforceOnScopedInstance_TenantViewSetsTenant(t *testing.T) {
	e := &scopedEntity{}
	err := viewer.EnforceOnScopedInstance(ctxWith(&testViewer{tenant: true, tenantID: "t-7"}), e)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if e.tenantID == nil || *e.tenantID != "t-7" {
		t.Errorf("tenantID 应被置为 t-7，got %v", e.tenantID)
	}
}

func TestEnforceOnScopedInstance_PlatformViewLeavesUntouched(t *testing.T) {
	e := &scopedEntity{}
	if err := viewer.EnforceOnScopedInstance(ctxWith(&testViewer{platform: true}), e); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if e.tenantID != nil {
		t.Errorf("平台视图不应写 tenantID，got %q", *e.tenantID)
	}
}

func TestEnforceOnScopedInstance_NonScopedIsNoop(t *testing.T) {
	if err := viewer.EnforceOnScopedInstance(ctxWith(&testViewer{tenant: true, tenantID: "t-1"}), &unscopedEntity{}); err != nil {
		t.Errorf("非 scoped 实体应直接返回 nil，got %v", err)
	}
}

func TestEnforceOnScopedInstance_NilInstance(t *testing.T) {
	var e *scopedEntity
	if err := viewer.EnforceOnScopedInstance(ctxWith(&testViewer{tenant: true, tenantID: "t-1"}), e); err != nil {
		t.Errorf("nil 实例应返回 nil，got %v", err)
	}
}

func TestEnforceOnScopedInstance_MissingViewerFailsClosed(t *testing.T) {
	err := viewer.EnforceOnScopedInstance(context.Background(), &scopedEntity{})
	if !errors.Is(err, viewer.ErrMissingViewer) {
		t.Fatalf("scoped 实体缺身份必须 fail-closed，got %v", err)
	}
}

func TestEnforceOnScopedInstanceAny_TenantViewSetsTenant(t *testing.T) {
	e := &scopedEntity{}
	if err := viewer.EnforceOnScopedInstanceAny(ctxWith(&testViewer{tenant: true, tenantID: "t-3"}), e); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if e.tenantID == nil || *e.tenantID != "t-3" {
		t.Errorf("tenantID 应被置为 t-3，got %v", e.tenantID)
	}
}

func TestEnforceOnScopedInstanceAny_NilIsNoop(t *testing.T) {
	if err := viewer.EnforceOnScopedInstanceAny(context.Background(), nil); err != nil {
		t.Errorf("nil 应返回 nil，got %v", err)
	}
}

func TestEnforceOnScopedInstanceAny_NonScopedIsNoop(t *testing.T) {
	if err := viewer.EnforceOnScopedInstanceAny(context.Background(), &unscopedEntity{}); err != nil {
		t.Errorf("非 scoped 应返回 nil，got %v", err)
	}
}

// TestEnforceOnScopedInstanceAny_TypedNilPointer 探 typed-nil 装箱的边界：
// 泛型版（instance *T）能识别 nil 指针，但 any 版收到的 `(*scopedEntity)(nil)`
// 装箱后 `instance == nil` 为 false（带类型信息），故会继续走到 SetTenantID。
// 本测试记录实际行为，不 panic 即视为安全。
func TestEnforceOnScopedInstanceAny_TypedNilPointer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("any 版对 typed-nil 指针 panic 了（应返回 error 而非崩溃）：%v", r)
		}
	}()
	var e *scopedEntity
	err := viewer.EnforceOnScopedInstanceAny(ctxWith(&testViewer{tenant: true, tenantID: "t-1"}), e)
	t.Logf("typed-nil 指针：err=%v（无 panic）", err)
}
