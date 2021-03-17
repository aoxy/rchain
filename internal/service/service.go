package service

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	casperV1 "rchain/api/pb"
	"rchain/internal/utils"
	"strconv"
	"sync"
	"time"
)

var t *Transaction

func initConfig() {
	hosts := viper.GetStringSlice("observer.hosts")
	grpcPort := viper.GetString("observer.grpcPort")
	httpPort := viper.GetString("observer.httpPort")
	grpcUrl := hosts[0] + ":" + grpcPort
	httpUrl := "http://" + hosts[3] + ":" + httpPort + "/getTransaction/"

	DSClient, conn, err := utils.NewDeployServiceClient(grpcUrl)
	if err != nil {
		fmt.Println("NewDeployServiceClient error ", err)
		return
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	t = &Transaction{
		DSClient:       DSClient,
		httpClient:     client,
		Conn:           conn,
		TransactionUrl: httpUrl,
		BlockInfoCh:    make(chan *casperV1.LightBlockInfo, 1),
	}

	// 读取最新blockNumber文件[文件来源：获取交易数据后会更新该文件内容]，初始化LastBlockNo
	// 这个blockNumber是上一次以及处理的最后一个blockNumber
	// 目的：当发生错误、重启等操作后，为了保证交易信息的连续性
	data := utils.ReadFile(viper.GetString("global.output") + ".latest_blocknumber")
	fmt.Println("data: ", data)
	if data != nil {
		no, err := strconv.ParseInt(string(data), 10, 64)
		fmt.Println("data no: ", no)
		fmt.Println("err: ", err)
		if err == nil {
			t.LastBlockNo = no
		}
	}
}

func Run() {
	// 初始化配置
	initConfig()
	// 处理交易信息
	HandleTransaction()
}

func HandleTransaction() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	// 启动一个协程，获取最新的block信息，放入管道
	go GetBlockInfo(wg)
	// 从管道取出block信息，调用http接口获得交易信息，将交易信息写入文件
	go GetTransferInfo(wg)
	wg.Wait()
	log.Println("处理交易信息结束")
}

// 功能：循环获取block信息，放入管道，供其他协程使用
// todo：循环中断的错误处理、重新唤醒协程；
// todo：考虑使用定时器
// todo：考虑将block信息持久化
func GetBlockInfo(wg *sync.WaitGroup) {
	defer wg.Done()
	defer t.Conn.Close()
	defer close(t.BlockInfoCh)

	fmt.Println("开始获取block信息...", t.LastBlockNo)
	// 定义当前需要处理的blockNumber
	var needBlockNumber int64
	for {
		// 获取最新验证通过的block信息
		blockInfo, err := utils.LastFinalizedBlock(t.DSClient)
		if err != nil {
			log.Println("501 获取block信息失败，", err)
			break
		}

		// 如果是第一次获取，则获取最新的10个区块
		// t.LastBlockNo第一次获取初始化为0，当写文件成功后更新为最新已处理的blockNumber
		if t.LastBlockNo == 0 {
			t.LastBlockNo = blockInfo.BlockNumber - 10
		}

		// 判断blockNumber是否连续
		// 如果连续，则放入管道交给其他协程处理，然后继续获取下一个
		if needBlockNumber == 0 {
			needBlockNumber = t.LastBlockNo + 1
		}
		fmt.Printf("numbers: %d\t%d\t%d", t.LastBlockNo, needBlockNumber, blockInfo.BlockNumber)
		// file := "./docs/output/log.txt"
		if blockInfo.BlockNumber == needBlockNumber {
			//s := "11\t" + strconv.Itoa(int(blockInfo.BlockNumber)) + "\t" + blockInfo.BlockHash + "\n"
			//utils.WriteFile(file, s)
			t.BlockInfoCh <- blockInfo
			fmt.Println("1-加入管道No.\t", blockInfo.BlockNumber)
			time.Sleep(time.Second)
			needBlockNumber = blockInfo.BlockNumber + 1
			continue
		}

		// 如果不连续，则先把漏掉的block取出进行处理
		if blockInfo.BlockNumber > needBlockNumber {
			fmt.Println("开始调用接口：GetBlocksByHeightClient")
			// 根据blockNumber区间获取block信息
			// 经测试发现：区间超过100就会直接返回空，区间100以内，随着区间的大小，请求时间成正比
			// todo：如何在保证连续性的情况下，设置合理的区间大小？
			blockInfos, err := utils.GetBlocksByHeightClient(t.DSClient, needBlockNumber, blockInfo.BlockNumber)
			if err != nil {
				log.Println("502 获取block信息失败，", err)
				break
			}
			fmt.Println("GetBlocksByHeightClient调用成功")

			if len(blockInfos) == 0 {
				log.Println("503 获取block信息数量为0")
				break
			}

			for _, blockInfo = range blockInfos {
				//s := "22\t" + strconv.Itoa(int(blockInfo.BlockNumber)) + "\t" + blockInfo.BlockHash + "\n"
				//utils.WriteFile(file, s)
				t.BlockInfoCh <- blockInfo
				fmt.Println("2-加入管道No.\t", blockInfo.BlockNumber)
				time.Sleep(time.Second)
			}
			needBlockNumber = blockInfo.BlockNumber + 1
		}
	}
}

