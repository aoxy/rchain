package service

import (
	"fmt"
	"github.com/pkg/errors"
	"rchain/internal/utils"
)

func PullBlocks() {
	// 拉取最近的block信息, 拿到最新的blockNumber
	var topBlockNumber int64
	fmt.Println("开始拉取块")
	utils.Invoke(4, 1, func(n int) error {
		blockInfos := utils.GetBlocks(n, 1)
		if len(blockInfos) <= 0 {
			return errors.New("non block")
		}
		topBlockNumber = blockInfos[0].BlockNumber
		fmt.Println("拉取成功", topBlockNumber)
		return nil
	})

	fmt.Println("最近的BlockNumber是：", topBlockNumber)
}
