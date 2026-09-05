package sorting

import (
	"encoding/json"
	"strings"

	storev1 "github.com/kalandramo/bald/bconf/gen/go/bald/store/v1"
	"github.com/tx7do/go-utils/stringcase"
	"go.einride.tech/aip/ordering"
)

type OrderByStringConverter struct {
}

func NewOrderByStringConverter() *OrderByStringConverter {
	return &OrderByStringConverter{}
}

// Convert 将排序字符串转换为排序对象列表
func (obc OrderByStringConverter) Convert(orderBy string) ([]*storev1.Sorting, error) {
	if len(orderBy) == 0 {
		return nil, nil
	}

	if strings.HasPrefix(orderBy, "[") && strings.HasSuffix(orderBy, "]") {
		// JSON 格式
		return obc.ParseJsonString(orderBy)
	}

	// AIP 格式
	return obc.ParseAIPString(orderBy)
}

// ParseJsonString 解析 JSON 格式的排序字符串
func (obc OrderByStringConverter) ParseJsonString(orderByJson string) ([]*storev1.Sorting, error) {
	if len(orderByJson) == 0 {
		return nil, nil
	}

	var strSlice []string
	var sortings []*storev1.Sorting

	// 反序列化
	err := json.Unmarshal([]byte(orderByJson), &strSlice)
	if err != nil {
		return nil, err
	}

	var isDesc bool
	var field string
	for _, item := range strSlice {
		item = strings.TrimSpace(item)
		if len(item) == 0 {
			continue
		}

		field = item

		if strings.HasPrefix(item, "-") {
			// 降序
			field = item[1:]
			if len(field) == 0 {
				continue
			}

			isDesc = true

		} else {
			// 升序
			if len(field) == 0 {
				continue
			}

			isDesc = false
		}

		field = strings.TrimSpace(field)
		field = stringcase.ToSnakeCase(field)

		if !isDesc {
			sortings = append(sortings, &storev1.Sorting{
				Field:     field,
				Direction: storev1.Sorting_ASC,
			})
			continue
		} else {
			sortings = append(sortings, &storev1.Sorting{
				Field:     field,
				Direction: storev1.Sorting_DESC,
			})
		}
	}

	return sortings, err
}

// ParseAIPString 解析 AIP 格式的排序字符串
func (obc OrderByStringConverter) ParseAIPString(orderByString string) ([]*storev1.Sorting, error) {
	if len(orderByString) == 0 {
		return nil, nil
	}

	var actual ordering.OrderBy
	err := actual.UnmarshalString(orderByString)
	if err != nil {
		return nil, err
	}

	var sortings []*storev1.Sorting
	for _, item := range actual.Fields {
		var direction storev1.Sorting_Direction
		if item.Desc {
			direction = storev1.Sorting_DESC
		} else {
			direction = storev1.Sorting_ASC
		}

		sortings = append(sortings, &storev1.Sorting{
			Field:     item.Path,
			Direction: direction,
		})
	}

	return sortings, nil
}
