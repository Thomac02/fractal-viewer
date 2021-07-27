package main

import "testing"

func BenchmarkDrawMandelbrot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		drawMandelbrot()
	}
}
