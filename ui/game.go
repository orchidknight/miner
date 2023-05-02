package ui

import (
	"fmt"
	"image"
	"image/color"

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

func (c *Client) newGameScene() scene.Scene {
	fmt.Println("new game scene")
	size := c.size.GridSize()
	fmt.Println(size, c.difficulty.ToInt())
	err := c.game.Start(size, c.difficulty.ToInt())
	if err != nil {
		c.log.Error("game", "Start: %v", err)
	}
	grid := &Grid{size: size,
		cellMap: make(map[int]*cellButton, size*size)}
	c.grid = grid

	s := scene.Scene{
		Start: func(ctx *scene.Context) {
			size := c.size.GridSize()
			grid := c.grid
			cellSize := cellSizes[size]
			offset := calcOffset(size, cellSize)
			c.NewBackButton(ctx, Position{0, 0}, Shape{20, 480}, cyan, grey, 1)
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					x := offset.x + float64(i*cellSize)
					y := offset.y + float64(j*cellSize)
					grid.cellMap[i*size+j] = c.newCellButton(ctx, i, j, float64(x), float64(y), float64(cellSize-1), float64(cellSize-1), cellColor, grey, 3)

				}
			}
		}}
	return s
}

func (c *Client) newCellButton(ctx *scene.Context, ix, iy int, x, y, w, h float64, clr, hclr color.RGBA, layer int) *cellButton {
	hb := &cellButton{
		button: button{
			color:      clr,
			hoverColor: hclr,
		},
		x:        ix,
		y:        iy,
		Position: Position{x, y},
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
		size := c.size.GridSize()
		if box.revealed {
			return 0
		}
		box.ColorBoxR.Color = image.NewUniform(color.RGBA{128, 128, 128, 128})
		me.StopPropagation = true
		cells, state, err := c.game.Reveal(hb.x, hb.y)
		if err != nil {
			c.log.Error("game", "Reveal: %v", err)
		}
		if state == loseState {
			for _, cell := range cells {
				cb := c.grid.cellMap[cell.X()*size+cell.Y()]
				cb.revealed = true
				if cell.HasBomb() {
					cb.ColorBoxR.Color = image.NewUniform(red)
					continue
				}
				cb.ColorBoxR.Color = image.NewUniform(grey)
				if cell.Count() > 0 {
					render.Draw(render.NewText(fmt.Sprintf("%d", cell.Count()), cb.Position.x+float64(cellSizes[size]/2-5), cb.Position.y+float64(cellSizes[size]/2-9)))
				}
			}
			ctx.DrawStack.Draw(c.font.NewText("YOU LOSE!", 250, 15))
			return 0
		}
		if state == winState {
			for _, cell := range cells {
				cb := c.grid.cellMap[cell.X()*size+cell.Y()]
				cb.revealed = true
				if cell.HasBomb() {
					cb.ColorBoxR.Color = image.NewUniform(red)
					continue
				}
				cb.ColorBoxR.Color = image.NewUniform(grey)
				if cell.Count() > 0 {
					render.Draw(render.NewText(fmt.Sprintf("%d", cell.Count()), cb.Position.x+float64(cellSizes[size]/2), cb.Position.y+float64(cellSizes[size]/2)))
				}
			}
			ctx.DrawStack.Draw(c.font.NewText("CONGRATULATIONS!", 250, 15))
			return 0
		} else {
			for _, cell := range cells {
				cb := c.grid.cellMap[cell.X()*size+cell.Y()]
				cb.revealed = true
				cb.ColorBoxR.Color = image.NewUniform(black)
				if cell.Count() > 0 {
					render.Draw(render.NewText(fmt.Sprintf("%d", cell.Count()), cb.Position.x+float64(cellSizes[size]/2-5), cb.Position.y+float64(cellSizes[size]/2-9)))
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
