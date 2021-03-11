package utils

import (
	"context"
	casperV1 "dawn1806/rchain/api/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func NewDeployServiceClient(target string) (casperV1.DeployServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	DSClient := casperV1.NewDeployServiceClient(conn)
	return DSClient, conn, nil
}

func GetBlocksByHeightClient(DSClient casperV1.DeployServiceClient, start int64, end int64) ([]*casperV1.LightBlockInfo, error) {

	blocksClient, err := DSClient.GetBlocksByHeights(context.Background(), &casperV1.BlocksQueryByHeight{
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
func LastFinalizedBlock(DSClient casperV1.DeployServiceClient) (*casperV1.LightBlockInfo, error) {
	resp, err := DSClient.LastFinalizedBlock(context.Background(), &casperV1.LastFinalizedBlockQuery{})
	if err != nil {
		log.Fatal("LastFinalizedBlock error ", err)
		return nil, err
	}
	return resp.GetBlockInfo().GetBlockInfo(), nil
}

func GetBlocks(DSClient casperV1.DeployServiceClient) {
	// 获取最近的区块
	//_, err := t.DSClient.GetBlocks(context.Background(), &casperV1.BlocksQuery{Depth: 20})
	//if err != nil {
	//	log.Fatal("GetBlocks error ", err)
	//}
}
