package main

import (
	"github.com/miner/logger"
	"github.com/miner/ui"
)

func main() {
	log := logger.NewLog()
	client := ui.NewClient(log)
	err := client.Run()
	if err != nil {
		log.Fatal("client", "Run : %v", err)
	}
}
