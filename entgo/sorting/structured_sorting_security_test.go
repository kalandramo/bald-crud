package sorting

import (
	"strings"
	"testing"

	"entgo.io/ent/dialect/sql"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
)

func buildSortingSQL(t *testing.T, table, field string) string {
	t.Helper()
	sel, err := NewStructuredSorting().BuildSelector([]*storev1.Sorting{
		{Field: field, Direction: storev1.Sorting_DESC},
	})
	if err != nil {
		t.Fatalf("BuildSelector error: %v", err)
	}
	s := sql.Select("*").From(sql.Table(table))
	sel(s)
	q, _ := s.Query()
	return q
}

// TestSorting_HostileFieldDropped 验证含 SQL 元字符的排序字段被硬性校验拒绝，
// 不会进入 ORDER BY（列白名单对未注册表 fail-open，故此校验必须独立生效）。
func TestSorting_HostileFieldDropped(t *testing.T) {
	for _, field := range []string{
		"id) OR (1=1 --",
		"(select version())",
		`name" --`,
		"id\xfff' OR '1'='1",
	} {
		for _, table := range []string{"users", "custom_table"} {
			q := buildSortingSQL(t, table, field)
			for _, d := range []string{"(", "'", `"`, "--", "version()", ";"} {
				if strings.Contains(q, d) {
					t.Errorf("hostile sort field %q leaked %q into SQL on %s: %q", field, d, table, q)
				}
			}
			if strings.Contains(q, "ORDER BY") {
				t.Errorf("hostile sort field %q must be dropped on %s, got ORDER BY: %q", field, table, q)
			}
		}
	}
}

// TestSorting_ValidFieldStillWorks 验证已注册表的合法列排序不受影响。
func TestSorting_ValidFieldStillWorks(t *testing.T) {
	q := buildSortingSQL(t, "users", "name")
	if !strings.Contains(q, "ORDER BY") {
		t.Fatalf("expected ORDER BY for valid field, got %q", q)
	}
	if !strings.Contains(q, "name") {
		t.Fatalf("expected column name in ORDER BY, got %q", q)
	}
}

// TestSorting_DefaultFieldValidated 验证默认排序字段同样过硬性校验。
func TestSorting_DefaultFieldValidated(t *testing.T) {
	sel, err := NewStructuredSorting().BuildSelectorWithDefaultField(nil, "(select version())", true)
	if err != nil {
		t.Fatalf("BuildSelectorWithDefaultField error: %v", err)
	}
	s := sql.Select("*").From(sql.Table("custom_table"))
	sel(s)
	q, _ := s.Query()
	if strings.Contains(q, "ORDER BY") || strings.Contains(q, "version()") {
		t.Errorf("hostile default sort field must be dropped, got %q", q)
	}
}

// TestSorting_UnknownColumnDropped 验证已注册表上「合法标识符但不在该表白名单」
// 的排序字段被 columnAllowed **拒绝**。
//
// 与 HostileFieldDropped 的区别：那些用例的载荷含 SQL 元字符，被
// IsValidFieldName 硬性校验拦下；本用例的 "password" 是**完全合法的标识符**，
// 只有列白名单能拦住它。此路径此前无测试覆盖（登记 #3g）——探针证实：
// 把 ent.CheckColumn 改为恒返回 nil 时，本包既有测试全部报 ok（未抓到）。
func TestSorting_UnknownColumnDropped(t *testing.T) {
	for _, col := range []string{"password", "secret", "nonexistent_col"} {
		q := buildSortingSQL(t, "users", col)
		if strings.Contains(q, "ORDER BY") {
			t.Errorf("column %q is not in users whitelist, must not appear in ORDER BY, got %q", col, q)
		}
	}
}

// TestSorting_MixedValidAndUnknownFields 验证混合输入下只丢弃未知字段，
// 合法字段仍被保留。
func TestSorting_MixedValidAndUnknownFields(t *testing.T) {
	sel, err := NewStructuredSorting().BuildSelector([]*storev1.Sorting{
		{Field: "password", Direction: storev1.Sorting_ASC},
		{Field: "name", Direction: storev1.Sorting_DESC},
	})
	if err != nil {
		t.Fatalf("BuildSelector error: %v", err)
	}
	s := sql.Select("*").From(sql.Table("users"))
	sel(s)
	q, _ := s.Query()
	if !strings.Contains(q, "ORDER BY") || !strings.Contains(q, "name") {
		t.Errorf("valid field 'name' must survive, got %q", q)
	}
	if strings.Contains(q, "password") {
		t.Errorf("unknown field 'password' must be dropped, got %q", q)
	}
}

// TestSorting_UnknownTableFailOpen 验证未注册表上列白名单 fail-open——
// 合法标识符字段仍被保留（契约的相反方向）。
func TestSorting_UnknownTableFailOpen(t *testing.T) {
	q := buildSortingSQL(t, "custom_table", "whatever")
	if !strings.Contains(q, "ORDER BY") || !strings.Contains(q, "whatever") {
		t.Errorf("unknown table must fail-open (keep field), got %q", q)
	}
}

// TestSorting_DefaultFieldUnknownColumnDropped 验证默认排序字段走同一条
// 白名单拒绝路径——BuildSelectorWithDefaultField 有独立代码分支，
// 若只测 BuildSelector 会漏掉这一处。
func TestSorting_DefaultFieldUnknownColumnDropped(t *testing.T) {
	sel, err := NewStructuredSorting().BuildSelectorWithDefaultField(nil, "password", true)
	if err != nil {
		t.Fatalf("BuildSelectorWithDefaultField error: %v", err)
	}
	s := sql.Select("*").From(sql.Table("users"))
	sel(s)
	q, _ := s.Query()
	if strings.Contains(q, "ORDER BY") || strings.Contains(q, "password") {
		t.Errorf("default sort field 'password' not in whitelist must be dropped, got %q", q)
	}
}
