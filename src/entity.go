package main

import "github.com/faiface/pixel"

type entity interface {
	update(dt float64)
	draw(drawTarget *pixel.Target)
}
