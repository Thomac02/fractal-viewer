package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	maxIteration = 10000
	width        = 1920
	height       = 1080
	concurrent   = true
	workers      = 16
)

var mandelbrot pixel.Picture

func main() {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	palette := calculatePalette()
	heightChunk := int(height / workers)

	fmt.Println("Creating image...")
	startTime := time.Now()
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	if concurrent && workers != 1 {
		fmt.Printf("Using %v workers...\n", workers)
		wg := &sync.WaitGroup{}
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func(starty int, wg *sync.WaitGroup) {
				for Px := 0; Px < width; Px++ {
					for Py := starty; Py < starty+heightChunk; Py++ {
						x0, y0 := mandelbrotScale(Px, Py)
						var x, y float64
						iteration := uint32(0)
						for x*x+y*y <= 2*2 && iteration < maxIteration {
							xtemp := x*x - y*y + x0
							y = 2*x*y + y0
							x = xtemp
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
	} else {
		for Px := 0; Px < width; Px++ {
			for Py := 0; Py < height; Py++ {
				x0, y0 := mandelbrotScale(Px, Py)
				var x, y float64
				iteration := uint32(0)
				for x*x+y*y <= 2*2 && iteration < maxIteration {
					xtemp := x*x - y*y + x0
					y = 2*x*y + y0
					x = xtemp
					iteration++
				}
				color := palette[iteration-1]
				img.Set(Px, Py, color)
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("images/image.png")
	png.Encode(f, img)
	fmt.Printf("Done, took %s\n", time.Since(startTime))
	mandelbrot = pixel.PictureDataFromImage(img)

	pixelgl.Run(run)
}

func mapVal(x, imin, imax, omin, omax float64) float64 {
	return (x-imin)*(omax-omin)/(imax-imin) + omin
}

func mandelbrotScale(x, y int) (float64, float64) {
	scaledx := mapVal(float64(x), 0.0, float64(width), -2.5, 1.0)
	scaledy := mapVal(float64(y), 0.0, float64(height), 1.0, -1.0)
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
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(mandelbrot, mandelbrot.Bounds())

	win.Clear(colornames.Greenyellow)

	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	for !win.Closed() {
		win.Update()
	}
}
