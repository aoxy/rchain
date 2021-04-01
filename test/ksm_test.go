package test

import (
	"fmt"
	"testing"
	"time"
)

func TestGetBlock(t *testing.T) {

	var globalNumber int64

	go func() {
		for {
			globalNumber += 2
			time.Sleep(time.Second * 2)
		}
	}()

	// var ch = make(chan int64, 1)

	var waitIndex int64 = 3
	var handleIndex int64 = 0
	var flag = true
	for {
		curIndex := globalNumber
		if handleIndex > waitIndex {
			if flag {
				handleIndex -= waitIndex
			}
		}

		if handleIndex > curIndex {
			time.Sleep(time.Second)
			continue
		}

		for blockIndex := handleIndex + 1; blockIndex <= curIndex; blockIndex++ {
			fmt.Println("遍历", blockIndex, handleIndex, curIndex)

			time.Sleep(time.Second)
			handleIndex = blockIndex
		}

		time.Sleep(time.Second)
		flag = false
		fmt.Println("下一轮", handleIndex)
		fmt.Println()
	}
}












