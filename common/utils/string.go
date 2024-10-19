package utils

import (
	"fmt"
)

// 拼接in查询参数
func JoinInParamsForNumber[T int | int64 | string](list []T) string {
	var sql string
	for i, item := range list {
		if i == len(list)-1 {
			sql += fmt.Sprintf("%v", item)
		} else {
			sql += fmt.Sprintf("%v,", item)
		}
	}
	return sql
}
