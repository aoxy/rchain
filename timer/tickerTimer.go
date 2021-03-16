package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	closeCh := make(chan struct{})
	timer := time.NewTicker(time.Second)
	showBlockInfo := func() error {
		fmt.Println("blockInfo")
		return nil
	}

	if err := showBlockInfo(); err != nil {
		log.Fatal("showBlockInfo error ", err)
	}

	for {
		select {
		case <-timer.C:
			if err := showBlockInfo(); err != nil {
				log.Fatal("showBlockInfo error ", err)
				panic(err)
			}
		case <-closeCh:
			timer.Stop()
			return
		}
	}
}
