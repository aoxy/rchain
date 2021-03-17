package service

import (
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	casperV1 "rchain/api/pb"
	"rchain/internal/utils"
	"reflect"
	"strings"
)

type Transaction struct {
	DSClient       casperV1.DeployServiceClient
	httpClient     *http.Client
	Conn           *grpc.ClientConn
	TransactionUrl string
	BlockInfoCh    chan *casperV1.LightBlockInfo
	LastBlockNo    int64
	ttl            int
}

type BlockTransferInfo struct {
	FromAddr      string `json:"fromAddr"`
	ToAddr        string `json:"toAddr"`
	Amount        int64  `json:"amount"`
	RetUnforeable string `json:"retUnforeable"`
	Deploy        struct {
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
	} `json:"deploy"`
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (b BlockTransferInfo) ToStringByTab() string {
	return strings.Join(b.ToSlice(), "\t")
}

func (b BlockTransferInfo) ToSlice() []string {
	v := reflect.ValueOf(b)
	s := make([]string, 0)
	s = utils.FormatStruct(v, s)
	return s
}

func (b BlockTransferInfo) ToSlice1() []string {
	v := reflect.ValueOf(b)
	s := make([]string, v.NumField())
	for i := range s {
		s[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return s
}
