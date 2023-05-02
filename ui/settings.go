package ui

import (
	"image"
	"image/color"

	"github.com/miner/game"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

type Size string

func (s Size) undefined() bool {
	return s == ""
}
func (s Size) String() string {
	return string(s)
}

func (s Size) GridSize() int {
	return gridSizes[s]
}

type Difficulty string

func (d Difficulty) undefined() bool {
	return d == ""
}

func (d Difficulty) String() string {
	return string(d)
}

func (d Difficulty) ToInt() int {
	return difficulties[d]
}

const (
	sizeSmall  Size = "small"
	sizeMedium Size = "medium"
	sizeLarge  Size = "large"

	difficultyEasy   Difficulty = "easy"
	difficultyNormal Difficulty = "normal"
	difficultyHard   Difficulty = "hard"
)

var (
	gridSizes = map[Size]int{
		sizeSmall:  10,
		sizeMedium: 14,
		sizeLarge:  20,
	}
	cellSizes = map[int]int{
		10: cellSizeLarge,
		14: cellSizeMedium,
		20: cellSizeSmall,
	}
	difficulties = map[Difficulty]int{
		difficultyEasy:   10,
		difficultyNormal: 20,
		difficultyHard:   30,
	}
	sizeButtons       = make(map[Size]*sizeButton)
	difficultyButtons = make(map[Difficulty]*difficultyButton)

	green     = color.RGBA{178, 222, 39, 1}
	yellow    = color.RGBA{249, 215, 28, 1}
	red       = color.RGBA{236, 100, 75, 1}
	grey      = color.RGBA{128, 128, 128, 128}
	cyan      = color.RGBA{20, 205, 200, 1}
	black     = color.RGBA{0, 0, 0, 0}
	cellColor = color.RGBA{100, 255, 255, 255}
	fontColor = color.RGBA{255, 255, 255, 1}
)

func (c *Client) NewBackButton(ctx *scene.Context, p Position, s Shape, color, hoverColor color.RGBA, layer int) {
	sb := &backButton{}
	sb.id = ctx.Register(sb)
	sb.ColorBoxR = render.NewColorBoxR(int(s.width), int(s.height), color)
	sb.ColorBoxR.SetPos(p.x, p.y)

	sp := collision.NewSpace(p.x, p.y, s.width, s.height, sb.id)
	sp.SetZLayer(float64(layer))

	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)

	render.Draw(sb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, sb, func(sb *backButton, me *mouse.Event) event.Response {
		me.StopPropagation = true
		c.size = ""
		c.difficulty = ""
		c.grid = nil
		ctx.Window.GoToScene("settings")
		return 0
	})
	event.Bind(ctx, mouse.Start, sb, func(sb *backButton, me *mouse.Event) event.Response {
		sb.ColorBoxR.Color = image.NewUniform(hoverColor)
		me.StopPropagation = true
		return 0
	})
	event.Bind(ctx, mouse.Stop, sb, func(sb *backButton, me *mouse.Event) event.Response {
		sb.ColorBoxR.Color = image.NewUniform(color)
		me.StopPropagation = true
		return 0
	})
}

type Position struct {
	x, y float64
}
type Shape struct {
	width, height float64
}

func (c *Client) newSizeButton(ctx *scene.Context, p Position, s Shape, color, hoverColor color.RGBA, layer int, size Size, m map[Size]*sizeButton) {
	var text render.Renderable
	sb := &sizeButton{
		size:        size,
		sizeButtons: m,
		selected:    false,
	}
	sizeButtons[size] = sb
	sb.id = ctx.Register(sb)
	sb.ColorBoxR = render.NewColorBoxR(int(s.width), int(s.height), color)
	sb.ColorBoxR.SetPos(p.x, p.y)
	sp := collision.NewSpace(p.x, p.y, s.width, s.height, sb.id)
	sp.SetZLayer(float64(layer))
	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)
	render.Draw(sb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, sb, func(sb *sizeButton, me *mouse.Event) event.Response {
		me.StopPropagation = true
		if sb.selected {
			return 0
		}
		sb.ShiftX(-20)
		sb.selected = true
		c.size = size
		for s, button := range sb.sizeButtons {
			if s != sb.size {
				if button.selected {
					button.ShiftX(20)
					button.selected = false
				}
			}
		}
		return 0
	})
	event.Bind(ctx, mouse.Start, sb, func(sb *sizeButton, me *mouse.Event) event.Response {
		sb.ColorBoxR.Color = image.NewUniform(hoverColor)
		me.StopPropagation = true
		text, _ = render.Draw(c.font.NewText(size.String(), p.x+s.width/2-20, p.y+s.height/2-10))
		return 0
	})
	event.Bind(ctx, mouse.Stop, sb, func(sb *sizeButton, me *mouse.Event) event.Response {
		sb.ColorBoxR.Color = image.NewUniform(color)
		me.StopPropagation = true
		text.Undraw()
		return 0
	})
}

