package utils

import (
	"fmt"
	"strconv"
	"strings"
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

func ConvertIntsToStrings(ints []int) []string {
	var strings []string
	for _, num := range ints {
		str := strconv.Itoa(num)
		strings = append(strings, str)
	}
	return strings
}

func ConvertStringsToInts(strings []string) ([]int, error) {
	var ints []int
	for _, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		ints = append(ints, num)
	}
	return ints, nil
}

// JoinString 拼接字符串
func JoinString(s ...string) string {
	// strings.Join(s, "")
	var b strings.Builder
	for _, str := range s {
		b.WriteString(str)
	}
	return b.String()
}
