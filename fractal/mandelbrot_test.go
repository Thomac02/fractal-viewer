package fractal

import (
	"fmt"
	"math"
	"testing"
)

func BenchmarkDraw60Workers(b *testing.B) {
	mand := &mandelbrot{
		resX:         1280,
		resY:         720,
		workers:      60,
		maxIteration: 1000,
		xMapScale:    (mandelbrotXMax - mandelbrotXMin) / float64(1280),
		yMapScale:    (mandelbrotYMin - mandelbrotYMax) / float64(720),
	}
	mand.palette = mand.calculatePalette()
	for i := 0; i < b.N; i++ {
		mand.Draw(1.0, 0.0, 0.0)
	}
}

func BenchmarkDraw120Workers(b *testing.B) {
	mand := &mandelbrot{
		resX:         1280,
		resY:         720,
		workers:      120,
		maxIteration: 1000,
		xMapScale:    (mandelbrotXMax - mandelbrotXMin) / float64(1280),
		yMapScale:    (mandelbrotYMin - mandelbrotYMax) / float64(720),
	}
	mand.palette = mand.calculatePalette()
	for i := 0; i < b.N; i++ {
		mand.Draw(1.0, 0.0, 0.0)
	}
}

func BenchmarkDraw240Workers(b *testing.B) {
	mand := &mandelbrot{
		resX:         1280,
		resY:         720,
		workers:      240,
		maxIteration: 1000,
		xMapScale:    (mandelbrotXMax - mandelbrotXMin) / float64(1280),
		yMapScale:    (mandelbrotYMin - mandelbrotYMax) / float64(720),
	}
	mand.palette = mand.calculatePalette()
	for i := 0; i < b.N; i++ {
		mand.Draw(1.0, 0.0, 0.0)
	}
}

func BenchmarkDrawChangingScale(b *testing.B) {
	mand := &mandelbrot{
		resX:         1280,
		resY:         720,
		workers:      60,
		maxIteration: 1000,
		xMapScale:    (mandelbrotXMax - mandelbrotXMin) / float64(1280),
		yMapScale:    (mandelbrotYMin - mandelbrotYMax) / float64(720),
	}
	mand.palette = mand.calculatePalette()
	scale := math.Pow(1.1, 10)
	fmt.Println(scale)
	for i := 0; i < b.N; i++ {
		mand.Draw(scale, 0.0, 0.0)
	}
}
