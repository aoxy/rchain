package main

import (
	"dawn1806/rchain/bitcoin"
	"dawn1806/rchain/config"
)

func main() {

	config.LoadConfig()

	// service.Run()

	bc := bitcoin.NewBlockchain()
	defer bc.Db.Close()

	cli := bitcoin.CLI{Bc: bc}
	cli.Run()
}
