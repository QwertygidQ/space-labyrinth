package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type worldMap struct {
	tiles      *[][]bool
	tilesW     int
	tilesH     int
	tileSize   float64
	startTile  pixel.Vec
	startDelta pixel.Vec
	startPos   pixel.Vec
	imd        *imdraw.IMDraw
}

func parseJSONMap(path string) *[][]bool {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var mapData [][]bool
	err = json.Unmarshal(bytes, &mapData)
	if err != nil {
		panic(err)
	}

	return &mapData
}

func checkMapSize(mapData *[][]bool) error {
	if len(*mapData) == 0 {
		return errors.New("Map is empty")
	}

	lineLength := len((*mapData)[0])
	for _, line := range *mapData {
		if len(line) != lineLength {
			return errors.New("Map line length is inconsistent")
		}
	}

	if lineLength == 0 {
		return errors.New("Map is empty")
	}

	return nil
}

func newWorldMap(tiles *[][]bool, tileSize float64, startTile pixel.Vec, startDelta pixel.Vec) *worldMap {
	err := checkMapSize(tiles)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	for y, tileLine := range *tiles {
		for x, tilePresent := range tileLine {
			if !tilePresent {
				continue
			}

			var (
				worldX = float64(x) * tileSize
				worldY = float64(len(*tiles)-1-y) * tileSize
			)

			imd.Color = pixel.RGB(1, 1, 1)
			imd.Push(pixel.V(worldX, worldY))
			imd.Push(pixel.V(worldX+tileSize, worldY))
			imd.Push(pixel.V(worldX+tileSize, worldY+tileSize))
			imd.Push(pixel.V(worldX, worldY+tileSize))
			imd.Polygon(0)
		}
	}

	startTile.Y = float64(len(*tiles) - 1 - int(startTile.Y))
	startPos := startTile.Scaled(tileSize).Add(startDelta)

	return &worldMap{
		tiles:      tiles,
		tileSize:   tileSize,
		tilesW:     len((*tiles)[0]),
		tilesH:     len(*tiles),
		startTile:  startTile,
		startDelta: startDelta,
		startPos:   startPos,
		imd:        imd,
	}
}

func (wm *worldMap) isIntersecting(rect pixel.Rect) bool {
	// If rect is fully outside of the map
	if rect.Max.X < 0 ||
		rect.Min.X > float64(wm.tilesW)*wm.tileSize ||
		rect.Max.Y < 0 ||
		rect.Min.Y > float64(wm.tilesH)*wm.tileSize {
		return false
	}

	// Naive approach: iterate over every wall and check for intersection
	for y, tileLine := range *wm.tiles {
		for x, tilePresent := range tileLine {
			if !tilePresent {
				continue
			}

			var (
				worldX = float64(x) * wm.tileSize
				worldY = float64(len(*wm.tiles)-1-y) * wm.tileSize
			)

			tileRect := pixel.R(worldX, worldY, worldX+wm.tileSize, worldY+wm.tileSize)
			if tileRect.Intersects(rect) {
				return true
			}
		}
	}

	return false
}

func (wm *worldMap) draw(drawTarget *pixel.Target) {
	wm.imd.Draw(*drawTarget)
}
