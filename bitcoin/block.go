package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// 创世区块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block.", []byte{})
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	// 工作量证明计算hash
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 序列化block，用于持久化存储
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		panic(err)
	}
	return result.Bytes()
}

// 反序列化block
func DeserializeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&block); err != nil {
		panic(err)
	}
	return &block
}

// 弃用，改用工作量证明计算hash
// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
// 	hash := sha256.Sum256(headers)
// 	b.Hash = hash[:]
// }
