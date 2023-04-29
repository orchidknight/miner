package ui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/miner/game"
	"github.com/miner/logger"
	"github.com/oakmound/oak/v4/collision"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	cellSizeSmall  = 20
	cellSizeMedium = 25
	cellSizeLarge  = 30
	loseState      = "lose"
	winState       = "win"
	windowHeight   = 480
	windowWidth    = 640
)

type Grid struct {
	size    int
	cells   [][]*cellButton
	cellMap map[int]*cellButton
}

func calcOffset(gridSize int, cellSize int) Position {
	return Position{
		0.5 * (windowWidth - float64(gridSize*cellSize)),
		0.5 * (windowHeight - float64(gridSize*cellSize)),
	}
}

func newGameScene(game *game.Miner, log logger.Logger) scene.Scene {
	err := game.Start()
	if err != nil {
		log.Error("game", "Start: %v", err)
	}
	grid := &Grid{size: game.Size,
		cellMap: make(map[int]*cellButton, game.Size*game.Size)}

	s := scene.Scene{
		Start: func(ctx *scene.Context) {
			cellSize := cellSizes[game.Size]
			offset := calcOffset(game.Size, cellSize)
			NewBackButton(ctx, Position{0, 0}, Shape{20, 480}, cyan, grey, 1, game)
			for i := 0; i < game.Size; i++ {
				for j := 0; j < game.Size; j++ {
					x := offset.x + float64(i*cellSize)
					y := offset.y + float64(j*cellSize)
					cell := game.Grid.GetCell(i, j)
					grid.cellMap[cell.X()*game.Size+cell.Y()] = newCellButton(ctx, cell.X(), cell.Y(), cell.Count(), float64(x), float64(y), float64(cellSize-1), float64(cellSize-1), cellColor, grey, 3, game, grid)

				}
			}
		}}
	return s
}

func newCellButton(ctx *scene.Context, ix, iy int, count int, x, y, w, h float64, clr, hclr color.RGBA, layer int, game *game.Miner, grid *Grid) *cellButton {
	hb := &cellButton{
		button: button{
			color:      clr,
			hoverColor: hclr,
		},
		x:        ix,
		y:        iy,
		Position: Position{x, y},
		count:    count,
	}
	hb.id = ctx.Register(hb)
	hb.ColorBoxR = render.NewColorBoxR(int(w), int(h), clr)
	hb.ColorBoxR.SetPos(x, y)

	sp := collision.NewSpace(x, y, w, h, hb.id)
	sp.SetZLayer(float64(layer))

	mouse.Add(sp)
	mouse.PhaseCollision(sp, ctx.Handler)

	render.Draw(hb.ColorBoxR, layer)

	event.Bind(ctx, mouse.ClickOn, hb, func(box *cellButton, me *mouse.Event) event.Response {
		if box.revealed {
			return 0
		}
		box.ColorBoxR.Color = image.NewUniform(color.RGBA{128, 128, 128, 128})
		me.StopPropagation = true
		cells, state := game.Reveal(hb.x, hb.y)
		if state == loseState {
			ctx.Window.GoToScene(loseState)
			return 0
		}
		if state == winState {
			ctx.Window.GoToScene(winState)
			return 0
		} else {
			for _, cell := range cells {
				cb := grid.cellMap[cell.X()*game.Size+cell.Y()]
				cb.revealed = true
				cb.ColorBoxR.Color = image.NewUniform(black)
				if cb.count > 0 {
					render.Draw(render.NewText(fmt.Sprintf("%d", cb.count), cb.Position.x+float64(cellSizes[game.Size]/2)-2, cb.Position.y+float64(cellSizes[game.Size]/2)-2))
				}
			}
		}
		return 0
	})
	event.Bind(ctx, mouse.Start, hb, func(box *cellButton, me *mouse.Event) event.Response {
		if box.revealed {
			return 0
		}
		box.ColorBoxR.Color = image.NewUniform(box.hoverColor)
		me.StopPropagation = true
		return 0
	})
	event.Bind(ctx, mouse.Stop, hb, func(box *cellButton, me *mouse.Event) event.Response {
		if box.revealed {
			return 0
		}
		box.ColorBoxR.Color = image.NewUniform(clr)
		me.StopPropagation = true
		return 0
	})
	return hb
}
