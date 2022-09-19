package main

import (
	crrand "crypto/rand"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/lucasb-eyer/go-colorful"
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

type Button struct {
	name   string
	sprite *pixel.Sprite
	x      float64
	y      float64
	scaleX float64
	scaleY float64
	rotate float64
	action func()
}

type ParticleRules struct {
	fuzionParticle1 int
	fuzionParticle2 int
	force           float64
}

type ParticleCreation struct {
	amount     int
	color      color.RGBA
	fuzionTemp int
}

var screenWidth float64 = 1600
var screenHeight float64 = 900
var tempature int = 0
var particleSize = 2
var red ParticleCreation = ParticleCreation{200, colornames.Red, 20}
var yellow ParticleCreation = ParticleCreation{200, colornames.Yellow, 40}
var green ParticleCreation = ParticleCreation{200, colornames.Green, 100}
var startingParticles = []ParticleCreation{red, yellow, green}
var isPauzed = false

var fuzionRules = []ParticleRules{}
var normalRules = []ParticleRules{}

var particles = []Group{}
var fuzionGroupIndexes = []int{}

var backgroundcolor color.RGBA = color.RGBA{0, 0, 0, 0}
var speedIndex float64 = 0.5

func calculateSpeedIndex() {
	if tempature < 0 {
		speedIndex -= 0.011
	} else {
		speedIndex += 0.011
	}
	if tempature <= -273 {
		speedIndex = 0
	}
}

func createRandom(number int, color color.RGBA, fuzetemp int) int {
	var group = []Particle{}

	for i := 0; i < number; i++ {
		var particle = Particle{(rand.Float64() * screenWidth), (rand.Float64() * screenHeight), 0, 0}
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
	createFuzionRules()
	return len(particles) - 1
}

func createFuzionRules() {
	for i := 0; i < len(fuzionGroupIndexes); i++ {
		for j := 0; j < len(particles); j++ {
			if RandInt(1, 3) > 1 {
				fuzionRules = append(fuzionRules, ParticleRules{i, j, (float64(RandInt(-12, 12)) * 0.03)})
			}
		}
	}
}

func RemoveIndex(s []Particle, index int) []Particle {
	return append(s[:index], s[index+1:]...)
}

func RandInt(lower, upper int) int {
	seed, err := crrand.Int(crrand.Reader, big.NewInt(27))
	if err == nil {
		rand.Seed(seed.Int64())
	}
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

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func increaseTemp() {
	tempature += 10
	if tempature > 160 {
		backgroundcolor.R += 5
	}
	// backgroundcolor.B -= 1
	calculateSpeedIndex()
}
func decreaseTemp() {
	if tempature > -273 {
		tempature -= 10
		calculateSpeedIndex()
		if backgroundcolor.R <= 0 {
			backgroundcolor.B += 1
		} else {
			backgroundcolor.R -= 1
		}
	}

}

func reset() {
	fuzionRules = []ParticleRules{}
	normalRules = []ParticleRules{}
	particles = []Group{}
	fuzionGroupIndexes = []int{}
	tempature = 0
	speedIndex = 0.5
	backgroundcolor.R = 0
	backgroundcolor.G = 0
	backgroundcolor.B = 0
	backgroundcolor.A = 0
	for i := 0; i < len(startingParticles); i++ {
		createRandom(startingParticles[i].amount, startingParticles[i].color, startingParticles[i].fuzionTemp)
	}
	for i := 0; i < len(particles); i++ {
		for j := 0; j < len(particles); j++ {
			if RandInt(1, 3) > 1 {
				normalRules = append(normalRules, ParticleRules{i, j, (float64(RandInt(-12, 12)) * 0.03)})
			}
		}
	}

}

func pauze() {
	isPauzed = true
}

func start() {
	isPauzed = false
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
				r, g, b, a := colorful.FastHappyColor().Clamped().RGBA()
				fuzeColor.R = uint8(r)
				fuzeColor.G = uint8(g)
				fuzeColor.B = uint8(b)
				fuzeColor.A = uint8(a)
				fuze = true
			}
			if d > 0 && d < 60 {
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
			if particles[groupIndex1].group[i].x < 2 || particles[groupIndex1].group[i].x > screenWidth-2*float64(particleSize) {
				particles[groupIndex1].group[i].vx *= -1
			}
			if particles[groupIndex1].group[i].y < 2 || particles[groupIndex1].group[i].y > screenHeight-2*float64(particleSize) {
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
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//text creation

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	tempatureText := text.New(pixel.V(50, screenHeight-90), basicAtlas)
	FuzionCountText := text.New(pixel.V(50, screenHeight-120), basicAtlas)

	spritesheet, err := loadPicture("./sprites/menuButtons.png")

	//button creation
	buttons := []Button{}
	var buttonFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 18 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 16 {
			buttonFrames = append(buttonFrames, pixel.R(x, y, x+18, y+16))
		}
	}

	tempForwardButton := Button{"tempForwardButton", pixel.NewSprite(spritesheet, buttonFrames[0]), 220, (screenHeight - 87), 1, 1, 0, increaseTemp}
	tempBackButton := Button{"tempBackButton", pixel.NewSprite(spritesheet, buttonFrames[4]), 140, (screenHeight - 87), 1, 1, 0, decreaseTemp}
	resetButton := Button{"tempBackButton", pixel.NewSprite(spritesheet, buttonFrames[90]), 135, (screenHeight - 47), 1, 1, 0, reset}
	pauzeButton := Button{"tempBackButton", pixel.NewSprite(spritesheet, buttonFrames[34]), 95, (screenHeight - 47), 1, 1, 0, pauze}
	startButton := Button{"tempBackButton", pixel.NewSprite(spritesheet, buttonFrames[12]), 55, (screenHeight - 47), 1, 1, 0, start}
	buttons = append(buttons, tempBackButton)
	buttons = append(buttons, tempForwardButton)
	buttons = append(buttons, resetButton)
	buttons = append(buttons, pauzeButton)
	buttons = append(buttons, startButton)

	//generating starting particles
	for i := 0; i < len(startingParticles); i++ {
		createRandom(startingParticles[i].amount, startingParticles[i].color, startingParticles[i].fuzionTemp)
	}

	//aplying rules
	for i := 0; i < len(particles); i++ {
		for j := 0; j < len(particles); j++ {
			normalRules = append(normalRules, ParticleRules{i, j, (float64(RandInt(-20, 20)) * 0.03)})
		}
	}

	imd := imdraw.New(nil)
	for !win.Closed() {

		win.Clear(backgroundcolor)

		//draw text
		FuzionCountText.Draw(win, pixel.IM)
		tempatureText.Draw(win, pixel.IM)

		//draw particles
		imd.Draw(win)

		//draw buttons
		for i := 0; i < len(buttons); i++ {
			button := buttons[i]
			sprite := button.sprite
			sprite.Draw(win, pixel.IM.Moved(pixel.V(float64(button.x), float64(button.y))).Rotated(pixel.ZV, float64(button.rotate)).ScaledXY(pixel.ZV, pixel.V(button.scaleX, button.scaleY)))
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			for i := 0; i < len(buttons); i++ {
				if win.MousePosition().X < buttons[i].x+12 && win.MousePosition().X > buttons[i].x-12 && win.MousePosition().Y < buttons[i].y+12 && win.MousePosition().Y > buttons[i].y-12 {
					buttons[i].action()
				}
			}
		}

		win.Update()
		tempatureText.Clear()
		FuzionCountText.Clear()
		fmt.Fprintf(tempatureText, "Tempature       %s 'C", strconv.Itoa(tempature))
		fmt.Fprintf(FuzionCountText, "New fuzions: %s", strconv.Itoa(len(fuzionGroupIndexes)))

		if !isPauzed {
			for i := 0; i < len(fuzionRules); i++ {

				rule(fuzionRules[i].fuzionParticle1, fuzionRules[i].fuzionParticle2, fuzionRules[i].force)
			}
			for i := 0; i < len(normalRules); i++ {
				rule(normalRules[i].fuzionParticle1, normalRules[i].fuzionParticle2, normalRules[i].force)
			}

			//Cleaning screen between frames
			imd.Reset()
			imd.Clear()

			//Change positions of particles for drawing
			for i := 0; i < len(particles); i++ {
				for j := 0; j < len(particles[i].group); j++ {
					imd.Color = particles[i].color
					imd.EndShape = imdraw.SharpEndShape
					imd.Push(pixel.V(particles[i].group[j].x-float64(particleSize), particles[i].group[j].y-float64(particleSize)), pixel.V(particles[i].group[j].x+float64(particleSize), particles[i].group[j].y+float64(particleSize)))
					imd.Rectangle(0)
				}
			}

		}
	}

}

func main() {
	pixelgl.Run(run)
}
