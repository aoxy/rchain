package test

import (
	"fmt"
	"testing"
)

func TestBlockchain(t *testing.T) {
	// b := bitcoin.NewBlockchain()
	// b.AddBlock("send 1BTC to a")
	// b.AddBlock("send 0.5BTC to b")

	// 对block进行验证
	// for _, block := range b.Blocks {
	// 	fmt.Printf("Prev: %s\n", block.PrevBlockHash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	pow := bitcoin.NewProofOfWork(block)
	// 	fmt.Printf("Pow: %s\n\n", strconv.FormatBool(pow.Validate()))
	// }
}

func TestHash(t *testing.T) {
	// s := "abc"
	// b := []byte(s)
	b := []byte{97, 98, 99}
	fmt.Printf("hash: %s\n", b)
}

func TestWork(t *testing.T) {
	// var n int64 = 1
	// a := big.NewInt(n)
	// fmt.Println(a)
	// a.Lsh(a, 2)
	// fmt.Println(bitcoin.IntToHex(a.Int64()))
	// fmt.Println(a)
}
