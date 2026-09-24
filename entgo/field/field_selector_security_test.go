package field

import (
	"strings"
	"testing"

	"entgo.io/ent/dialect/sql"
)

func buildSelectSQL(t *testing.T, table string, fields []string) string {
	t.Helper()
	s := sql.Select("*").From(sql.Table(table))
	NewFieldSelector().BuildSelect(s, fields)
	q, _ := s.Query()
	return q
}

// TestFieldSelector_HostilePathDropped 验证 FieldMask 路径的每一段（此前仅校验
// 第一个点之后的部分，前缀可携带注入载荷）都过硬性标识符校验，敌意路径
// 不会进入 SELECT 列表。
func TestFieldSelector_HostilePathDropped(t *testing.T) {
	hostile := []string{
		"(select version()) -- \xff.name", // H-2 原始攻击链：非法 UTF-8 前缀 + 合法后缀列
		"id) OR (1=1 --",
		"a`.name",
		"name AS (select version())",
		"user name",
		"id; DROP TABLE users",
	}
	for _, path := range hostile {
		for _, table := range []string{"users", "custom_table"} {
			q := buildSelectSQL(t, table, []string{path})
			// 注意：ent 正常输出会用反引号引用表/列名（`users`），
			// 这里只检查注入载荷相关字符与子查询标记。
			for _, d := range []string{"(", ")", "'", ";", "--", " AS ", "version()", "1=1", "DROP TABLE"} {
				if strings.Contains(q, d) {
					t.Errorf("hostile mask path %q leaked %q into SQL on %s: %q", path, d, table, q)
				}
			}
		}
	}
}

// TestFieldSelector_ValidPathStillWorks 验证合法路径不受影响。
// 注：entgo 的 NormalizePaths 会把 "user.name" 拍平为 "user_name"（既有行为），
// 此处仅断言合法输入仍产生正常 SELECT 且无注入残留。
func TestFieldSelector_ValidPathStillWorks(t *testing.T) {
	q := buildSelectSQL(t, "users", []string{"name"})
	if !strings.Contains(q, "name") {
		t.Fatalf("expected column name in SELECT, got %q", q)
	}
	// 带点路径：每段均为标识符 → 通过校验（后缀列校验 + 前缀硬性校验），
	// 最终按既有行为拍平为单列名
	q2 := buildSelectSQL(t, "users", []string{"user.name"})
	if strings.ContainsAny(q2, "();'") || strings.Contains(q2, " AS ") {
		t.Fatalf("valid dotted path must not introduce metacharacters, got %q", q2)
	}
	if !strings.Contains(q2, "user") || !strings.Contains(q2, "name") {
		t.Fatalf("expected flattened column from dotted path, got %q", q2)
	}
}

// TestFieldSelector_UnknownColumnDropped 验证已注册表上「合法标识符但不在该表
// 白名单」的列被 columnAllowed **拒绝**。
//
// 与 HostilePathDropped 的区别：那些用例的载荷含 SQL 元字符，被
// IsValidFieldPath 硬性校验拦下；本用例的 "password" 是**完全合法的标识符**，
// 只有列白名单能拦住它。此路径此前无测试覆盖（登记 #3g）——探针证实：
// 把 ent.CheckColumn 改为恒返回 nil 时，本包既有测试全部报 ok（未抓到）。
func TestFieldSelector_UnknownColumnDropped(t *testing.T) {
	// users 是已注册表，其白名单列为 id / tenant_id / name / age。
	for _, col := range []string{"password", "secret", "nonexistent_col"} {
		q := buildSelectSQL(t, "users", []string{col})
		if strings.Contains(q, col) {
			t.Errorf("column %q is not in users whitelist, must be dropped, got %q", col, q)
		}
	}
}

// TestFieldSelector_MixedValidAndUnknownColumns 验证混合输入下只丢弃未知列，
// 合法列仍被保留——拒绝逻辑不能误伤同批次的其他字段。
func TestFieldSelector_MixedValidAndUnknownColumns(t *testing.T) {
	q := buildSelectSQL(t, "users", []string{"name", "password"})
	if !strings.Contains(q, "name") {
		t.Errorf("valid column 'name' must survive, got %q", q)
	}
	if strings.Contains(q, "password") {
		t.Errorf("unknown column 'password' must be dropped, got %q", q)
	}
}

// TestFieldSelector_UnknownTableFailOpen 验证未注册表上列白名单 fail-open
// （保持无白名单时的旧行为）——合法标识符列仍被保留。
//
// 这是与上一条**相反**的方向：契约规定未注册表放行、已注册表的未知列拒绝。
// 两个方向都断言，才能钉住 columnAllowed 的分支语义。
func TestFieldSelector_UnknownTableFailOpen(t *testing.T) {
	q := buildSelectSQL(t, "custom_table", []string{"whatever"})
	if !strings.Contains(q, "whatever") {
		t.Errorf("unknown table must fail-open (keep column), got %q", q)
	}
}
