package ui

import (
	"image"

	"github.com/oakmound/oak/v4/render"
)

func newFont() (*render.Font, error) {
	fg := render.DefaultFontGenerator
	fg.Color = image.NewUniform(fontColor)
	fg.FontOptions.Size = 20
	return fg.Generate()
}
