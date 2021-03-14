package utils

import (
	"fmt"
	"reflect"
)

// 递归解析结构体，输出切片数组，解决嵌套结构体解析问题
func FormatStruct(v reflect.Value, s []string) []string {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Deploy" {
			s = FormatStruct(v.Field(i), s)
			continue
		}
		if t.Field(i).Name == "Term" {
			continue
		}
		s = append(s, fmt.Sprintf("%v", v.Field(i)))
		fmt.Println("sss: ", s)
	}
	return s
}
