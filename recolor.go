package main

import (
	"fmt"
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

func recolor(img image.Image, thresholds []luminanceThreshold) image.Image {
	b := img.Bounds()
	recolored := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := img.At(x, y)

			var colorPicked bool
			for i := len(thresholds) - 1; i >= 0; i-- {
				thresholdLuminance := float64(thresholds[i].luminancePercent)
				if luminancePercent(color) >= thresholdLuminance {
					recolored.Set(x, y, thresholds[i].color)
					colorPicked = true
					break
				}
			}
			// The parseThresholds function should always error if
			// there is no zero percent luminance threshold, so
			// this panic is here since this should only ever
			// happen from a programmer error
			if !colorPicked {
				r, g, b, _ := color.RGBA()
				panic(fmt.Sprintf(
					"color %d,%d,%d luminance %f "+
						"below all luminance thresholds %v",
					r, g, b, luminancePercent(color), thresholds,
				))
			}
		}
	}
	return recolored
}
