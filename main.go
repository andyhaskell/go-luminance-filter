package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
