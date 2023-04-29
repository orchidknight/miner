package ui

import (
	"github.com/miner/game"
	"github.com/miner/logger"
	"github.com/oakmound/oak/v4"
)

func Run(log logger.Logger) {
	g := game.NewGame()
	font, err := newFont()
	if err != nil {
		log.Error("ui", "NewFont: %v", err)
	}
	oak.AddScene("settings", NewSettingsScene(g, font, log))
	oak.AddScene("error", NewErrorScene(g, font, log))
	oak.AddScene("lose", NewLoseScene(g, font, log))
	oak.AddScene("win", NewWinScene(g, font, log))
	oak.Init("settings")
}
