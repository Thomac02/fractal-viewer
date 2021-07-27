package fractal

import (
	"fmt"
	"image"
)

var fractals = map[string]string{"mandelbrot": ""}

type Fractal interface {
	// Draw returns an image of a fractal at a given scale and offset
	Draw(scale, offsetX, offsetY float64) *image.RGBA
}

type FractalConfig struct {
	Choice       string
	Workers      uint32
	MaxIteration uint32
	ScreenRes    image.Point
}

func mapVal(x, imin, imax, omin, omax float64) float64 {
	return (x-imin)*(omax-omin)/(imax-imin) + omin
}

func NewFractal(config FractalConfig) Fractal {
	if _, ok := fractals[config.Choice]; !ok {
		return nil
	}

	switch config.Choice {
	case "mandelbrot":
		return NewMandelbrot(uint32(config.ScreenRes.X), uint32(config.ScreenRes.Y), config.Workers, config.MaxIteration)
	default:
		fmt.Printf("Unknown fractal choice: %s", config.Choice)
	}
	return nil
}
