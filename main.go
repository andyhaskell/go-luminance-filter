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

	"github.com/spf13/cobra"
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

var cmd = &cobra.Command{
	Use:   "go-luminance-filter raw.jpg -o recolored.jpg -t 0,50 -c 80C9AC,221F20",
	Short: "Edit images by recoloring pixels by luminance",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func init() {
	// [TODO] Add flag params
}

// shouldOutputToFile checks whether we should output our recolored image to
// the filepath passed in. If the file exists, then the user is prompted to
// decide whether to replace it. If an unexpected error occurs, then the error
// is fatally logged.
func shouldOutputToFile(path string) bool {
	if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
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
			return false
		}
	}
	return true
}

func run(cmd *cobra.Command, args []string) {
	inputPath := args[0]
	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("error loading %s: %v", inputPath, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("error decoding file to an Image: %v", err)
	}

	if !shouldOutputToFile("recolored.jpg") {
		return
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

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error running luminance filter: %v", err)
	}
}
