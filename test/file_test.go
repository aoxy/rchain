package test

import (
	"fmt"
	"rchain/internal/utils"
	"strconv"
	"testing"
)

func TestOpFile(t *testing.T) {
	file := "E:\\work\\rchain\\docs\\output\\.latest_blocknumber"
	data := utils.ReadFile(file)
	if data != nil {
		fmt.Println(data)
		fmt.Println(string(data) + "aaa")
		fmt.Println(string(data[:]) + "bbb")

		s := string(data)
		no, _ := strconv.ParseInt(s, 10, 64)
		fmt.Println("no: ", no)
	}
}
