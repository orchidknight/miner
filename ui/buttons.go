package ui

import (
	"image/color"

	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/mouse"
	"github.com/oakmound/oak/v4/render"
)

type button struct {
	id         event.CallerID
	Position   Position
	Shape      Shape
	hoverColor color.RGBA
	color      color.RGBA
	mouse.CollisionPhase
	*render.ColorBoxR
}

type cellButton struct {
	button
	Selected bool
	x, y     int
	revealed bool
	count    int
	Position Position
}

type difficultyButton struct {
	button
	difficultyButtons map[Difficulty]*difficultyButton
	selected          bool
	difficulty        Difficulty
}

type sizeButton struct {
	sizeButtons map[Size]*sizeButton
	selected    bool
	size        Size
	button
}

type startButton struct {
	button
}

type backButton struct {
	button
}
