package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

func loadPicture(path string) *pixel.PictureData {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return pixel.PictureDataFromImage(img)
}

func createPartialSprite(picData *pixel.PictureData, frame pixel.Rect) *pixel.Sprite {
	return pixel.NewSprite(picData, frame)
}

func createFullSprite(picData *pixel.PictureData) *pixel.Sprite {
	return createPartialSprite(picData, picData.Bounds())
}
