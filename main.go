package main

import (
	"rchain/bitcoin"
	"rchain/config"
)

func main() {

	config.LoadConfig()
	// service.Run()
	cli := bitcoin.CLI{}
	cli.Run()
}
