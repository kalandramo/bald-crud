package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
)

func TestProcessor_BuildOpenSearchQuery_AllOperators(t *testing.T) {
	proc := NewProcessor()
	ops := []struct {
		op   storev1.Operator
		val  string
		vals []any
		want string // 断言结构类型
	}{
		{storev1.Operator_EQ, "v", nil, "term"},
		{storev1.Operator_NEQ, "v", nil, "bool"},
		{storev1.Operator_IN, "", []any{"a", "b"}, "terms"},
		{storev1.Operator_NIN, "", []any{"a", "b"}, "bool"},
		{storev1.Operator_GTE, "1", nil, "range"},
		{storev1.Operator_GT, "1", nil, "range"},
		{storev1.Operator_LTE, "1", nil, "range"},
		{storev1.Operator_LT, "1", nil, "range"},
		{storev1.Operator_BETWEEN, "", []any{"1", "2"}, "range"},
		{storev1.Operator_IS_NULL, "", nil, "bool"},
		{storev1.Operator_IS_NOT_NULL, "", nil, "exists"},
		{storev1.Operator_CONTAINS, "v", nil, "match_phrase"},
		{storev1.Operator_ICONTAINS, "v", nil, "match_phrase"},
		{storev1.Operator_STARTS_WITH, "v", nil, "prefix"},
		{storev1.Operator_ISTARTS_WITH, "v", nil, "prefix"},
		{storev1.Operator_ENDS_WITH, "v", nil, "wildcard"},
		{storev1.Operator_IENDS_WITH, "v", nil, "wildcard"},
		{storev1.Operator_EXACT, "v", nil, "term"},
		{storev1.Operator_IEXACT, "v", nil, "term"},
		{storev1.Operator_REGEXP, "v", nil, "regexp"},
		{storev1.Operator_IREGEXP, "v", nil, "regexp"},
		{storev1.Operator_SEARCH, "v", nil, "query_string"},
	}
	for _, tc := range ops {
		cond := &storev1.FilterCondition{
			Field:      "f",
			Op:         tc.op,
			ValueOneof: &storev1.FilterCondition_Value{Value: tc.val},
			Values:     toStringSlice(tc.vals),
		}
		got := proc.buildCond(cond)
		if tc.want != "" {
			if assert.NotNil(t, got, "op %v", tc.op) {
				found := false
				for k := range got {
					if k == tc.want {
						found = true
						break
					}
				}
				assert.True(t, found, "expect key %s for op %v, got %v", tc.want, tc.op, got)
			}
		}
	}
}

func TestProcessor_BuildOpenSearchQuery_AND_OR(t *testing.T) {
	proc := NewProcessor()
	// AND
	andExpr := &storev1.FilterExpr{
		Type: storev1.ExprType_AND,
		Conditions: []*storev1.FilterCondition{
			{Field: "f1", Op: storev1.Operator_EQ, ValueOneof: &storev1.FilterCondition_Value{Value: "v1"}},
			{Field: "f2", Op: storev1.Operator_GT, ValueOneof: &storev1.FilterCondition_Value{Value: "2"}},
		},
	}
	q := proc.BuildOpenSearchQuery(andExpr)
	assert.NotNil(t, q)
	boolQ := q["bool"].(map[string]any)
	assert.Len(t, boolQ["must"], 2)
	// OR
	orExpr := &storev1.FilterExpr{
		Type: storev1.ExprType_OR,
		Conditions: []*storev1.FilterCondition{
			{Field: "f1", Op: storev1.Operator_EQ, ValueOneof: &storev1.FilterCondition_Value{Value: "v1"}},
			{Field: "f2", Op: storev1.Operator_GT, ValueOneof: &storev1.FilterCondition_Value{Value: "2"}},
		},
	}
	q2 := proc.BuildOpenSearchQuery(orExpr)
	assert.NotNil(t, q2)
	boolQ2 := q2["bool"].(map[string]any)
	assert.Len(t, boolQ2["should"], 2)
}

func toStringSlice(vals []any) []string {
	if vals == nil {
		return nil
	}
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = v.(string)
	}
	return out
}
