package filter

import (
	"testing"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
	"github.com/kalandramo/bald-crud/mongodb/query"
)

func TestProcessor_Process_ReturnsBuilder_NoPanic(t *testing.T) {
	proc := NewProcessor()

	ops := []storev1.Operator{
		storev1.Operator_EQ,
		storev1.Operator_NEQ,
		storev1.Operator_IN,
		storev1.Operator_NIN,
		storev1.Operator_GTE,
		storev1.Operator_GT,
		storev1.Operator_LTE,
		storev1.Operator_LT,
		storev1.Operator_BETWEEN,
		storev1.Operator_IS_NULL,
		storev1.Operator_IS_NOT_NULL,
		storev1.Operator_CONTAINS,
		storev1.Operator_ICONTAINS,
		storev1.Operator_STARTS_WITH,
		storev1.Operator_ISTARTS_WITH,
		storev1.Operator_ENDS_WITH,
		storev1.Operator_IENDS_WITH,
		storev1.Operator_EXACT,
		storev1.Operator_IEXACT,
		storev1.Operator_REGEXP,
		storev1.Operator_IREGEXP,
		storev1.Operator_SEARCH,
	}

	for _, op := range ops {
		qb := &query.Builder{}
		got := proc.Process(qb, op, "name", "val", []any{"a", "b"})
		if got == nil {
			t.Fatalf("Process returned nil for op %v", op)
		}
		if got != qb {
			t.Fatalf("Process should return the same builder pointer for op %v", op)
		}
	}
}

func TestProcessor_SpecificCases_ReturnsBuilder(t *testing.T) {
	proc := NewProcessor()

	t.Run("Equal_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.Equal(qb, "name", "tom")
		if got == nil || got != qb {
			t.Fatalf("Equal should return the same non-nil builder")
		}
	})

	t.Run("In_JSONArray_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.In(qb, "name", `["a","b"]`, nil)
		if got == nil || got != qb {
			t.Fatalf("In should return the same non-nil builder")
		}
	})

	t.Run("NotIn_WithValues_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.NotIn(qb, "status", "", []any{"x", "y"})
		if got == nil || got != qb {
			t.Fatalf("NotIn should return the same non-nil builder")
		}
	})

	t.Run("Range_Between_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.Range(qb, "created_at", `["2020-01-01","2021-01-01"]`, nil)
		if got == nil || got != qb {
			t.Fatalf("Range should return the same non-nil builder")
		}
	})

	t.Run("IsNull_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.IsNull(qb, "deleted_at")
		if got == nil || got != qb {
			t.Fatalf("IsNull should return the same non-nil builder")
		}
	})

	t.Run("Contains_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.Contains(qb, "title", "go")
		if got == nil || got != qb {
			t.Fatalf("Contains should return the same non-nil builder")
		}
	})

	t.Run("JsonField_DotPath_ReturnsBuilder", func(t *testing.T) {
		qb := &query.Builder{}
		got := proc.Process(qb, storev1.Operator_EQ, "preferences.daily_email", "1", nil)
		if got == nil || got != qb {
			t.Fatalf("Processing JSON dot path should return the same non-nil builder")
		}
	})
}
