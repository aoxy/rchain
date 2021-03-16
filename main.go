package main

import (
	"dawn1806/rchain/bitcoin"
	"dawn1806/rchain/config"
)

func main() {

	config.LoadConfig()

	// service.Run()

	cli := bitcoin.CLI{}
	cli.Run()
}
