package filter

import (
	"testing"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
	"github.com/tx7do/go-utils/trans"
)

func TestConverterStringToOperator(t *testing.T) {
	cases := map[string]storev1.Operator{
		"eq":             storev1.Operator_EQ,
		"EQ":             storev1.Operator_EQ,
		"equal":          storev1.Operator_EQ,
		"equals":         storev1.Operator_EQ,
		"ne":             storev1.Operator_NEQ,
		"not-equal":      storev1.Operator_NEQ,
		"not_equal":      storev1.Operator_NEQ,
		"gt":             storev1.Operator_GT,
		"greater-than":   storev1.Operator_GT,
		"gte":            storev1.Operator_GTE,
		"less_than":      storev1.Operator_LT,
		"like":           storev1.Operator_LIKE,
		"iLike":          storev1.Operator_ILIKE,
		"i_like":         storev1.Operator_ILIKE,
		"in":             storev1.Operator_IN,
		"notin":          storev1.Operator_NIN,
		"isNotNull":      storev1.Operator_IS_NOT_NULL,
		"isnull":         storev1.Operator_IS_NULL,
		"between":        storev1.Operator_BETWEEN,
		"regexp":         storev1.Operator_REGEXP,
		"iregex":         storev1.Operator_IREGEXP,
		"contains":       storev1.Operator_CONTAINS,
		"icontains":      storev1.Operator_ICONTAINS,
		"startsWith":     storev1.Operator_STARTS_WITH,
		"ends_with":      storev1.Operator_ENDS_WITH,
		"json_contains":  storev1.Operator_JSON_CONTAINS,
		"array_contains": storev1.Operator_ARRAY_CONTAINS,
		"exists":         storev1.Operator_EXISTS,
		"search":         storev1.Operator_SEARCH,
		"exact":          storev1.Operator_EXACT,
		"iexact":         storev1.Operator_IEXACT,

		// unknown / empty -> unspecified
		"":       storev1.Operator_OPERATOR_UNSPECIFIED,
		"foobar": storev1.Operator_OPERATOR_UNSPECIFIED,
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			got := ConverterStringToOperator(input)
			if got != want {
				t.Fatalf("ConverterStringToOperator(%q) = %v, want %v", input, got, want)
			}
		})
	}
}

func TestIsValidOperatorString(t *testing.T) {
	valid := []string{"eq", "not_equal", "i_like", "search"}
	for _, s := range valid {
		if !IsValidOperatorString(s) {
			t.Fatalf("IsValidOperatorString(%q) = false, want true", s)
		}
	}

	invalid := []string{"", "unknown_op", "blah"}
	for _, s := range invalid {
		if IsValidOperatorString(s) {
			t.Fatalf("IsValidOperatorString(%q) = true, want false", s)
		}
	}
}

func TestConverterStringToDatePart(t *testing.T) {
	cases := map[string]*storev1.DatePart{
		"date":         trans.Ptr(storev1.DatePart_DATE),
		"Date":         trans.Ptr(storev1.DatePart_DATE),
		"DATE":         trans.Ptr(storev1.DatePart_DATE),
		"year":         trans.Ptr(storev1.DatePart_YEAR),
		"yr":           trans.Ptr(storev1.DatePart_YEAR),
		"iso_year":     trans.Ptr(storev1.DatePart_ISO_YEAR),
		"iso-year":     trans.Ptr(storev1.DatePart_ISO_YEAR),
		"quarter":      trans.Ptr(storev1.DatePart_QUARTER),
		"month":        trans.Ptr(storev1.DatePart_MONTH),
		"week":         trans.Ptr(storev1.DatePart_WEEK),
		"week_day":     trans.Ptr(storev1.DatePart_WEEK_DAY),
		"week-day":     trans.Ptr(storev1.DatePart_WEEK_DAY),
		"weekday":      trans.Ptr(storev1.DatePart_WEEK_DAY),
		"iso_week_day": trans.Ptr(storev1.DatePart_ISO_WEEK_DAY),
		"iso-week-day": trans.Ptr(storev1.DatePart_ISO_WEEK_DAY),
		"day":          trans.Ptr(storev1.DatePart_DAY),
		"time":         trans.Ptr(storev1.DatePart_TIME),
		"hour":         trans.Ptr(storev1.DatePart_HOUR),
		"minute":       trans.Ptr(storev1.DatePart_MINUTE),
		"min":          trans.Ptr(storev1.DatePart_MINUTE),
		"second":       trans.Ptr(storev1.DatePart_SECOND),
		"sec":          trans.Ptr(storev1.DatePart_SECOND),
		"microsecond":  trans.Ptr(storev1.DatePart_MICROSECOND),

		// unknown / empty -> unspecified
		"":       nil,
		"foobar": nil,
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			got := ConverterStringToDatePart(input)
			if want == nil && got != nil {
				t.Fatalf("ConverterStringToDatePart(%q) = %v, want %v", input, got, want)
			} else if want != nil && *got != *want {
				t.Fatalf("ConverterStringToDatePart(%q) = %v, want %v", input, got, want)
			}
		})
	}
}

func TestIsValidDatePartString(t *testing.T) {
	valid := []string{"date", "year", "iso_year", "minute", "microsecond"}
	for _, s := range valid {
		if !IsValidDatePartString(s) {
			t.Fatalf("IsValidDatePartString(%q) = false, want true", s)
		}
	}

	invalid := []string{"", "not_a_part", "blah"}
	for _, s := range invalid {
		if IsValidDatePartString(s) {
			t.Fatalf("IsValidDatePartString(%q) = true, want false", s)
		}
	}
}

func TestConverterDatePartToString(t *testing.T) {
	ptr := func(v storev1.DatePart) *storev1.DatePart { return &v }

	cases := []struct {
		in   *storev1.DatePart
		want []string
	}{
		{in: ptr(storev1.DatePart_DATE), want: []string{"date"}},
		{in: ptr(storev1.DatePart_YEAR), want: []string{"year", "yr"}},
		{in: ptr(storev1.DatePart_ISO_YEAR), want: []string{"iso_year", "iso-year"}},
		{in: ptr(storev1.DatePart_QUARTER), want: []string{"quarter"}},
		{in: ptr(storev1.DatePart_MONTH), want: []string{"month"}},
		{in: ptr(storev1.DatePart_WEEK), want: []string{"week"}},
		{in: ptr(storev1.DatePart_WEEK_DAY), want: []string{"week_day", "week-day", "weekday"}},
		{in: ptr(storev1.DatePart_ISO_WEEK_DAY), want: []string{"iso_week_day", "iso-week-day"}},
		{in: ptr(storev1.DatePart_DAY), want: []string{"day"}},
		{in: ptr(storev1.DatePart_TIME), want: []string{"time"}},
		{in: ptr(storev1.DatePart_HOUR), want: []string{"hour"}},
		{in: ptr(storev1.DatePart_MINUTE), want: []string{"minute", "min"}},
		{in: ptr(storev1.DatePart_SECOND), want: []string{"second", "sec"}},
		{in: ptr(storev1.DatePart_MICROSECOND), want: []string{"microsecond"}},

		// nil input -> empty string
		{in: nil, want: []string{""}},

		// unknown enum value -> empty string
		{in: ptr(storev1.DatePart(9999)), want: []string{""}},
	}

	for _, c := range cases {
		got := ConverterDatePartToString(c.in)
		ok := false
		for _, w := range c.want {
			if got == w {
				ok = true
				break
			}
		}
		if !ok {
			t.Fatalf("ConverterDatePartToString(%v) = %q, want one of %v", c.in, got, c.want)
		}
	}
}
