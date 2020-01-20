package main

import (
	"github.com/faiface/pixel"
)

// Animation on a horizontal spreadsheet
type animation struct {
	partialSprite   *pixel.Sprite
	frameSize       pixel.Vec
	currentFrame    int
	totalFrames     int
	secondsPerFrame float64
	currentSeconds  float64
	looping         bool
	isFinished      bool
}

func newAnimation(partialSprite *pixel.Sprite, frameSize pixel.Vec, secondsPerFrame float64, looping bool) *animation {
	totalFrames := int(partialSprite.Picture().Bounds().W()) / int(frameSize.Y)
	return &animation{
		partialSprite:   partialSprite,
		frameSize:       frameSize,
		currentFrame:    0,
		totalFrames:     totalFrames,
		secondsPerFrame: secondsPerFrame,
		looping:         looping,
	}
}

func (a *animation) advance(dt float64) {
	if a.isFinished {
		return
	}

	a.currentSeconds += dt
	if a.currentSeconds < a.secondsPerFrame {
		return
	}

	framesToAdvance := int(a.currentSeconds / a.secondsPerFrame)
	a.currentSeconds -= a.secondsPerFrame * float64(framesToAdvance)

	for i := 0; i < framesToAdvance && !a.isFinished; i++ {
		movedRight := a.partialSprite.Frame().Moved(pixel.V(a.frameSize.X, 0))

		// if movedRight is completely inside the spritesheet
		if a.partialSprite.Picture().Bounds().Intersect(movedRight) == movedRight {
			a.partialSprite.Set(a.partialSprite.Picture(), movedRight)
			a.currentFrame++
		} else if a.looping {
			firstFrame := pixel.R(0, 0, a.frameSize.X, a.frameSize.Y)
			a.partialSprite.Set(a.partialSprite.Picture(), firstFrame)
		} else {
			a.isFinished = true
		}
	}
}
