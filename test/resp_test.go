package test

import (
	"encoding/json"
	"fmt"
	"testing"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestResp(t *testing.T) {
	str := "[[{\"name\": \"dawn\", \"age\": 30}], [{\"name\": \"dawn002\", \"age\": 20}]]"
	users := make([][]user, 1)
	err := json.Unmarshal([]byte(str), &users)
	if err != nil {
		fmt.Println("json.Unmarshal error ", err)
	}
	fmt.Println("pass")
}
