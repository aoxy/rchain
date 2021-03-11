package service

import (
	"dawn1806/rchain/internal/utils"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"sync"
	"time"
)

var grpcUrl string
var httpUrl string

func InitConfig() {
	hosts := viper.GetStringSlice("observer.hosts")
	grpcPort := viper.GetString("observer.grpcPort")
	httpPort := viper.GetString("observer.httpPort")
	grpcUrl = hosts[0] + ":" + grpcPort
	httpUrl = "http://" + hosts[3] + ":" + httpPort + "/getTransaction/"
}

func Run() {
	DSClient, conn, err := utils.NewDeployServiceClient(grpcUrl)
	if err != nil {
		fmt.Println("NewDeployServiceClient error ", err)
		return
	}
	defer conn.Close()

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	t := &Transaction{
		DSClient:       DSClient,
		httpClient:     client,
		TransactionUrl: httpUrl,
		BlockHash:      make(chan string, 1),
	}
	HandleTransaction(t)
}

func HandleTransaction(t *Transaction) {

	wg := sync.WaitGroup{}
	wg.Add(2)

	// 启动一个协程，获取最新的blockHash，放入管道
	go func() {
		for {
			if err := t.GetBlockHash(); err != nil {
				break
			}
			time.Sleep(time.Second)
		}
		close(t.BlockHash)
		wg.Done()
	}()

	// 从管道取出blockHash，调用http接口获得交易信息，将交易信息写入文件
	go func() {
		for {
			if err := t.GetTransferInfo(); err != nil {
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func init() {
	viper.SetConfigFile("./config/rchain-dev.toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("viper read file error ", err)
		return
	}
	fmt.Fprintln(os.Stdout, "using config file:", viper.ConfigFileUsed())
}
