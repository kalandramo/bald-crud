package query

import (
	"fmt"
	"reflect"
	"strings"
)

// RawValue 显式豁免格式化的裸表达式值（如 `now() - interval '15 minutes'`）。
// CD1 引入：BuildQueryWithParams 默认对 value 安全格式化（引号转义），
// 需要拼接 Flux 表达式的调用方必须显式包装——与 bald bundle.Raw 同款
//「默认安全 + 显式豁免」哲学。误用风险自负：内容将原样进入查询。
type RawValue string

// FormatValue 根据类型格式化值；slice 会被格式化为 "(v1,v2,...)"。
// 导出于 CD1 修复：utils.BuildQueryWithParams 原以 %v 裸拼 value——
// 注入面（值含引号/分号即改写查询语义）；统一走本安全格式化。
func FormatValue(v any) string {
	if v == nil {
		return "NULL"
	}

	switch t := v.(type) {
	case RawValue:
		return string(t)
	case string:
		return fmt.Sprintf("'%s'", escapeString(t))
	case bool:
		if t {
			return "true"
		}
		return "false"
	case fmt.Stringer:
		return fmt.Sprintf("'%s'", escapeString(t.String()))
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			parts := make([]string, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				parts = append(parts, FormatValue(rv.Index(i).Interface()))
			}
			return "(" + strings.Join(parts, ",") + ")"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%v", v)
		default:
			// fallback to quoted string
			return fmt.Sprintf("'%s'", escapeString(fmt.Sprintf("%v", v)))
		}
	}
}

// formatRegex 将值包装为 /.../ 形式，用于 =~ 操作
func formatRegex(v any) string {
	s := ""
	switch t := v.(type) {
	case string:
		s = t
	default:
		s = fmt.Sprintf("%v", v)
	}
	// 简单转义斜线
	s = strings.ReplaceAll(s, "/", "\\/")
	return fmt.Sprintf("/%s/", s)
}

// escapeString 转义单引号等
func escapeString(s string) string {
	// 转义单引号和反斜杠
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}
