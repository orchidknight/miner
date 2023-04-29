package game

type Game interface {
	Reveal(x, y int) ([]Cell, string)
	Reset()
	Start()
	GetCell() Cell
	SetSize(size int) error
	SetDifficulty(diff int) error
}
