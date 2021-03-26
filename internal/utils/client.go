package utils

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	casperV1 "rchain/api/pb"
	"time"
)

var Observers = []string{
	"observer-asia.services.mainnet.rchain.coop:4040",
	"observer-eu.services.mainnet.rchain.coop:40401",
	"observer-us.services.mainnet.rchain.coop:40401",
}

var getTransactionUrl = "http://observer-exch2.services.mainnet.rchain.coop:7070/getTransaction/"

func Invoke(conRetryNum int, t time.Duration, fn func(n int) error) {
	for i := 0; i < conRetryNum; i++ {
		err := fn(i)
		if err == nil {
			return
		}
		if i == conRetryNum-1 {
			fmt.Println("\n连接失败，请检查连接地址或网络")
			return
		}
		fmt.Println("\n", t.String(), "秒后重试连接")
		time.Sleep(t * time.Second)
	}
}

func NewDeployServiceClient(target string) (casperV1.DeployServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	DSClient := casperV1.NewDeployServiceClient(conn)
	return DSClient, conn, nil
}

func GetBlocksByHeightClient(i int, start int64, end int64) ([]*casperV1.LightBlockInfo, error) {
	client, conn, _ := NewDeployServiceClient(Observers[i%len(Observers)])
	defer conn.Close()

	blocksClient, err := client.GetBlocksByHeights(context.Background(), &casperV1.BlocksQueryByHeight{
		StartBlockNumber: start,
		EndBlockNumber:   end,
	})

	if err != nil {
		log.Fatal("GetBlocksByHeights error ", err)
		return nil, err
	}

	var blockInfos = make([]*casperV1.LightBlockInfo, 0)

	for {
		blockInfoResp, err := blocksClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("blocksClient.Recv error ", err)
			return nil, err
		}

		lightBlockInfo := blockInfoResp.GetBlockInfo()
		blockInfos = append(blockInfos, lightBlockInfo)
	}

	return blockInfos, nil
}

// 获取最新验证通过的block
func LastFinalizedBlock(i int) (*casperV1.LightBlockInfo, error) {
	client, conn, _ := NewDeployServiceClient(Observers[i%len(Observers)])
	defer conn.Close()

	resp, err := client.LastFinalizedBlock(context.Background(), &casperV1.LastFinalizedBlockQuery{})
	if err != nil {
		log.Fatal("LastFinalizedBlock error ", err)
		return nil, err
	}
	return resp.GetBlockInfo().GetBlockInfo(), nil
}

func DoDeploy(i int, deploy *casperV1.DeployDataProto) (*casperV1.DeployResponse, error) {
	client, conn, _ := NewDeployServiceClient(Observers[i%len(Observers)])
	defer conn.Close()

	resp, err := client.DoDeploy(context.Background(), deploy)
	if err != nil {
		log.Fatal("DoDeploy error ", err)
		return nil, err
	}
	return resp, nil
}

// 获取最近的区块
func GetBlocks(i int, depth int32) []*casperV1.LightBlockInfo {
	fmt.Println("连接地址：", Observers[i%len(Observers)])
	client, conn, _ := NewDeployServiceClient(Observers[i%len(Observers)])
	defer conn.Close()

	blocksClient, err := client.GetBlocks(context.Background(), &casperV1.BlocksQuery{Depth: depth})
	if err != nil {
		log.Println("GetBlocks error ", err)
		return nil
	}

	var blockInfos = make([]*casperV1.LightBlockInfo, 0)
	for {
		resp, err := blocksClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("GetBlocks Recv error", err)
			return nil
		}

		blockInfo := resp.GetBlockInfo()
		blockInfos = append(blockInfos, blockInfo)
	}
	return blockInfos
}
