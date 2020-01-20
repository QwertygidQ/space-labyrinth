package main

import "github.com/faiface/pixel"

// Animation on a horizontal spreadsheet
type animation struct {
	partialSprite *pixel.Sprite
	frameSize     pixel.Vec
	currentFrame  int
	totalFrames   int
}

func (a *animation) advance() {
	movedRight := a.partialSprite.Frame().Moved(pixel.V(a.frameSize.X, 0))

	// if movedRight is completely inside the spritesheet
	if a.partialSprite.Picture().Bounds().Intersect(movedRight) == movedRight {
		a.partialSprite.Set(a.partialSprite.Picture(), movedRight)
		a.currentFrame++
	} else {
		firstFrame := pixel.R(0, 0, a.frameSize.X, a.frameSize.Y)
		a.partialSprite.Set(a.partialSprite.Picture(), firstFrame)
	}
}
