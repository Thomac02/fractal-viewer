package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

const (
	maxIteration = 10000
	width        = 3840
	height       = 2160
)

func main() {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	palette := calculatePalette()

	fmt.Println("Creating image...")
	startTime := time.Now()
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
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

	// Encode as PNG.
	f, _ := os.Create("images/image.png")
	png.Encode(f, img)
	fmt.Printf("Done, took %s\n", time.Since(startTime))
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
