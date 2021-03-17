package test

import (
	"fmt"
	"rchain/internal/service"
	"rchain/internal/utils"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetBlocksByHeightClient(t *testing.T) {
	// numbers: 651551	651552	651684
	grpcUrl := "observer-asia.services.mainnet.rchain.coop:40401"
	DSClient, _, _ := utils.NewDeployServiceClient(grpcUrl)
	// 区间超过100就会直接返回空
	// 区间100以内，随着区间的大小，请求时间成正比
	blockInfos, err := utils.GetBlocksByHeightClient(DSClient, 6510, 6515)
	fmt.Println("err: ", err)
	fmt.Println("blockInfos: ", blockInfos)
}

func TestStructFormat(t *testing.T) {
	for i := 0; i < 1; i++ {
		b := service.BlockTransferInfo{
			FromAddr:      "fromaddr123456",
			ToAddr:        "toaddr123456",
			Amount:        560000,
			RetUnforeable: "abcd",
			Deploy: struct {
				Deployer              string `json:"deployer"`
				Term                  string `json:"term"`
				Timestamp             int64  `json:"timestamp"`
				Sig                   string `json:"sig"`
				SigAlgorithm          string `json:"sigAlgorithm"`
				PhloPrice             int64  `json:"phloPrice"`
				PhloLimit             int64  `json:"phloLimit"`
				ValidAfterBlockNumber int64  `json:"validAfterBlockNumber"`
				Cost                  int64  `json:"cost"`
				Errored               bool   `json:"errored"`
				SystemDeployError     string `json:"systemDeployError"`
			}{
				Deployer:              "deployer1111",
				Term:                  "TerTermTermTermTermTermTermTermTermTermTermTermTermTermTermTermTerm",
				Timestamp:             124235423534,
				Sig:                   "a",
				SigAlgorithm:          "b",
				PhloPrice:             10,
				PhloLimit:             10,
				ValidAfterBlockNumber: 10,
				Cost:              10,
				Errored:           false,
				SystemDeployError: "mark",
			},
			Success: false,
			Reason:  "",
		}

		var v = reflect.ValueOf(b)
		var s = make([]string, 0)
		fmt.Println("s: ", s)
		s = formatStruct(v, s)
		fmt.Println("s: ", s)

		//formatDeploy := func(v reflect.Value) string {
		//	t1 := v.Type()
		//	s1 := make([]string, v.NumField())
		//	for i := range s1 {
		//		if t1.Field(i).Name == "Term" {
		//			fmt.Println("Term: ", v.Field(i))
		//			continue
		//		}
		//		s1[i] = fmt.Sprintf("%v", v.Field(i))
		//	}
		//	return strings.Join(s1, "\t")
		//}

		//v := reflect.ValueOf(b)
		//t := reflect.TypeOf(b)
		//s := make([]string, v.NumField())
		//for i := range s {
		//	if t.Field(i).Name == "Deploy" {
		//		tempstr := formatDeploy(v.Field(i))
		//		s[i] = tempstr
		//		continue
		//	}
		//	s[i] = fmt.Sprintf("%v", v.Field(i))
		//}

		time.Sleep(time.Second)
		utils.WriteFile("E:\\work\\rchain\\docs\\output\\test", strings.Join(s, "\t")+"\n", false)
	}
}

func formatStruct(v reflect.Value, s []string) []string {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == "Deploy" {
			s = formatStruct(v.Field(i), s)
			continue
		}
		if t.Field(i).Name == "Term" {
			continue
		}
		s = append(s, fmt.Sprintf("%v", v.Field(i)))
		fmt.Println("sss: ", s)
	}
	return s
}

func TestGetBlockInfos(t *testing.T) {
	grpcUrl := "observer-asia.services.mainnet.rchain.coop:40401"
	DSClient, _, _ := utils.NewDeployServiceClient(grpcUrl)
	blockInfo, _ := utils.LastFinalizedBlock(DSClient)
	fmt.Println(blockInfo)
}
