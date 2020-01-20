package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Lander",
		Bounds: pixel.R(0, 0, 1024, 768),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	drawTarget := pixel.Target(win)

	playerIdleSprite := createFullSprite(loadPicture("img/idleShuttle.png"))
	playerRunningSprite := createFullSprite(loadPicture("img/runningShuttle.png"))
	playerInstance := newPlayer(playerIdleSprite, playerRunningSprite, win.Bounds().Center())

	keyboardEventManager := newEventManager()
	sub := subscriber(playerInstance)
	keyboardEventManager.subscribe("Left Pressed", &sub)
	keyboardEventManager.subscribe("Right Pressed", &sub)
	keyboardEventManager.subscribe("Forward Pressed", &sub)
	keyboardEventManager.subscribe("Forward Not Pressed", &sub)

	lastTime := time.Now()
	for !win.Closed() {
		win.Update()

		dt := time.Since(lastTime).Seconds()
		lastTime = time.Now()

		if win.Pressed(pixelgl.KeyW) || win.Pressed(pixelgl.KeyUp) {
			keyboardEventManager.notifySubscribers(&event{name: "Forward Pressed", data: dt})
		} else {
			keyboardEventManager.notifySubscribers(&event{name: "Forward Not Pressed"})
		}
		if win.Pressed(pixelgl.KeyS) || win.Pressed(pixelgl.KeyDown) {
			keyboardEventManager.notifySubscribers(&event{name: "Backward Pressed", data: dt})
		}
		if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyLeft) {
			keyboardEventManager.notifySubscribers(&event{name: "Left Pressed", data: dt})
		}
		if win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
			keyboardEventManager.notifySubscribers(&event{name: "Right Pressed", data: dt})
		}

		playerInstance.update(dt)

		win.Clear(colornames.Black)
		playerInstance.draw(&drawTarget)
	}
}

func main() {
	pixelgl.Run(run)
}
