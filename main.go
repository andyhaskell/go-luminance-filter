package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

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
	_ = img // TODO: Edit the image
	log.Println("We got an image!")
}
