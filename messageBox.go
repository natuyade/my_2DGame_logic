package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawMessage(screen *ebiten.Image, msg string) {
	f := &text.GoTextFace{
		Source:    japaneseFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      16,
	}

	// テキスト原点からの次行の原点までのスペース
	lineSpacing := 16.
	w, h := text.Measure(msg, f, lineSpacing)
	x, y := 0, screenSizeHeight-int(h)

	// msgボックスの背景
	vector.FillRect(screen, float32(x), float32(y), float32(w), float32(h), color.Black, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = lineSpacing
	text.Draw(screen, msg, f, op)
}