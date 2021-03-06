package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type player struct {
	idleSprite    *pixel.Sprite
	runningSprite *pixel.Sprite
	explosion     *animation
	isRunning     bool
	isDead        bool
	rect          pixel.Rect
	velocity      pixel.Vec
	direction     pixel.Vec
}

func newPlayer(
	idleSprite *pixel.Sprite,
	runningSprite *pixel.Sprite,
	explosion *animation,
	startPos pixel.Vec,
) *player {
	hitboxSize := idleSprite.Frame().Size().Scaled(.5)

	rect := pixel.R(-hitboxSize.X/2, -hitboxSize.Y/2, hitboxSize.X/2, hitboxSize.Y/2)
	rect = rect.Moved(startPos)
	velocity := pixel.V(0, 0)
	direction := pixel.V(0, 1)

	return &player{
		idleSprite:    idleSprite,
		runningSprite: runningSprite,
		explosion:     explosion,
		rect:          rect,
		velocity:      velocity,
		direction:     direction,
	}
}

func (p *player) notify(ev *event) {
	switch ev.name {
	case "Space Pressed":
		p.destroy()
	case "Left Pressed":
		p.rotate(1, ev.data.(float64)) // dt
	case "Right Pressed":
		p.rotate(-1, ev.data.(float64)) // dt
	case "Forward Pressed":
		p.accelerate(1, ev.data.(float64)) // dt
	case "Forward Not Pressed":
		p.isRunning = false
	}
}

func (p *player) destroy() {
	if !p.isDead {
		p.isRunning = false
		p.isDead = true
		p.velocity = pixel.V(0, 0)
		p.direction = pixel.V(0, 0)
	}
}

func (p *player) rotate(directionMultiplier, dt float64) {
	const rotAngle float64 = math.Pi / 3
	p.direction = p.direction.Rotated(directionMultiplier * rotAngle * dt)
}

func (p *player) accelerate(accelerationMultiplier, dt float64) {
	const acceleration float64 = 600
	p.velocity = p.velocity.Add(p.direction.Scaled(accelerationMultiplier * acceleration * dt))
	p.isRunning = true
}

func (p *player) update(dt float64) {
	g := pixel.V(0, 98)

	if !p.isDead {
		p.velocity = p.velocity.Sub(g.Scaled(dt))
		p.rect = p.rect.Moved(p.velocity.Scaled(dt))
	} else if !p.explosion.isFinished {
		p.explosion.advance(dt)
	}
}

func (p *player) draw(drawTarget *pixel.Target) {
	mat := pixel.IM
	mat = mat.Rotated(pixel.ZV, p.direction.Angle())
	mat = mat.Moved(p.rect.Center())

	if p.isDead {
		if !p.explosion.isFinished {
			p.explosion.partialSprite.Draw(*drawTarget, mat)
		}
	} else if p.isRunning {
		p.runningSprite.Draw(*drawTarget, mat)
	} else {
		p.idleSprite.Draw(*drawTarget, mat)
	}

	if debug {
		imd := imdraw.New(nil)

		imd.Push(p.rect.Center())
		imd.Push(p.rect.Center().Add(p.direction.Scaled(100)))
		imd.Line(3)

		for _, vert := range p.rect.Vertices() {
			imd.Push(vert)
		}
		imd.Polygon(3)

		imd.Draw(*drawTarget)
	}
}
