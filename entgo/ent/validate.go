package ent

import (
	"regexp"
	"strings"
)

// validFieldNameRe 单段 SQL 标识符硬性白名单：字母/下划线开头，仅含字母、数字、下划线。
// 过滤字段名、排序字段名、FieldMask 路径在进入任何 SQL 构建器之前必须通过该校验。
// 注意：不能依赖 stringcase.ToSnakeCase 之类的归一化做清洗——它对非法 UTF-8 输入
// 原样透传，SQL 元字符（引号/括号/注释符）可存活；也不能依赖列白名单——列白名单
// 对未注册表 fail-open。此校验是独立于二者的最后一道防线。
var validFieldNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// IsValidFieldName 报告 name 是否为安全的单段 SQL 标识符。
// 不做 TrimSpace 等宽容处理：调用方应原样拒绝而非修剪后使用。
func IsValidFieldName(name string) bool {
	return validFieldNameRe.MatchString(name)
}

// IsValidFieldPath 报告点分隔路径的每一段是否均为安全标识符（如 FieldMask 的 "user.name"）。
func IsValidFieldPath(path string) bool {
	if path == "" {
		return false
	}
	for _, p := range strings.Split(path, ".") {
		if !IsValidFieldName(p) {
			return false
		}
	}
	return true
}

// CheckColumn 报告 column 是否存在于已注册白名单的表 table 中。
//
// 这是 ent.go（生成代码）中私有 checkColumn 的**导出别名**——手写的
// field / filter / sorting 三包在进入 SQL 构建器之前用它做列白名单校验。
// 委托而非复制：白名单注册表（menu / user 的 Table + ValidColumn）只有一份
// 真相源，生成代码与手写代码不会各自判断而漂移。
//
// 返回值语义是**契约**，调用方依赖它区分两种拒绝：
//   - 表未注册 → 错误含 "unknown table" → 调用方 **fail-open**（保留字段，
//     保持无白名单时的旧行为，避免误伤非 ent 管理的表）
//   - 表已注册但列不存在 → 错误含 "unknown column" → 调用方 **拒绝**该字段
//
// 该文本由 entgo.io/ent 的 sql.NewColumnCheck 产生，不由本包控制——
// 故有 TestCheckColumn_UnknownTable 钉住它。
func CheckColumn(table, column string) error {
	return checkColumn(table, column)
}
