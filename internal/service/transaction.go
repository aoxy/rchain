package service

import (
	casperV1 "dawn1806/rchain/api/pb"
	"dawn1806/rchain/internal/dto"
	"dawn1806/rchain/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var lastBlockNumber int64

type Transaction struct {
	DSClient       casperV1.DeployServiceClient
	httpClient     *http.Client
	TransactionUrl string
	BlockHash      chan string
	LastBlockNo    int64
	ttl            int
}

func (t *Transaction) GetBlockHash() error {

	blockInfo, err := utils.LastFinalizedBlock(t.DSClient)
	if err != nil {
		log.Println("501 获取block信息失败，", err)
		return err
	}

	// 如果是第一次获取，则获取最新的10个区块
	if t.LastBlockNo == 0 {
		t.LastBlockNo = blockInfo.BlockNumber - 10
	}

	// 先判断blockNumber是否连续
	// 取出上一次处理的blockNumber进行对比，如果是它的下一个block，则放入管道进行处理
	// 如果不是，则先把漏掉的block取出进行处理
	needBlockNumber := t.LastBlockNo + 1
	file := "./docs/output/log.txt"
	if blockInfo.BlockNumber == needBlockNumber {
		s := "11\t" + strconv.Itoa(int(blockInfo.BlockNumber)) + "\t" + blockInfo.BlockHash + "\n"
		utils.WriteFile(file, s)
		t.BlockHash <- blockInfo.BlockHash
	} else if blockInfo.BlockNumber > needBlockNumber {
		// 把漏掉的block取出来进行处理
		blockInfos, err := utils.GetBlocksByHeightClient(t.DSClient, needBlockNumber, blockInfo.BlockNumber)
		if err != nil {
			log.Println("502 获取block信息失败，", err)
			return err
		}

		if len(blockInfos) == 0 {
			close(t.BlockHash)
			return errors.New("获取block失败")
		}

		for _, blockInfo = range blockInfos {
			s := "22\t" + strconv.Itoa(int(blockInfo.BlockNumber)) + "\t" + blockInfo.BlockHash + "\n"
			utils.WriteFile(file, s)
			t.BlockHash <- blockInfo.BlockHash
		}
	} else {
		log.Println("最新的blockNumber小于已处理的blockNumber")
	}
	return nil
}

func (t *Transaction) GetTransferInfo() error {
	blockHash := <-t.BlockHash
	if blockHash == "" {
		return errors.New("取不到blockHash")
	}
	url := t.TransactionUrl + blockHash
	resp, err := t.httpClient.Get(url)
	if err != nil {
		log.Fatal("501 GetTransferInfo error ", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("502 GetTransferInfo error ", err)
		return err
	}

	var transInfos = make([][]dto.BlockTransferInfo, 0)
	err = json.Unmarshal(body, &transInfos)
	if err != nil {
		log.Fatal("503 GetTransferInfo error ", err)
		return err
	}

	for _, infos := range transInfos {
		for _, info := range infos {
			s := fmt.Sprintf("blockHash: %s\nfrom: %s\nto: %s\namount: %d\ntimestamp: %d\nvalidAfterBlockNumber: %d\n\n",
				blockHash, info.FromAddr, info.ToAddr, info.Amount, info.Deploy.Timestamp, info.Deploy.ValidAfterBlockNumber)
			file := viper.GetString("global.output") + "block_transaction.txt"
			utils.WriteFile(file, s)
		}
	}
	return nil
}
