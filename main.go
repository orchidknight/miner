package main

import (
	"github.com/miner/logger"
	"github.com/miner/ui"
)

func main() {
	log := logger.NewLog()
	ui.Run(log)
}
