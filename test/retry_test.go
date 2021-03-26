package test

import (
	"fmt"
	"github.com/pkg/errors"
	"rchain/examples/retry"
	"testing"
)

func TestRetry(t *testing.T) {
	retry.FormatNumber()
	retry.GetInfos("Tom")
	retry.SumNumber(1, 2)
}

func TestRetry2(t *testing.T) {
	retry.Invoke(4, 1, func(num int) error {
		fmt.Println("num ", num)
		retry.SumNumber(2, 3)
		return errors.New("a")
	})
}
