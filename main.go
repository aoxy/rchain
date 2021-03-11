package main

import "dawn1806/rchain/internal/service"

func main() {
	service.InitConfig()
	service.Run()
}
