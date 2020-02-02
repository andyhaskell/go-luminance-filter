package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func recolor(img image.Image) {
	b := img.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := img.At(x, y)

			// we have the color of the pixel at (x, y). Now we
			// just need to figure out how bright that pixel is,
			// and recolor it accordingly!
			_ = color
		}
	}
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

	// we need this to know the coordinates of the top-left pixel, which isn't
	// necessarily (0, 0)
	b := img.Bounds()

	color := img.At(b.Min.X, b.Min.Y)
	log.Printf("%#v\n", color)
}
