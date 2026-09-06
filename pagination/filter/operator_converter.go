package filter

import (
	"strings"

	"github.com/kalandramo/bald-utils/stringcase"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
)

var operatorMap = map[string]storev1.Operator{
	"eq":     storev1.Operator_EQ,
	"equal":  storev1.Operator_EQ,
	"equals": storev1.Operator_EQ,

	"ne":         storev1.Operator_NEQ,
	"neq":        storev1.Operator_NEQ,
	"not":        storev1.Operator_NEQ,
	"not_equal":  storev1.Operator_NEQ,
	"not_equals": storev1.Operator_NEQ,
	"not-equal":  storev1.Operator_NEQ,

	"gt":           storev1.Operator_GT,
	"greater_than": storev1.Operator_GT,
	"greater-than": storev1.Operator_GT,

	"gte":                   storev1.Operator_GTE,
	"greater_than_or_equal": storev1.Operator_GTE,
	"greater_equals":        storev1.Operator_GTE,
	"greater_or_equal":      storev1.Operator_GTE,
	"greater-or-equal":      storev1.Operator_GTE,

	"lt":        storev1.Operator_LT,
	"less_than": storev1.Operator_LT,
	"less-than": storev1.Operator_LT,

	"lte":                storev1.Operator_LTE,
	"less_than_or_equal": storev1.Operator_LTE,
	"less_equals":        storev1.Operator_LTE,
	"less_or_equal":      storev1.Operator_LTE,
	"less-or-equal":      storev1.Operator_LTE,

	"like": storev1.Operator_LIKE,

	"ilike":  storev1.Operator_ILIKE,
	"i_like": storev1.Operator_ILIKE,

	"not_like": storev1.Operator_NOT_LIKE,
	"notlike":  storev1.Operator_NOT_LIKE,

	"in": storev1.Operator_IN,

	"nin":    storev1.Operator_NIN,
	"not_in": storev1.Operator_NIN,
	"notin":  storev1.Operator_NIN,

	"is_null": storev1.Operator_IS_NULL,
	"isnull":  storev1.Operator_IS_NULL,

	"is_not_null": storev1.Operator_IS_NOT_NULL,
	"isnot_null":  storev1.Operator_IS_NOT_NULL,
	"isnotnull":   storev1.Operator_IS_NOT_NULL,
	"not_isnull":  storev1.Operator_IS_NOT_NULL,

	"between": storev1.Operator_BETWEEN,
	"range":   storev1.Operator_BETWEEN,

	"regexp": storev1.Operator_REGEXP,
	"regex":  storev1.Operator_REGEXP,

	"iregexp":  storev1.Operator_IREGEXP,
	"i_regexp": storev1.Operator_IREGEXP,
	"iregex":   storev1.Operator_IREGEXP,

	"contains": storev1.Operator_CONTAINS,

	"icontains":  storev1.Operator_ICONTAINS,
	"i_contains": storev1.Operator_ICONTAINS,

	"starts_with": storev1.Operator_STARTS_WITH,
	"startswith":  storev1.Operator_STARTS_WITH,

	"istarts_with":  storev1.Operator_ISTARTS_WITH,
	"i_starts_with": storev1.Operator_ISTARTS_WITH,
	"istartswith":   storev1.Operator_ISTARTS_WITH,

	"ends_with": storev1.Operator_ENDS_WITH,
	"endswith":  storev1.Operator_ENDS_WITH,

	"iends_with":  storev1.Operator_IENDS_WITH,
	"i_ends_with": storev1.Operator_IENDS_WITH,
	"iendswith":   storev1.Operator_IENDS_WITH,

	"json_contains":  storev1.Operator_JSON_CONTAINS,
	"array_contains": storev1.Operator_ARRAY_CONTAINS,
	"exists":         storev1.Operator_EXISTS,
	"search":         storev1.Operator_SEARCH,
	"exact":          storev1.Operator_EXACT,

	"iexact":  storev1.Operator_IEXACT,
	"i_exact": storev1.Operator_IEXACT,
}

// ConverterStringToOperator 将字符串转换为 storev1.Operator 枚举值
func ConverterStringToOperator(str string) storev1.Operator {
	key := strings.ToLower(stringcase.ToSnakeCase(str))
	if v, ok := operatorMap[key]; ok {
		return v
	}
	return storev1.Operator_OPERATOR_UNSPECIFIED
}

// IsValidOperatorString 检查字符串是否为有效的 storev1.Operator 枚举值
func IsValidOperatorString(str string) bool {
	op := ConverterStringToOperator(str)
	return op != storev1.Operator_OPERATOR_UNSPECIFIED
}

var datePartMap = map[string]storev1.DatePart{
	"date": storev1.DatePart_DATE,

	"year": storev1.DatePart_YEAR,
	"yr":   storev1.DatePart_YEAR,

	"iso_year": storev1.DatePart_ISO_YEAR,
	"iso-year": storev1.DatePart_ISO_YEAR,

	"quarter": storev1.DatePart_QUARTER,
	"month":   storev1.DatePart_MONTH,
	"week":    storev1.DatePart_WEEK,

	"week_day": storev1.DatePart_WEEK_DAY,
	"week-day": storev1.DatePart_WEEK_DAY,
	"weekday":  storev1.DatePart_WEEK_DAY,

	"iso_week_day": storev1.DatePart_ISO_WEEK_DAY,
	"iso-week-day": storev1.DatePart_ISO_WEEK_DAY,

	"day":  storev1.DatePart_DAY,
	"time": storev1.DatePart_TIME,
	"hour": storev1.DatePart_HOUR,

	"minute": storev1.DatePart_MINUTE,
	"min":    storev1.DatePart_MINUTE,

	"second": storev1.DatePart_SECOND,
	"sec":    storev1.DatePart_SECOND,

	"microsecond": storev1.DatePart_MICROSECOND,
}

// ConverterStringToDatePart 将字符串转换为 storev1.DatePart 枚举
func ConverterStringToDatePart(s string) *storev1.DatePart {
	key := strings.ToLower(stringcase.ToSnakeCase(s))
	if v, ok := datePartMap[key]; ok {
		return &v
	}
	return nil
}

// ConverterDatePartToString 将 storev1.DatePart 枚举转换为字符串
func ConverterDatePartToString(datePart *storev1.DatePart) string {
	if datePart == nil {
		return ""
	}
	for k, v := range datePartMap {
		if v == *datePart {
			return k
		}
	}
	return ""
}

// IsValidDatePartString 检查字符串是否为有效的 storev1.DatePart 枚举值
func IsValidDatePartString(str string) bool {
	dp := ConverterStringToDatePart(str)
	return dp != nil
}
