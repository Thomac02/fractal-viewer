package main

import (
	"flag"
	"fmt"
	"image"

	"github.com/Thomac02/fractal-viewer/fractal"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var (
	fract        fractal.Fractal
	offsetX      float64
	offsetY      float64
	scale        = 1.0
	scaleInverse float64
	startPanX    uint32
	startPanY    uint32
)

func main() {

	workers := flag.Int("workers", 1, "number of workers to use")
	fractalChoice := flag.String("fractal", "", "fractal to display")
	maxIteration := flag.Int("max-iteration", 1000, "maximum iteration for escape time")
	flag.Parse()

	config := fractal.FractalConfig{
		Choice:       *fractalChoice,
		Workers:      uint32(*workers),
		MaxIteration: uint32(*maxIteration),
		ScreenRes:    image.Point{screenWidth, screenHeight},
	}

	fract = fractal.NewFractal(config)

	fmt.Printf("Using %v workers...\n", *workers)
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "fract-viewer",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	picture := pixel.PictureDataFromImage(fract.Draw(scale, offsetX, offsetY))
	sprite := pixel.NewSprite(picture, picture.Bounds())

	win.Clear(colornames.Greenyellow)

	startMatrix := pixel.IM.Moved(win.Bounds().Center())
	sprite.Draw(win, startMatrix)
	currentMatrix := startMatrix

	for !win.Closed() {
		mousePos := win.MousePosition()
		mouseXPos := uint32(mousePos.X)
		mouseYPos := uint32(-mousePos.Y + screenHeight)
		win.Clear(colornames.White)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			startPanX = mouseXPos
			startPanY = mouseYPos
		}

		if win.Pressed(pixelgl.MouseButtonLeft) {
			offsetX -= (float64(mouseXPos) - float64(startPanX)) / scale
			offsetY -= (float64(mouseYPos) - float64(startPanY)) / scale
			startPanX = mouseXPos
			startPanY = mouseYPos
		}

		if win.Pressed(pixelgl.KeySpace) {
			offsetX, offsetY = 0, 0
			scale = 1.0
		}

		mouseXBeforeZoom, mouseYBeforeZoom := screenToWorld(mouseXPos, mouseYPos)
		if win.MouseScroll().Y < 0 {
			scale *= 1.1
			scaleInverse = 1 / scale
		}

		if win.MouseScroll().Y > 0 {
			scale *= 0.9
			scaleInverse = 1 / scale
		}
		mouseXAfterZoom, mouseYAfterZoom := screenToWorld(mouseXPos, mouseYPos)
		offsetX += mouseXBeforeZoom - mouseXAfterZoom
		offsetY += mouseYBeforeZoom - mouseYAfterZoom

		picture = pixel.PictureDataFromImage(fract.Draw(scale, offsetX, offsetY))
		sprite.Set(picture, picture.Bounds())
		sprite.Draw(win, currentMatrix)

		win.Update()
	}
}

func screenToWorld(screenX, screenY uint32) (worldX, worldY float64) {
	return float64(screenX)*scaleInverse + offsetX, float64(screenY)*scaleInverse + offsetY
}

func worldToScreen(worldX, worldY float64) (screenX, screenY uint32) {
	return uint32((worldX - offsetX) * scale), uint32((worldY - offsetY) * scale)
}
