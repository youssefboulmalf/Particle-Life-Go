package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Particle struct {
	x  float64
	y  float64
	vx float64
	vy float64
}

type Group struct {
	color    color.RGBA
	group    []Particle
	fuzetemp int
}

type Response struct {
	state bool
	index int
}

var screenWidth float64 = 1024
var screenHeight float64 = 768
var tempature int = 40
var speedIndex float64 = 0.5

var particles = []Group{}
var fuzionGroupIndexes = []int{}

func createRandom(number int, color color.RGBA, fuzetemp int) int {
	var group = []Particle{}

	for i := 0; i < number; i++ {
		var particle = Particle{((rand.Float64() * 924) + 50), ((rand.Float64() * 628) + 50), 0, 0}
		group = append(group, particle)
	}
	particles = append(particles, Group{color, group, fuzetemp})
	return len(particles) - 1
}

func createFuzionGroup(number int, color color.RGBA, fuzetemp int, x float64, y float64) int {
	var group = []Particle{}

	for i := 0; i < number; i++ {
		var particle = Particle{x, y, 0, 0}
		group = append(group, particle)
	}
	particles = append(particles, Group{color, group, fuzetemp})
	fuzionGroupIndexes = append(fuzionGroupIndexes, len(particles)-1)
	return len(particles) - 1
}

func RemoveIndex(s []Particle, index int) []Particle {
	return append(s[:index], s[index+1:]...)
}

func RandInt(lower, upper int) int {
	rand.Seed(time.Now().UnixNano())
	rng := upper - lower
	return rand.Intn(rng) + lower
}

func groupcolorInParticles(a color.RGBA, list []Group) Response {
	for i := 0; i < len(list); i++ {
		if list[i].color == a {
			return Response{state: true, index: i}
		}
	}
	return Response{state: false, index: 0}
}

func rule(groupIndex1 int, groupIndex2 int, g float64) {

	for i := 0; i < len(particles[groupIndex1].group); i++ {

		var fx float64 = 0
		var fy float64 = 0
		var fuze bool = false
		var fuzeColor color.RGBA = colornames.Blue
		for j := 0; j < len(particles[groupIndex2].group); j++ {

			a := particles[groupIndex1].group[i]
			b := particles[groupIndex2].group[j]
			Ta := particles[groupIndex1].fuzetemp
			Tb := particles[groupIndex2].fuzetemp

			dx := a.x - b.x
			dy := a.y - b.y
			d := math.Sqrt(dx*dx + dy*dy)
			if d > 0 && d < 3 && tempature >= Ta && tempature >= Tb && particles[groupIndex1].color != particles[groupIndex2].color {
				particles[groupIndex2].group = RemoveIndex(particles[groupIndex2].group, j)
				fuzeColor.R = particles[groupIndex1].color.R + particles[groupIndex2].color.R
				fuzeColor.G = particles[groupIndex1].color.G + particles[groupIndex2].color.G
				fuzeColor.B = particles[groupIndex1].color.B + particles[groupIndex2].color.B
				fuzeColor.A = particles[groupIndex1].color.A + particles[groupIndex2].color.A
				fuze = true
			}
			if d > 0 && d < 80 {
				F := g * 1 / d
				fx += (F * dx)
				fy += (F * dy)
			}
		}
		if !fuze {
			particles[groupIndex1].group[i].vx = (particles[groupIndex1].group[i].vx + fx) * speedIndex
			particles[groupIndex1].group[i].vy = (particles[groupIndex1].group[i].vy + fy) * speedIndex
			particles[groupIndex1].group[i].x += particles[groupIndex1].group[i].vx
			particles[groupIndex1].group[i].y += particles[groupIndex1].group[i].vy
			if particles[groupIndex1].group[i].x <= 0 || particles[groupIndex1].group[i].x >= screenWidth {
				particles[groupIndex1].group[i].vx *= -1
			}
			if particles[groupIndex1].group[i].y <= 0 || particles[groupIndex1].group[i].y >= screenHeight {
				particles[groupIndex1].group[i].vy *= -1
			}
		}
		if fuze {
			var knownParticle Response = groupcolorInParticles(fuzeColor, particles)
			if knownParticle.state {
				particles[knownParticle.index].group = append(particles[knownParticle.index].group, Particle{particles[groupIndex1].group[i].x, particles[groupIndex1].group[i].y, 0, 0})
				particles[groupIndex1].group = RemoveIndex(particles[groupIndex1].group, i)
			}
			if !knownParticle.state {
				createFuzionGroup(1, fuzeColor, (particles[groupIndex1].fuzetemp + particles[groupIndex2].fuzetemp), particles[groupIndex1].group[i].x, particles[groupIndex1].group[i].y)
				particles[groupIndex1].group = RemoveIndex(particles[groupIndex1].group, i)
			}
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

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	tempatureText := text.New(pixel.V(50, screenHeight-50), basicAtlas)
	FuzionCountText := text.New(pixel.V(50, screenHeight-70), basicAtlas)

	red := createRandom(200, colornames.Red, 50)
	yellow := createRandom(200, colornames.Yellow, 30)
	white := createRandom(200, colornames.Forestgreen, 100)
	blue := createRandom(200, colornames.Blue, 20)

	imd := imdraw.New(nil)
	for !win.Closed() {

		win.Clear(colornames.Black)
		FuzionCountText.Draw(win, pixel.IM)
		tempatureText.Draw(win, pixel.IM)

		imd.Draw(win)

		win.Update()
		tempatureText.Clear()
		FuzionCountText.Clear()
		fmt.Fprintf(tempatureText, "Tempature: %s 'C", strconv.Itoa(tempature))
		fmt.Fprintf(FuzionCountText, "New fuzions: %s", strconv.Itoa(len(fuzionGroupIndexes)))

		rule(white, white, -0.14)
		rule(white, red, -0.17)
		rule(white, yellow, 0.1)
		rule(red, red, 0.1)
		rule(yellow, red, -0.05)
		rule(red, white, 0.1)
		rule(yellow, yellow, 0.15)
		rule(yellow, white, -0.16)
		rule(blue, white, -0.16)
		rule(yellow, blue, -0.17)
		rule(blue, red, 0.11)

		for i := 0; i < len(fuzionGroupIndexes); i++ {
			for j := 0; j < len(particles); j++ {
				rule(fuzionGroupIndexes[i], j, (float64(RandInt(-10, 10)) * 0.03))
			}
		}

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
