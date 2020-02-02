package main

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

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

func recolor(img image.Image) {
	b := img.Bounds()

	var lightPixels, darkPixels int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := img.At(x, y)

			// [TODO] Replace this with actually setting the color
			// of the current pixel.
			if luminancePercent(color) > 50 {
				lightPixels++
			} else {
				darkPixels++
			}
		}
	}
	log.Printf("This image has %d pixels above 50%% luminance and %d below 50%%\n",
		lightPixels, darkPixels)
}

func main() {
	f, err := os.Open("raw.jpg")
	if err != nil {
		log.Fatalf("error loading raw.jpg: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("error decoding file to an Image: %v", err)
	}

	recolor(img)
}
