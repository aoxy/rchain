package bitcoin

import (
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// 输出主要包含两部分：
// 一定量的比特币(Value)
// 一个锁定脚本(ScriptPubKey)，要花这笔钱，必须要解锁该脚本。
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// 比特币采用的是 UTXO 模型，并非账户模型，
// 并不直接存在“余额”这个概念，余额需要通过遍历整个交易历史得来。
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// 功能：给其他人发送一些币
// 创建一笔新的交易，将它放到一个块里，然后挖出这个块。
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var ins []TXInput
	var outs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds.")
	}

	// 创建input列表
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, out := range outs {
			input := TXInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			ins = append(ins, input)
		}
	}

	outs = append(outs, TXOutput{amount, to})
	if acc > amount {
		outs = append(outs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, ins, outs}
	tx.SetID()
	return &tx
}

// coinbase 交易: 是一种特殊的交易, 它不需要引用之前一笔交易的输出。
// 它“凭空”产生了币（也就是产生了新币），这是矿工获得挖出新块的奖励，也可以理解为“发行新币”。
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

// todo
func (tx *Transaction) IsCoinBase() bool {
	return tx.Vin[0].Vout == -1
}

// todo
func (tx *Transaction) SetID() {
	tx.ID = []byte{}
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
