package fractal

import (
	"image"
	"image/color"
	"sync"
)

const (
	mandelbrotXMin = -2.5
	mandelbrotYMin = -1.0
	mandelbrotXMax = 1.0
	mandelbrotYMax = 1.0
)

type mandelbrot struct {
	resX         uint32
	resY         uint32
	workers      uint32
	maxIteration uint32
	xMapScale    float64
	yMapScale    float64
	palette      []color.RGBA
}

func NewMandelbrot(resX, resY, workers, maxIteration uint32) (m *mandelbrot) {
	m = &mandelbrot{
		resX:         resX,
		resY:         resY,
		workers:      workers,
		maxIteration: maxIteration,
		xMapScale:    (mandelbrotXMax - mandelbrotXMin) / float64(resX),
		yMapScale:    (mandelbrotYMin - mandelbrotYMax) / float64(resY),
	}
	m.palette = m.calculatePalette()
	return
}

func (m *mandelbrot) Draw(scale, offsetX, offsetY float64) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(m.resX), int(m.resY)}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	heightChunk := int(m.resY / m.workers)
	wg := &sync.WaitGroup{}
	wg.Add(int(m.workers))

	// Calculate scale inverse to avoid unnecessary divisions in the loop.
	scaleInverse := 1 / scale
	for i := 0; i < int(m.workers); i++ {
		startY := i * heightChunk
		endY := startY + heightChunk
		if i == int(m.workers)-1 && endY < int(m.resY) {
			endY = int(m.resY)
		}
		go func(starty, endy int, wg *sync.WaitGroup) {
			for Px := 0; Px < int(m.resX); Px++ {
				for Py := starty; Py < endy; Py++ {
					x0, y0 := m.scale(uint32(Px), uint32(Py), scaleInverse, offsetX, offsetY)

					// Cardioid checking
					y02 := y0 * y0
					q := (x0-0.25)*(x0-0.25) + y02
					cardioid := q*(q+(x0-0.25)) <= 0.25*y02
					bulb := (x0+1)*(x0+1)+y02 <= 0.0625
					if cardioid || bulb {
						color := m.palette[m.maxIteration-1]
						img.Set(Px, Py, color)
						continue
					}

					// Optimised escape time
					var x, y, x2, y2 float64
					iteration := uint32(0)
					for x2+y2 <= 4 && iteration < m.maxIteration {
						y = (x+x)*y + y0
						x = x2 - y2 + x0
						x2 = x * x
						y2 = y * y
						iteration++
					}
					color := m.palette[iteration-1]
					img.Set(Px, Py, color)
				}
			}
			wg.Done()
		}(startY, endY, wg)
	}
	wg.Wait()
	return img
}

func (m *mandelbrot) scale(x, y uint32, scale, offsetX, offsetY float64) (float64, float64) {
	x_f, y_f := float64(x)*scale+offsetX, float64(y)*scale+offsetY
	scaledx := x_f*m.xMapScale + mandelbrotXMin
	scaledy := y_f*m.yMapScale + mandelbrotYMax
	return scaledx, scaledy
}

func (m *mandelbrot) calculatePalette() []color.RGBA {
	pal := make([]color.RGBA, m.maxIteration, m.maxIteration)
	for i := 0; i < int(m.maxIteration); i++ {
		inew := uint32(mapVal(float64(i), 0.0, float64(m.maxIteration-1), float64(0xFFFFFF), 0.0))
		pal[i] = color.RGBA{uint8(inew >> 16), uint8(inew >> 8), uint8(inew), 0xFF}
	}
	return pal
}
