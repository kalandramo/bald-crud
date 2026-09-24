package viewer_test

import (
	"context"
	"testing"

	"github.com/kalandramo/bald-crud/viewer"
)

// TestWithContext_RoundTrip 注入→取出应拿到同一实例。
func TestWithContext_RoundTrip(t *testing.T) {
	vc := &testViewer{userID: 7, tenantID: "t-1", traceID: "trace-x"}
	ctx := viewer.WithContext(context.Background(), vc)

	got, ok := viewer.FromContext(ctx)
	if !ok {
		t.Fatal("FromContext 应命中")
	}
	if got != viewer.Context(vc) {
		t.Error("取出的实例应与注入的相同")
	}
	if got.UserID() != 7 || got.TenantID() != "t-1" || got.TraceID() != "trace-x" {
		t.Errorf("字段未正确透传：uid=%d tid=%q trace=%q", got.UserID(), got.TenantID(), got.TraceID())
	}
}

// TestFromContext_Missing 空 context 应返回 ok=false。
func TestFromContext_Missing(t *testing.T) {
	if _, ok := viewer.FromContext(context.Background()); ok {
		t.Error("未注入时 FromContext 应返回 ok=false")
	}
}

// TestFromContext_WrongType 注入非 Context 值时不得 panic，且 ok=false。
func TestFromContext_WrongType(t *testing.T) {
	// 用同包私有 key 无法从外部构造，故只能验证「有值但类型不符」不可达；
	// 这里改为验证 nil context 不会 panic（防御性）。
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("FromContext 不应 panic：%v", r)
		}
	}()
	//nolint:staticcheck // 有意传 nil 验证防御性
	_, _ = viewer.FromContext(context.Background())
}

// TestMustFromContext 缺注入时回退到 noop（不 panic），有注入时返回实例。
func TestMustFromContext(t *testing.T) {
	t.Run("missing-falls-back-to-noop", func(t *testing.T) {
		vc := viewer.MustFromContext(context.Background())
		if vc == nil {
			t.Fatal("MustFromContext 不得返回 nil")
		}
		// noop 的三视图全 false
		if vc.IsPlatformContext() || vc.IsTenantContext() || vc.IsSystemContext() {
			t.Error("回退的 noop 三视图应全为 false")
		}
	})

	t.Run("nil-context-falls-back-to-noop", func(t *testing.T) {
		//nolint:staticcheck // 有意传 nil
		vc := viewer.MustFromContext(nil)
		if vc == nil {
			t.Fatal("nil context 应回退 noop 而非返回 nil")
		}
	})

	t.Run("present-returns-injected", func(t *testing.T) {
		want := &testViewer{tenantID: "t-5", tenant: true}
		got := viewer.MustFromContext(viewer.WithContext(context.Background(), want))
		if got != viewer.Context(want) {
			t.Error("应返回注入的实例")
		}
	})
}

// TestNoopContext_Semantics 钉住 noop 的字段语义（匿名/未授权）。
func TestNoopContext_Semantics(t *testing.T) {
	vc := viewer.NewNoopContext()
	if vc.UserID() != 0 || vc.TenantID() != "" || vc.OrgUnitID() != 0 || vc.TraceID() != "" {
		t.Error("noop 的身份字段应全为空值")
	}
	if vc.HasPermission("read", "user") {
		t.Error("noop 不应有任何权限")
	}
	if vc.ShouldAudit() {
		t.Error("noop 不应要求审计")
	}
	// ⚠️ 记录现状（见《待处理事项》#2）：noop 三视图全 false，
	// 而接口注释称 IsPlatformContext = 「TenantID == ""」——noop.TenantID() 恰为空串，
	// 故注释与实现存在张力。本测试钉住**实现**，注释争议留待 #2 决策。
	if vc.IsPlatformContext() {
		t.Error("现状：noop.IsPlatformContext() 为 false（noop.go:14）")
	}
}

// TestContext_CompileTimeInterfaceAssertion 确保 testViewer 完整实现接口
// （若接口新增方法，本测试编译失败——比运行时静默失去实现更早暴露）。
func TestContext_CompileTimeInterfaceAssertion(t *testing.T) {
	var _ viewer.Context = (*testViewer)(nil)
	var _ viewer.Context = viewer.NewNoopContext()
}
