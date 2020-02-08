package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
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

func recolor(img image.Image) image.Image {
	b := img.Bounds()
	recolored := image.NewRGBA(b)

	// Add color.Color variables for the colors we're recoloring to
	charcoal := color.RGBA{R: 34, G: 31, B: 32, A: 255}
	blueGreen := color.RGBA{R: 128, G: 201, B: 172, A: 255}

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

	if _, err := os.Stat("recolored.jpg"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("unexpected error checking if recolored.jpg already exists")
	} else if err == nil {
		fmt.Println("recolored.jpg already exists. Replace it? (y/n)")

		replace, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalf(
				"unexpected error reading whether to replace the file: %v",
				err,
			)
		} else if strings.TrimRight(strings.ToLower(replace), "\n") != "y" {
			return
		}
	}

	out, err := os.Create("recolored.jpg")
	if err != nil {
		log.Fatalf("error creating output file recolored.jpg: %v", err)
	}
	defer out.Close()

	recolored := recolor(img)
	if err := jpeg.Encode(out, recolored, nil); err != nil {
		log.Fatalf("error outputting recolored image: %v", err)
	}
}
