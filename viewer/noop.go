package viewer

// noopContext 实现 Context 接口，用于表示匿名或未授权用户。
//
// 语义（2026-09-25，见《待处理事项》#2）：noop 三视图**全 false**——它既非
// 平台视图也非系统视图，且租户为空。故经 EnforceTenant 判定为「身份不完整」
// → fail-closed（ErrMissingViewer），而非落入租户分支注入 `tenant_id = ''`。
// 这是「平台身份必须显式声明」纪律的体现：匿名绝不被推断为平台视图。
type noopContext struct{}

func (noopContext) UserID() uint64                 { return 0 }
func (noopContext) TenantID() string               { return "" }
func (noopContext) OrgUnitID() uint64              { return 0 }
func (noopContext) Permissions() []string          { return nil }
func (noopContext) Roles() []string                { return nil }
func (noopContext) DataScope() []DataScope         { return nil }
func (noopContext) TraceID() string                { return "" }
func (noopContext) HasPermission(_, _ string) bool { return false }
func (noopContext) IsPlatformContext() bool        { return false }
func (noopContext) IsTenantContext() bool          { return false }
func (noopContext) IsSystemContext() bool          { return false }
func (noopContext) ShouldAudit() bool              { return false }

// NewNoopContext 创建一个匿名上下文实例
func NewNoopContext() Context {
	return noopContext{}
}
