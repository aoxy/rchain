package bitcoin

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) Run() {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	from := addBlockCmd.String("from", "", "Transaction data")
	to := addBlockCmd.String("to", "", "Transaction data")
	amount := addBlockCmd.Int("amount", 0, "Transaction data")

	switch os.Args[1] {
	case "addblock":
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	default:
		// cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *from == "" || *to == "" || *amount == 0 {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.send(*from, *to, *amount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) send(from, to string, amount int) {
	fmt.Printf("开始创建区块链...from: %s\tto: %s\tamount: %d\n", from, to, amount)
	bc := NewBlockchain(from)
	fmt.Println("链成功创建：", bc)
	defer bc.Db.Close()

	cli.Bc = bc

	fmt.Println("开始进行交易...")
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
	fmt.Println("add block success!")
}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
