package retry

import (
	"fmt"
	"time"
)

type Info struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Desc string `json:"desc"`
}

func Invoke(retryNum int, sec time.Duration, fn func(num int) error) {
	for i := 0; i < retryNum; i++ {
		err := fn(i)
		if err == nil {
			return
		}
		time.Sleep(sec * time.Second)
	}
}

func SumNumber(a int, b int) int {
	fmt.Println("a + b =", a+b)
	return a + b
}

func GetInfos(name string) []*Info {
	var infos = make([]*Info, 0)
	info := &Info{
		Name: name,
		Age:  30,
		Desc: "测试数据",
	}
	infos = append(infos, info)
	fmt.Println("infos ", infos)
	return infos
}

func FormatNumber() {
	fmt.Println("FormatNumber")
}
