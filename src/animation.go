package main

import (
	"time"

	"github.com/faiface/pixel"
)

// Animation on a horizontal spreadsheet
type animation struct {
	partialSprite   *pixel.Sprite
	frameSize       pixel.Vec
	currentFrame    int
	totalFrames     int
	secondsPerFrame float64
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

func (a *animation) run() {
	go (func() {
		duration := time.Duration(float64(time.Second) * a.secondsPerFrame)
		ticker := time.Tick(duration)
		for !a.isFinished {
			select {
			case <-ticker:
				if a.isFinished {
					break
				}
				a.advance()
			default:
			}
		}
	})()
}

func (a *animation) advance() {
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
