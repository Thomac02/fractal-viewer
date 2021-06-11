package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	maxIteration   = 1000
	screenWidth    = 1280
	screenHeight   = 720
	mandelbrotXMin = -2.5
	mandelbrotYMin = -1.0
	mandelbrotXMax = 1.0
	mandelbrotYMax = 1.0
)

var (
	mandelbrot    pixel.Picture
	palette       []color.RGBA
	workers       *int
	offsetX       float64
	offsetY       float64
	scaleX        = 1.0
	scaleY        = 1.0
	scaleXInverse float64
	scaleYInverse float64
	startPanX     uint32
	startPanY     uint32
	xMapScale     float64
	yMapScale     float64
)

func main() {

	palette = calculatePalette()
	workers = flag.Int("workers", 1, "number of workers to use")
	flag.Parse()
	xMapScale = (mandelbrotXMax - mandelbrotXMin) / float64(screenWidth)
	yMapScale = (mandelbrotYMin - mandelbrotYMax) / float64(screenHeight)

	fmt.Printf("Using %v workers...\n", *workers)
	img := drawMandelbrot()
	mandelbrot = pixel.PictureDataFromImage(img)
	pixelgl.Run(run)
}

func mapVal(x, imin, imax, omin, omax float64) float64 {
	return (x-imin)*(omax-omin)/(imax-imin) + omin
}

func drawMandelbrot() *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{screenWidth, screenHeight}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	heightChunk := int(screenHeight / *workers)
	wg := &sync.WaitGroup{}
	wg.Add(*workers)
	for i := 0; i < *workers; i++ {
		go func(starty int, wg *sync.WaitGroup) {
			for Px := 0; Px < screenWidth; Px++ {
				for Py := starty; Py < starty+heightChunk; Py++ {
					x0, y0 := mandelbrotScale(uint32(Px), uint32(Py))

					// Cardioid checking
					y02 := y0 * y0
					q := (x0-0.25)*(x0-0.25) + y02
					cardioid := q*(q+(x0-0.25)) <= 0.25*y02
					bulb := (x0+1)*(x0+1)+y02 <= 0.0625
					if cardioid || bulb {
						color := palette[maxIteration-1]
						img.Set(Px, Py, color)
						continue
					}

					// Optimised escape time
					var x, y, x2, y2 float64
					iteration := uint32(0)
					for x2+y2 <= 4 && iteration < maxIteration {
						y = (x+x)*y + y0
						x = x2 - y2 + x0
						x2 = x * x
						y2 = y * y
						iteration++
					}
					color := palette[iteration-1]
					img.Set(Px, Py, color)
				}
			}
			wg.Done()
		}(i*heightChunk, wg)
	}
	wg.Wait()
	return img
}

func mandelbrotScale(x, y uint32) (float64, float64) {
	x_f, y_f := screenToWorld(x, y)
	scaledx := x_f*xMapScale + mandelbrotXMin
	scaledy := y_f*yMapScale + mandelbrotYMax
	return scaledx, scaledy
}

func calculatePalette() []color.RGBA {
	pal := make([]color.RGBA, maxIteration, maxIteration)
	for i := 0; i < maxIteration; i++ {
		inew := uint32(mapVal(float64(i), 0.0, float64(maxIteration-1), float64(0xFFFFFF), 0.0))
		pal[i] = color.RGBA{uint8(inew >> 16), uint8(inew >> 8), uint8(inew), 0xFF}
	}
	return pal
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "mandelbrot-viewer",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	mandelbrot = pixel.PictureDataFromImage(drawMandelbrot())
	sprite := pixel.NewSprite(mandelbrot, mandelbrot.Bounds())

	win.Clear(colornames.Greenyellow)

	startMatrix := pixel.IM.Moved(win.Bounds().Center())
	sprite.Draw(win, startMatrix)
	currentMatrix := startMatrix

	scaleX, scaleY = 1.0, 1.0
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
			offsetX -= (float64(mouseXPos) - float64(startPanX)) / scaleX
			offsetY -= (float64(mouseYPos) - float64(startPanY)) / scaleY
			startPanX = mouseXPos
			startPanY = mouseYPos
		}

		if win.Pressed(pixelgl.KeySpace) {
			offsetX, offsetY = 0, 0
			scaleX, scaleY = 1, 1
		}

		mouseXBeforeZoom, mouseYBeforeZoom := screenToWorld(mouseXPos, mouseYPos)
		if win.MouseScroll().Y < 0 {
			scaleX *= 1.1
			scaleY *= 1.1
			scaleXInverse = 1 / scaleX
			scaleYInverse = 1 / scaleY
		}

		if win.MouseScroll().Y > 0 {
			scaleX *= 0.9
			scaleY *= 0.9
			scaleXInverse = 1 / scaleX
			scaleYInverse = 1 / scaleY
		}
		mouseXAfterZoom, mouseYAfterZoom := screenToWorld(mouseXPos, mouseYPos)
		offsetX += mouseXBeforeZoom - mouseXAfterZoom
		offsetY += mouseYBeforeZoom - mouseYAfterZoom

		//fmt.Printf("startPanX: %v startPanY: % v offsetX: %v offsetY: %v scaleX: %v scaleY: %v\n", startPanX, startPanY, offsetX, offsetY, scaleX, scaleY)
		//fmt.Printf("mouseX: %v, mouseY: %v\n", mouseXPos, mouseYPos)

		renderTimeStart := time.Now()
		mandelbrot = pixel.PictureDataFromImage(drawMandelbrot())
		renderTime := time.Since(renderTimeStart)
		fmt.Printf("render time: %s\n", renderTime)
		sprite.Set(mandelbrot, mandelbrot.Bounds())
		sprite.Draw(win, currentMatrix)

		win.Update()
	}
}

func screenToWorld(screenX, screenY uint32) (worldX, worldY float64) {
	return float64(screenX)*scaleXInverse + offsetX, float64(screenY)*scaleYInverse + offsetY
}

func worldToScreen(worldX, worldY float64) (screenX, screenY uint32) {
	return uint32((worldX - offsetX) * scaleX), uint32((worldY - offsetY) * scaleY)
}
