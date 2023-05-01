package game

type Game interface {
	Reveal(x, y int) ([]Cell, GameState, error)
	Start(size, difficulty int) error
}

type GameState string

const (
	Lose       GameState = "lose"
	Win        GameState = "win"
	InProgress GameState = "in progress"
)
