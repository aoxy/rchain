package test

import (
	"dawn1806/rchain/internal/utils"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBlockInfo(t *testing.T) {
	grpcUrl := "observer-asia.services.mainnet.rchain.coop:40401"
	DSClient, _, _ := utils.NewDeployServiceClient(grpcUrl)
	blockInfo, _ := utils.LastFinalizedBlock(DSClient)
	fmt.Println(blockInfo)
}

func TestDeployService(t *testing.T) {

	var transCh = make(chan string, 1)

	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 3; i++ {
			s := "111"
			fmt.Println(s)
			transCh <- s
			time.Sleep(time.Second)
		}
		close(transCh)
		wg.Done()
	}()

	go func() {
		for {
			s := <-transCh
			if s == "" {
				break
			}
			s = s + "222"
			fmt.Println(s)
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("over.")
}
