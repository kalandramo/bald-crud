package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/kalandramo/bald-crud/entgo/rule"
)

// TenantID 是租户隔离的 ent mixin。
//
// 类型统一为 string（2026-09-24）：此前为 `TenantID[IDT uint32 | uint64]` 的
// 伪泛型——`IDT` 仅传给 rule.TenantPrivacy，字段硬编码 `field.Uint32`，实际只
// 支持 uint32。现去掉未使用的类型参数，字段统一 `field.String`，与 bald 生态
// （contextx/authn/jwt/audit 全为 string）对齐，并消除「非数字租户 ID 解析失败
// 被静默当作平台视图」的安全隐患。
type TenantID struct{ mixin.Schema }

func (TenantID) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").
			Comment("租户ID").
			Immutable().
			Default("").
			Nillable().
			Optional(),
	}
}

func (TenantID) Policy() ent.Policy {
	return rule.TenantPrivacy{}
}
