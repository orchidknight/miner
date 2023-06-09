package game

import (
	"errors"
	"math/rand"
	"time"
)

var ErrInvalidPosition = errors.New("invalid x y position")
var ErrInvalidSettings = errors.New("invalid size or difficulty")

type Miner struct {
	Size          int
	Difficulty    int
	BombsCount    int
	RevealedCount int
	Bombs         map[int]Position
	Grid          *Grid
}

func NewGame() *Miner {
	return &Miner{
		Bombs: make(map[int]Position),
	}
}

func (g *Miner) cells() []Cell {
	all := make([]Cell, 0, g.Size*g.Size)
	for _, cells := range g.Grid.cells {
		all = append(all, cells...)
	}
	return all
}

// Reveal checks the given cell with incoming coordinates
// If bomb - returns all cells for revealing, game state - lose.
// If cell is empty, recursively collects all adjacent empty cells to reveal, game state is in progress. Recursive cell traversal ends if cell bomb count is greater than 0
// If all possible cells are revealed, the game state is win, returns all cells to be revealed.
func (g *Miner) Reveal(x, y int) ([]Cell, GameState, error) {
	if !g.Grid.validatedPosition(x, y) {
		return nil, InProgress, ErrInvalidPosition
	}
	revealedCells := make([]Cell, 0)
	if g.Grid.getCell(x, y).HasBomb() {
		return g.cells(), Lose, nil
	}
	revealed := make(map[int]Position)
	g.check(x, y, revealed)
	for k, p := range revealed {
		revealedCells = append(revealedCells, g.Grid.getCell(p.x, p.y))
		g.Grid.revealed[k] = p
	}
	if g.Size*g.Size-len(g.Grid.revealed) == g.BombsCount {
		return revealedCells, Win, nil
	}
	return revealedCells, InProgress, nil
}

// Start initiate the game with the given settings. Cannot be created if the settings are null. Cell matrix with uniform distribution is created
func (g *Miner) Start(size, difficulty int) error {
	if size <= 0 || difficulty <= 0 {
		return ErrInvalidSettings
	}
	g.Size, g.Difficulty = size, difficulty
	g.BombsCount = (g.Size * g.Size * g.Difficulty) / 100
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < g.BombsCount; i++ {
		b := rand.Intn(g.Size * g.Size)
		if _, ok := g.Bombs[b]; !ok {
			g.Bombs[b] = Position{b / g.Size, b % g.Size}
		} else {
			i--
		}
	}
	grid := &Grid{
		cells:    make([][]Cell, g.Size),
		revealed: make(map[int]Position),
	}
	for x := range grid.cells {
		grid.cells[x] = make([]Cell, g.Size)
		for y := range grid.cells[x] {
			bomb := false
			if _, ok := g.Bombs[x*g.Size+y]; ok {
				bomb = true
			}
			grid.cells[x][y] = Cell{
				Position: Position{x, y},
				revealed: false,
				bomb:     bomb,
				count:    0,
			}

		}
	}
	for x := range grid.cells {
		for y := range grid.cells[x] {
			count := 0
			cellNumbers := grid.nearCells(x, y)
			for _, cellNumber := range cellNumbers {
				if grid.cells[cellNumber/g.Size][cellNumber%g.Size].HasBomb() {
					count++
				}
			}
			grid.cells[x][y].count = count
		}
	}
	g.Grid = grid
	return nil
}

type Position struct {
	x, y int
}

type Grid struct {
	cells    [][]Cell
	revealed map[int]Position
}

func (g *Grid) getCell(x, y int) Cell {
	return g.cells[x][y]
}

type Cell struct {
	Position
	revealed bool
	bomb     bool
	count    int
}

func (c Cell) X() int {
	return c.x
}
func (c Cell) Y() int {
	return c.y
}

func (c Cell) Count() int {
	return c.count
}

func (g *Grid) validatedPosition(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if x >= len(g.cells) || y >= len(g.cells) {
		return false
	}
	return true

}

func (g *Grid) nearCells(x, y int) []int {
	var i, j int
	size := len(g.cells)
	cellNumbers := make([]int, 0)
	i, j = x-1, y-1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x-1, y
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x-1, y+1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x, y-1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x, y+1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x+1, y-1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x+1, y
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	i, j = x+1, y+1
	if g.validatedPosition(i, j) {
		cellNumbers = append(cellNumbers, i*size+j)
	}
	return cellNumbers
}

func (g *Miner) check(x, y int, revealed map[int]Position) {
	cell := g.Grid.getCell(x, y)
	revealed[x*g.Size+y] = Position{x, y}
	if cell.count > 0 {
		return
	}

	for _, position := range g.Grid.nearCells(x, y) {
		if _, ok := revealed[position]; ok {
			continue
		}
		g.check(position/g.Size, position%g.Size, revealed)
	}
}

func (c Cell) HasBomb() bool {
	return c.bomb
}
