package ui

import (
	"github.com/miner/game"
	"github.com/miner/logger"
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/render"
)

type Client struct {
	size       Size
	difficulty Difficulty
	game       game.Game
	window     *oak.Window
	font       *render.Font
	log        logger.Logger
	grid       *Grid
}

func NewClient(log logger.Logger) *Client {
	window := oak.NewWindow()
	font, err := newFont()
	if err != nil {
		log.Error("ui", "NewFont: %v", err)
	}
	return &Client{
		window: window,
		font:   font,
		log:    log,
	}
}

func (c *Client) Run() error {
	var err error
	c.game = game.NewGame()

	err = c.window.AddScene("settings", c.newSettingScene())
	if err != nil {
		return err
	}
	err = c.window.AddScene("error", c.NewErrorScene())
	if err != nil {
		return err
	}
	err = c.window.AddScene("lose", c.NewLoseScene())
	if err != nil {
		return err
	}
	err = c.window.AddScene("win", c.NewWinScene())
	if err != nil {
		return err
	}
	err = c.window.Init("settings")
	if err != nil {
		return err
	}
	return nil
}
