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
	Use:   "go-luminance-filter raw.jpg -o recolored.jpg -t 0,000000,50,FFFFFF",
	Short: "Edit images by recoloring pixels by luminance",
	Args:  cobra.ExactArgs(1),
	Run:   run,
}

func init() {
	cmd.Flags().StringP(
		"output",
		"o",
		"recolored.jpg",
		"name of the file to output recolored image to.",
	)

	// [TODO] This felt clunky to explain how to use. Could there be a
	// more intuitive way to parse luminance thresholds than in this
	// format, or a simpler explanation of this CLI arg?
	cmd.Flags().StringP(
		"thresholds",
		"t",
		"0,000000,50,FFFFFF",
		"comma-separated list of pairs of luminance percentages and hexadecimal colors, "+
			"in the pattern 'luminancePercentage,color,luminancePercentage,color'. "+
			"If a given pixel's luminance is above the n-th percentage, but below "+
			"the next percentage in the list, then that pixel will be recolored to "+
			"the n-th color. For example in the argument \"0,000000,50,FFFFFF\" a "+
			"pixel with 25% luminance would be recolored to black (#000000) since "+
			"it's above 0% luminance, while a pixel with 60% luminance would be "+
			"recolored to white (#FFFFFF). Must have one luminance threshold of 0"+
			"as a catch-all value.",
	)
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

	outFilepath, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Fatalf("error getting output file path: %v", err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("error decoding file to an Image: %v", err)
	}

	if !shouldOutputToFile(outFilepath) {
		return
	}

	out, err := os.Create(outFilepath)
	if err != nil {
		log.Fatalf("error creating output file recolored.jpg: %v", err)
	}
	defer out.Close()

	thresholdsArg, err := cmd.Flags().GetString("thresholds")
	if err != nil {
		log.Fatalf("error getting luminance thresholds CLI arg: %v", err)
	}
	thresholds, err := parseThresholds(thresholdsArg)
	if err != nil {
		log.Fatalf("error parsing luminance thresholds: %v", err)
	}

	recolored := recolor(img, thresholds)
	if err := jpeg.Encode(out, recolored, nil); err != nil {
		log.Fatalf("error outputting recolored image: %v", err)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error running luminance filter: %v", err)
	}
}
