package viewer

import "context"

// Context 定义当前访问者（Viewer）上下文接口，适用于查询和更新等操作。
// 增加了权限相关方法，便于在 Hook/Policy 中判断是否允许更新/删除等操作。
type Context interface {
	// UserID 返回当前用户ID
	UserID() uint64

	// TenantID 返回租户ID。
	//
	// 类型为 string（2026-09-24 统一）：空串表示「无租户上下文」（身份不完整）。
	// 此前为 uint64，用 0 同时表达「平台视图」与「非数字租户 ID 解析失败」，
	// 导致后者被静默误判为平台视图而跳过租户强制（安全失效）。string 无解析
	// 失败路径，语义单一。
	//
	// 注意（2026-09-25，见《待处理事项》#2）：空租户**不等于**平台视图——
	// 平台身份必须显式声明（见 IsPlatformContext）。空租户经 EnforceTenant
	// 一律 fail-closed。
	TenantID() string

	// OrgUnitID 返回当前身份挂载的组织单元 ID
	OrgUnitID() uint64

	// Permissions 返回当前 Viewer 的权限列表（可用于细粒度判断）
	Permissions() []string

	// Roles 返回当前 Viewer 的角色列表（可选，用于审计或策略）
	Roles() []string

	// DataScope 返回当前身份的数据权限范围（用于 SQL 拼接）
	DataScope() []DataScope

	// TraceID 返回当前请求的 Trace ID（用于日志跟踪）
	TraceID() string

	// HasPermission 判断是否具有某个动作/资源的权限（如 "update:user"）
	HasPermission(action, resource string) bool

	// IsPlatformContext 当前是否处于平台管理视图（跨租户）。
	//
	// 语义（2026-09-25 收紧，见《待处理事项》#2）：**必须显式声明**，绝不基于
	// 「TenantID 为空」推断——空租户同时表示「匿名/未认证」与「平台视图」两种
	// 相反语义，靠推断会 fail-open（匿名请求被当作平台视图而看到全部租户数据）。
	// 实现者须以显式字段（如 SimpleViewer.Platform）表达平台身份。
	IsPlatformContext() bool

	// IsTenantContext 当前是否处于租户业务视图（有租户身份的普通业务请求）。
	IsTenantContext() bool

	// IsSystemContext 判断是否为系统后台任务
	IsSystemContext() bool

	// ShouldAudit 返回是否需要记录审计日志（便于在中间件/Hook 中快速判断）
	ShouldAudit() bool
}

type contextKey struct{}

// WithContext 将 Context 注入 context
func WithContext(ctx context.Context, vc Context) context.Context {
	return context.WithValue(ctx, contextKey{}, vc)
}

// FromContext 从 context 中提取 Context
func FromContext(ctx context.Context) (Context, bool) {
	v := ctx.Value(contextKey{})
	vc, ok := v.(Context)
	return vc, ok
}

// MustFromContext 从 context 中提取 Context，若不存在则返回一个默认的 NoopContext
func MustFromContext(ctx context.Context) Context {
	if ctx == nil {
		return NewNoopContext()
	}
	if v := ctx.Value(contextKey{}); v != nil {
		if vc, ok := v.(Context); ok && vc != nil {
			return vc
		}
	}
	return NewNoopContext()
}