func (c *Client) newStartButton(ctx *scene.Context, p Position, s Shape, color, hoverColor color.RGBA, layer int) {
	var text render.Renderable
	hb := &startButton{
		button{color: color, hoverColor: hoverColor},
	}
	hb.id = ctx.Register(hb)
	hb.ColorBoxR = render.NewColorBoxR(int(s.width), int(s.height), color)
	hb.ColorBoxR.SetPos(p.x, p.y)
	sp := collision.NewSpace(p.x, p.y, s.width, s.height, hb.id)
	sp.SetZLayer(float64(layer))
	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)
	render.Draw(hb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, hb, func(box *startButton, me *mouse.Event) event.Response {
		me.StopPropagation = true
		if c.difficulty.undefined() || c.size.undefined() {
			ctx.Window.GoToScene("error")
			return 0
		}
		ctx.Window.NextScene()
		return 0
	})
	event.Bind(ctx, mouse.Start, hb, func(box *startButton, me *mouse.Event) event.Response {
		box.ColorBoxR.Color = image.NewUniform(hoverColor)
		me.StopPropagation = true
		text, _ = render.Draw(c.font.NewText("Start", p.x+s.width/2-20, p.y+s.height/2-10))
		return 0
	})
	event.Bind(ctx, mouse.Stop, hb, func(box *startButton, me *mouse.Event) event.Response {
		box.ColorBoxR.Color = image.NewUniform(color)
		me.StopPropagation = true
		text.Undraw()
		return 0
	})
}

func (c *Client) newDifficultyButton(ctx *scene.Context, p Position, s Shape, color, hoverColor color.RGBA, layer int, diff Difficulty, m map[Difficulty]*difficultyButton) {
	var text render.Renderable
	sb := &difficultyButton{
		button: button{
			color:      color,
			hoverColor: hoverColor,
		},
		difficulty:        diff,
		difficultyButtons: m,
		selected:          false,
	}
	difficultyButtons[diff] = sb
	sb.id = ctx.Register(sb)
	sb.ColorBoxR = render.NewColorBoxR(int(s.width), int(s.height), color)
	sb.ColorBoxR.SetPos(p.x, p.y)

	sp := collision.NewSpace(p.x, p.y, s.width, s.height, sb.id)
	sp.SetZLayer(float64(layer))

	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)

	render.Draw(sb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, sb, func(sb *difficultyButton, me *mouse.Event) event.Response {
		me.StopPropagation = true
		if sb.selected {
			return 0
		}
		sb.ShiftX(20)
		sb.selected = true
		c.difficulty = diff
		for d, button := range sb.difficultyButtons {
			if d != sb.difficulty {
				if button.selected {
					button.ShiftX(-20)
					button.selected = false
				}
			}
		}
		return 0
	})
	event.Bind(ctx, mouse.Start, sb, func(sb *difficultyButton, me *mouse.Event) event.Response {
		sb.ColorBoxR.Color = image.NewUniform(sb.hoverColor)
		me.StopPropagation = true
		text, _ = render.Draw(c.font.NewText(diff.String(), p.x+s.width/2-20, p.y+s.height/2-10))
		return 0
	})
	event.Bind(ctx, mouse.Stop, sb, func(sb *difficultyButton, me *mouse.Event) event.Response {

		sb.ColorBoxR.Color = image.NewUniform(sb.color)
		me.StopPropagation = true
		text.Undraw()
		return 0
	})
}

func (c *Client) NewErrorScene() scene.Scene {
	return scene.Scene{Start: func(ctx *scene.Context) {
		ctx.DrawStack.Draw(c.font.NewText("Bad input!", 210, 240))
		c.NewBackButton(ctx, Position{0, 0}, Shape{20, 480}, cyan, grey, 1)
	}}

}

//func (c *Client) NewWinScene() scene.Scene {
//	return scene.Scene{Start: func(ctx *scene.Context) {
//		ctx.DrawStack.Draw(c.font.NewText("CONGRATULATIONS!", 250, 240))
//		c.NewBackButton(ctx, Position{0, 0}, Shape{20, 480}, cyan, grey, 1)
//	}}
//}

//func (c *Client) NewLoseScene() scene.Scene {
//	return scene.Scene{Start: func(ctx *scene.Context) {
//		ctx.DrawStack.Draw(c.font.NewText("YOU LOSE!", 250, 240))
//		c.NewBackButton(ctx, Position{0, 0}, Shape{20, 480}, cyan, grey, 1)
//	}}
//}

func (c *Client) newSettingScene() scene.Scene {
	return scene.Scene{
		Start: func(ctx *scene.Context) {
			s := Shape{200, 50}
			c.newSizeButton(ctx, Position{119, 50}, s, green, grey, 1, sizeSmall, sizeButtons)
			c.newSizeButton(ctx, Position{119, 102}, s, yellow, grey, 1, sizeMedium, sizeButtons)
			c.newSizeButton(ctx, Position{119, 154}, s, red, grey, 1, sizeLarge, sizeButtons)

			c.newDifficultyButton(ctx, Position{321, 50}, s, green, grey, 1, difficultyEasy, difficultyButtons)
			c.newDifficultyButton(ctx, Position{321, 102}, s, yellow, grey, 1, difficultyNormal, difficultyButtons)
			c.newDifficultyButton(ctx, Position{321, 154}, s, red, grey, 1, difficultyHard, difficultyButtons)

			c.newStartButton(ctx, Position{119, 206}, Shape{402, 50}, cyan, grey, 1)
		},
		End: func() (string, *scene.Result) {
			g := game.NewGame()
			c.game = g
			c.window.AddScene("game", c.newGameScene())
			return "game", nil //set the next scene to "game"
		},
	}
}
