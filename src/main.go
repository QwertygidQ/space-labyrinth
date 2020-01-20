package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Space Labyrinth",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	drawTarget := pixel.Target(win)

	mapData := parseJSONMap("misc/map.json")
	labyrinthMap := newWorldMap(mapData, 128, pixel.V(1, 1), pixel.V(128/2, -128/2))

	playerIdleSprite := createFullSprite(loadPicture("img/idleShuttle.png"))
	playerRunningSprite := createFullSprite(loadPicture("img/runningShuttle.png"))
	explosionPartialSprite := createPartialSprite(loadPicture("img/explosion.png"), pixel.R(0, 0, 32, 32))
	playerExplosion := newAnimation(explosionPartialSprite, pixel.V(32, 32), .1, false)
	playerInstance := newPlayer(playerIdleSprite, playerRunningSprite, playerExplosion, labyrinthMap.startPos)

	keyboardEventManager := newEventManager()
	sub := subscriber(playerInstance)
	keyboardEventManager.subscribe("Left Pressed", &sub)
	keyboardEventManager.subscribe("Right Pressed", &sub)
	keyboardEventManager.subscribe("Forward Pressed", &sub)
	keyboardEventManager.subscribe("Forward Not Pressed", &sub)
	keyboardEventManager.subscribe("Space Pressed", &sub)

	frames := 0
	ticker := time.Tick(time.Second)

	lastTime := time.Now()
	for !win.Closed() {
		win.Update()

		frames++
		select {
		case <-ticker:
			win.SetTitle(fmt.Sprintf("Space Labyrinth | %d FPS", frames))
			frames = 0
		default:
		}

		dt := time.Since(lastTime).Seconds()
		lastTime = time.Now()

		if win.Pressed(pixelgl.KeySpace) {
			keyboardEventManager.notifySubscribers(&event{name: "Space Pressed", data: dt})
		}
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

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(playerInstance.rect.Center()))
		win.SetMatrix(cam)

		win.Clear(colornames.Black)
		playerInstance.draw(&drawTarget)
		labyrinthMap.draw(&drawTarget)
	}
}

func main() {
	pixelgl.Run(run)
}
