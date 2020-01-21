package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type camera struct {
	canvas  *pixelgl.Canvas
	camPos  pixel.Vec
	camRot  float64
	camZoom pixel.Vec // more like stretch
}

func newCamera(canvas *pixelgl.Canvas, camPos pixel.Vec, camRot float64, camZoom pixel.Vec) *camera {
	return &camera{canvas: canvas, camPos: camPos, camRot: camRot, camZoom: camZoom}
}

func (c *camera) moveTo(newPos pixel.Vec) {
	c.camPos = newPos
}

func (c *camera) moveBy(deltaPos pixel.Vec) {
	c.camPos = c.camPos.Add(deltaPos)
}

func (c *camera) setRotation(angle float64) {
	c.camRot = angle
}

func (c *camera) rotateBy(angle float64) {
	c.camRot += angle
}

func (c *camera) setZoom(zoom pixel.Vec) {
	c.camZoom = zoom
}

func (c *camera) zoomBy(zoom pixel.Vec) {
	c.camZoom = c.camZoom.Add(zoom)
}

func (c *camera) getMatrix() pixel.Matrix {
	scaled := pixel.IM.ScaledXY(c.camPos, c.camZoom)
	rotated := scaled.Rotated(c.camPos, c.camRot)
	moved := rotated.Moved(c.canvas.Bounds().Center().Sub(c.camPos))
	return moved
}

func (c *camera) applyChanges() {
	c.canvas.SetMatrix(c.getMatrix())
}
