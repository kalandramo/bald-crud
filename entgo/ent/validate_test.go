package ent

import (
	"strings"
	"testing"
)

// TestCheckColumn_RegisteredTableValidColumn 已注册白名单的表 + 合法列 → 通过。
func TestCheckColumn_RegisteredTableValidColumn(t *testing.T) {
	for _, c := range []struct{ table, col string }{
		{"users", "name"},
		{"users", "id"},
		{"users", "tenant_id"},
		{"menus", "id"},
	} {
		if err := CheckColumn(c.table, c.col); err != nil {
			t.Errorf("CheckColumn(%q, %q) must pass, got %v", c.table, c.col, err)
		}
	}
}

// TestCheckColumn_RegisteredTableUnknownColumn 已注册表 + 不存在的列 →
// 报 "unknown column"（调用方据此**拒绝**该列）。
func TestCheckColumn_RegisteredTableUnknownColumn(t *testing.T) {
	for _, c := range []struct{ table, col string }{
		{"users", "nope"},
		{"users", "password"},
		{"menus", "nope"},
	} {
		err := CheckColumn(c.table, c.col)
		if err == nil {
			t.Errorf("CheckColumn(%q, %q) must be rejected", c.table, c.col)
			continue
		}
		if !strings.Contains(err.Error(), "unknown column") {
			t.Errorf("CheckColumn(%q, %q) error must mention \"unknown column\", got %v", c.table, c.col, err)
		}
	}
}

// TestCheckColumn_UnknownTable 未注册表 → 报 "unknown table"。
//
// 这是**契约级断言**：三个调用方（field/filter/sorting）都以
// `strings.Contains(err.Error(), "unknown table")` 判定「表未注册 → 列白名单
// fail-open」。若错误文本漂移，fail-open 会静默变成 fail-closed，
// 未注册表的合法字段会被全部丢弃。
func TestCheckColumn_UnknownTable(t *testing.T) {
	for _, table := range []string{"custom_table", "orders", "not_registered"} {
		err := CheckColumn(table, "whatever")
		if err == nil {
			t.Errorf("CheckColumn(%q, _) must be rejected as unknown table", table)
			continue
		}
		if !strings.Contains(err.Error(), "unknown table") {
			t.Errorf("CheckColumn(%q, _) error must mention \"unknown table\" (callers depend on it for fail-open), got %v", table, err)
		}
	}
}

// TestCheckColumn_ConsistentWithGenerated 与生成代码的私有 checkColumn 判定一致。
//
// 防「同一不变量的两处各自判断漂移」——本项目最高频错误（已发生四次）。
// 导出版供手写的 field/filter/sorting 调用，私有版供生成代码（Asc/Desc 等）调用；
// 两者必须对同一输入给出相同的接受/拒绝结论，否则生成路径与手写路径的安全语义会分叉。
func TestCheckColumn_ConsistentWithGenerated(t *testing.T) {
	cases := []struct{ table, col string }{
		{"users", "name"},
		{"users", "nope"},
		{"users", ""},
		{"menus", "id"},
		{"menus", "nope"},
		{"custom_table", "x"},
		{"", ""},
	}
	for _, c := range cases {
		got := CheckColumn(c.table, c.col)
		want := checkColumn(c.table, c.col)
		if (got == nil) != (want == nil) {
			t.Errorf("drift: CheckColumn(%q, %q)=%v but generated checkColumn=%v",
				c.table, c.col, got, want)
		}
	}
}
