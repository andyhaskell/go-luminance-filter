package main

import (
	"errors"
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"
)

type luminanceThreshold struct {
	luminancePercent int
	color            color.Color
}

type luminanceThresholdSlice []luminanceThreshold

func (l luminanceThresholdSlice) Less(i, j int) bool {
	return l[i].luminancePercent < l[j].luminancePercent
}

func parseThresholdPair(lumStr, colorStr string) (*luminanceThreshold, error) {
	luminancePercent, err := strconv.Atoi(strings.TrimRight(lumStr, "%%"))
	if err != nil {
		return nil, fmt.Errorf(
			"threshold %s was invalid int - error: %v", lumStr, err)
	} else if luminancePercent < 0 || luminancePercent > 100 {
		return nil, fmt.Errorf(
			"luminance threshold percent %d must be valid int 0-100",
			luminancePercent,
		)
	}

	colorHexStr := strings.TrimLeft(colorStr, "#")
	if len(colorHexStr) != 6 {
		return nil, fmt.Errorf(
			"hex color %s did not have 6 hexadecimal digits", colorHexStr)
	}
	colorHex, err := strconv.ParseInt(strings.TrimLeft(colorStr, "#"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"color %s was invalid hexadecimal int - error: %v", colorStr, err)
	} else if colorHex < 0 || colorHex > 0xFFFFFF {
		return nil, fmt.Errorf("color %x must be RGB hex color 000000-FFFFFF", colorHex)
	}

	c := color.RGBA{
		R: uint8((colorHex & 0xFF0000) >> 16),
		G: uint8((colorHex & 0x00FF00) >> 8),
		B: uint8((colorHex & 0x0000FF) >> 0),
		A: 255,
	}
	return &luminanceThreshold{luminancePercent: luminancePercent, color: c}, nil
}

func parseThresholds(thresholdsArg string) ([]luminanceThreshold, error) {
	thresholdStrings := strings.Split(thresholdsArg, ",")
	if len(thresholdStrings)%2 == 1 {
		return nil, errors.New("odd number of strings in thresholds arg")
	} else if len(thresholdStrings) < 4 {
		return nil, errors.New("at least 2 luminance thresholds must be in thresholds arg")
	}

	var hasZero bool
	allThresholds := make([]luminanceThreshold, 0, len(thresholdStrings)/2)
	for i := 0; i < len(thresholdStrings); i += 2 {
		luminance := thresholdStrings[i]
		color := thresholdStrings[i+1]

		th, err := parseThresholdPair(luminance, color)
		if err != nil {
			return nil, err
		}

		if th.luminancePercent == 0 {
			hasZero = true
		}
		for _, threshold := range allThresholds {
			if threshold.luminancePercent == th.luminancePercent {
				return nil, fmt.Errorf(
					"duplicate luminance threshold percent %d",
					th.luminancePercent,
				)
			}
		}

		allThresholds = append(allThresholds, *th)
	}

	if !hasZero {
		return nil, errors.New("no threshold pair for luminance threshold of 0")
	}
	sort.Slice(allThresholds, func(i, j int) bool {
		return allThresholds[i].luminancePercent < allThresholds[j].luminancePercent
	})
	return allThresholds, nil
}
