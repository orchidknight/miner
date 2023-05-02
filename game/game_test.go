package game

import (
	"errors"
	"testing"
)

func TestMiner_Reveal(t *testing.T) {
	tests := map[string]struct {
		size          int
		difficulty    int
		x, y          int
		expectedErr   error
		expectedState GameState
	}{
		"win check":  {size: 10, difficulty: 1, x: 0, y: 0, expectedErr: nil, expectedState: Win},
		"lose check": {size: 10, difficulty: 100, x: 0, y: 0, expectedErr: nil, expectedState: Lose},
		"wrong x y":  {size: 10, difficulty: 1, x: -1, y: 0, expectedErr: ErrInvalidPosition, expectedState: InProgress},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			game := NewGame()
			game.Start(tc.size, tc.difficulty)
			_, actualState, err := game.Reveal(tc.x, tc.y)
			if err != nil && tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected: %v, got: %v", tc.expectedErr, err)
				}
			}
			if actualState != tc.expectedState {
				t.Fatalf("expected: %v, got: %v", tc.expectedState, actualState)
			}
		})
	}
}

func TestMiner_Start(t *testing.T) {
	tests := map[string]struct {
		size        int
		difficulty  int
		expectedErr error
	}{
		"correct settings":          {size: 10, difficulty: 10, expectedErr: nil},
		"wrong settings difficulty": {size: 10, difficulty: 0, expectedErr: ErrInvalidSettings},
		"wrong settings size":       {size: 0, difficulty: 10, expectedErr: ErrInvalidSettings},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			game := NewGame()
			err := game.Start(tc.size, tc.difficulty)
			if err != nil && tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected: %v, got: %v", tc.expectedErr, err)
				}
			}
		})
	}

}

func TestCellGetters(t *testing.T) {
	x := 10
	y := 1
	count := 1
	c := Cell{
		Position: Position{x, y},
		count:    count,
	}
	if c.X() != x {
		t.Fatalf("expected: %v, got: %v", x, c.X())
	}
	if c.Y() != y {
		t.Fatalf("expected: %v, got: %v", y, c.Y())
	}
	if c.Count() != count {
		t.Fatalf("expected: %v, got: %v", count, c.Count())
	}
}
