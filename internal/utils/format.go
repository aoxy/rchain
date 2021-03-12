package utils

import (
	"fmt"
	"reflect"
)

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
