package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Particle struct {
	x  float64
	y  float64
	vx float64
	vy float64
}

type Group struct {
	color color.RGBA
	group []Particle
}

var screenWidth float64 = 1024
var screenHeight float64 = 768

var particles = []Group{}

func create(number int, color color.RGBA) int {
	var group = []Particle{}

	for i := 0; i < number; i++ {
		var particle = Particle{((rand.Float64() * 924) + 50), ((rand.Float64() * 628) + 50), 0, 0}
		group = append(group, particle)
	}
	particles = append(particles, Group{color, group})
	return len(particles) - 1
}

func rule(groupIndex1 int, groupIndex2 int, g float64) {

	for i := 0; i < len(particles[groupIndex1].group); i++ {

		var fx float64 = 0
		var fy float64 = 0
		for j := 0; j < len(particles[groupIndex2].group); j++ {

			a := particles[groupIndex1].group[i]
			b := particles[groupIndex2].group[j]

			dx := a.x - b.x
			dy := a.y - b.y
			d := math.Sqrt(dx*dx + dy*dy)
			if d > 0 && d < 100 {
				F := g * 1 / d
				fx += (F * dx)
				fy += (F * dy)

			}
		}
		particles[groupIndex1].group[i].vx = (particles[groupIndex1].group[i].vx + fx) * 0.5
		particles[groupIndex1].group[i].vy = (particles[groupIndex1].group[i].vy + fy) * 0.5
		particles[groupIndex1].group[i].x += particles[groupIndex1].group[i].vx
		particles[groupIndex1].group[i].y += particles[groupIndex1].group[i].vy
		if particles[groupIndex1].group[i].x <= 0 || particles[groupIndex1].group[i].x >= screenWidth {
			particles[groupIndex1].group[i].vx *= -1
		}
		if particles[groupIndex1].group[i].y <= 0 || particles[groupIndex1].group[i].y >= screenHeight {
			particles[groupIndex1].group[i].vy *= -1
		}

	}

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Artificial Life",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	red := create(200, colornames.Red)
	yellow := create(200, colornames.Yellow)
	green := create(200, colornames.Green)

	imd := imdraw.New(nil)
	for !win.Closed() {

		win.Clear(colornames.Black)
		imd.Draw(win)

		win.Update()
		rule(green, green, -0.32)
		rule(green, red, -0.17)
		rule(green, yellow, 0.34)
		rule(red, red, -0.10)
		rule(red, green, -0.34)
		rule(yellow, yellow, 0.15)
		rule(yellow, green, -0.20)

		imd.Reset()
		imd.Clear()
		for i := 0; i < len(particles); i++ {
			for j := 0; j < len(particles[i].group); j++ {
				imd.Color = particles[i].color
				imd.EndShape = imdraw.SharpEndShape
				imd.Push(pixel.V(particles[i].group[j].x-2, particles[i].group[j].y-2), pixel.V(particles[i].group[j].x+2, particles[i].group[j].y+2))
				imd.Rectangle(0)
			}
		}

	}

}

func main() {
	pixelgl.Run(run)
}
