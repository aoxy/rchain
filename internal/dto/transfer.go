package dto

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

func (b BlockTransferInfo) String() string {
	// todo
	return ""
}
