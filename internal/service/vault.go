package service

import (
	"bytes"
	casperV1 "dawn1806/rchain/api/pb"
	"dawn1806/rchain/internal/utils"
	"fmt"
	"html/template"
	"time"
)

type Vault struct {
	FromAddr  string
	ToAddr    string
	Amount    int
	Deployer  []byte // public key
	PhloPrice int64
	PhloLimit int64
}

func Transfer() {
	vault := &Vault{
		FromAddr:  "1111MBwE8Mt3Lz5aidYdgqREd1HnuWVJj6grkJkgW9UYCGZ8uWy6Z",
		ToAddr:    "11112iTYv5p32nyquVFaWqp5XnoNKkvc5oJZZh6KT6xApGE3oPLFPs",
		Amount:    140000000000,
		Deployer:  []byte("public key"),
		PhloPrice: 1,
		PhloLimit: 500000,
	}

	// 读取模板文件，将变量解析进去
	// The difference between `transfer_ensure` and `transfer` is that , if the to_addr is not created in the
	// chain, the `transfer` would hang until the to_addr successfully created in the change and the `transfer_ensure`
	// can be sure that if the `to_addr` is not existed in the chain the process would created the vault in the chain
	// and make the transfer successfully.
	tmpl, _ := template.ParseFiles("E:\\work\\rchain\\templates\\transferRhoTpl")
	// tmpl2, _ := template.ParseFiles("E:\\work\\rchain\\templates\\transferEnsureToRhoTpl")
	var tmplBytes bytes.Buffer
	tmpl.Execute(&tmplBytes, vault)
	fmt.Println(tmplBytes.String())

	term := tmplBytes.String()

	in := casperV1.DeployDataProto{
		Deployer:              vault.Deployer,           // public key
		Term:                  term,                     // rholang源代码
		Timestamp:             time.Now().Unix() * 1000, // millisecond timestamp
		SigAlgorithm:          "secp256k1",              // 签名算法名称
		PhloPrice:             vault.PhloPrice,
		PhloLimit:             vault.PhloLimit,
		ValidAfterBlockNumber: -1,
	}
	// todo: in.Sig = sign(key, in)

	utils.DoDeploy(t.DSClient, &in)
}
