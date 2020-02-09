package main

import (
	"image"
	"image/color"
)

// luminancePercent calculates the relative luminance of a pixel,
func luminancePercent(c color.Color) float64 {
	r, g, b, _ := c.RGBA()

	// We're dividing our pixel's red, green, and blue values by 2^16 because
	// in colors returned from Color.RGBA(), the maximum value for a color
	// is 2^16-1, or 65,535.
	redPercent := float64(r) / 65535 * 100
	greenPercent := float64(g) / 65535 * 100
	bluePercent := float64(b) / 65535 * 100

	return redPercent*0.2126 + greenPercent*0.7152 + bluePercent*0.0722
}

var (
	charcoal  = color.RGBA{R: 34, G: 31, B: 32, A: 255}
	blueGreen = color.RGBA{R: 128, G: 201, B: 172, A: 255}
)

func recolor(img image.Image) image.Image {
	b := img.Bounds()
	recolored := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := img.At(x, y)

			if luminancePercent(color) > 50 {
				recolored.Set(x, y, charcoal)
			} else {
				recolored.Set(x, y, blueGreen)
			}
		}
	}
	return recolored
}