// 功能：获取交易信息，并写入指定文件
// 循环从管道取出block信息，调用http接口可以得到交易信息
// todo：错误中断循环后再次唤醒
// todo: 经大量测试，只有少数几个block有交易数据，那么无交易数据的block该如何处理？丢弃还是其他？
func GetTransferInfo(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("开始获取交易信息...")
	// 输出文件目录
	output := viper.GetString("global.output")
	for {
		blockInfo := <-t.BlockInfoCh
		if blockInfo == nil {
			log.Println("GetTransferInfo BlockInfoCh为空了")
			break
		}
		fmt.Println("从管道取出No.\t", blockInfo.BlockNumber)

		url := t.TransactionUrl + blockInfo.BlockHash
		resp, err := t.httpClient.Get(url)
		if err != nil {
			log.Println("GetTransferInfo 请求接口失败 ", err)
			break
		}
		fmt.Println("接口调用成功")

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("GetTransferInfo 读取body失败 ", err)
			break
		}
		fmt.Println("读取body成功")

		// 将得到的交易信息，解析到结构体，方便程序做下一步处理
		var transInfos = make([][]BlockTransferInfo, 0)
		err = json.Unmarshal(body, &transInfos)
		if err != nil {
			log.Println("GetTransferInfo body反序列化失败 ", err)
			break
		}
		fmt.Println("body反序列化成功")

		if len(transInfos) > 0 {
			for _, infos := range transInfos {
				for _, info := range infos {
					//s := fmt.Sprintf("blockHash: %s\nblockNumber: %d\nfrom: %s\nto: %s\namount: %d\ntimestamp: %d\nvalidAfterBlockNumber: %d\n\n",
					//	blockInfo.BlockHash, blockInfo.BlockNumber, info.FromAddr, info.ToAddr, info.Amount, info.Deploy.Timestamp, info.Deploy.ValidAfterBlockNumber)
					fmt.Println("开始写文件\n", info.ToStringByTab())
					utils.WriteFile(output+"block_transaction", info.ToStringByTab()+"\n", false)
					fmt.Println("写文件结束")
				}
			}
		} else {
			fmt.Println("无交易数据")
			// 将无交易数据的blockHash持久化到文件
			utils.WriteFile(output+"no_transaction_block",
				blockInfo.BlockHash+"\t"+strconv.FormatInt(blockInfo.BlockNumber, 10)+"\n", false)
			// todo：每个一段时间（比如10分钟）重新获取交易数据
			// todo：经测试，就算间隔半天，这部分block也无交易数据，对这部分block如何处理？
		}

		// 更新t.LastBlockNo为最新的blockNumber
		t.LastBlockNo = blockInfo.BlockNumber
		// 将blockInfo.BlockNumber更新到文件持久化，保证下次获取
		utils.WriteFile(output+".latest_blocknumber", strconv.FormatInt(blockInfo.BlockNumber, 10), true)
	}
}
