package main

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	black     = color.RGBA{}
	white     = color.RGBA{R: 255, G: 255, B: 255}
	violet    = color.RGBA{R: 255, G: 0, B: 255}
	charcoal  = color.RGBA{R: 34, G: 31, B: 32, A: 255}
	blueGreen = color.RGBA{R: 128, G: 201, B: 172, A: 255}
)

func TestLuminance(t *testing.T) {
	assert.InDelta(t, 0, luminancePercent(black), 0.01)
	assert.InDelta(t, 100, luminancePercent(white), 0.01)
	assert.InDelta(t, 28.48, luminancePercent(violet), 0.01)
}

// A three by three image of gray colors. With each row and each column
// increasing, the brightness, and therefore luminance, increases as well.
var threeByThreeImg image.Image

func init() {
	img := image.NewGray(image.Rect(0, 0, 3, 3))
	grayValues := [][]uint8{
		{0, 32, 64},
		{96, 128, 160},
		{192, 224, 255},
	}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			img.Set(x, y, color.Gray{grayValues[y][x]})
		}
	}
	threeByThreeImg = img
}

func assertSameColor(t *testing.T, exp, got color.Color) bool {
	expR, expG, expB, _ := exp.RGBA()
	gotR, gotG, gotB, _ := got.RGBA()

	if expR != gotR || expG != gotG || expB != gotB {
		assert.Fail(
			t,
			"did not get expected colors",
			"expected RGB values %d,%d,%d, got %d,%d,%d",
			expR, expG, expB, gotR, gotG, gotB,
		)
		return false
	}
	return true
}

func TestRecolorTwoThresholds(t *testing.T) {
	recolored := recolor(threeByThreeImg, []luminanceThreshold{
		{luminancePercent: 0, color: blueGreen},
		{luminancePercent: 50, color: charcoal},
	})

	expectedColors := [][]color.Color{
		{blueGreen, blueGreen, blueGreen},
		{blueGreen, charcoal, charcoal},
		{charcoal, charcoal, charcoal},
	}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if !assertSameColor(t, expectedColors[y][x], recolored.At(x, y)) {
				t.Logf("color mismatched at pixel (%d, %d)", x, y)
				return
			}
		}
	}
}

func TestRecolorThreeThresholds(t *testing.T) {
	recolored := recolor(threeByThreeImg, []luminanceThreshold{
		{luminancePercent: 0, color: blueGreen},
		{luminancePercent: 50, color: charcoal},
		{luminancePercent: 75, color: violet},
	})

	expectedColors := [][]color.Color{
		{blueGreen, blueGreen, blueGreen},
		{blueGreen, charcoal, charcoal},
		{violet, violet, violet},
	}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if !assertSameColor(t, expectedColors[y][x], recolored.At(x, y)) {
				t.Logf("color mismatched at pixel (%d, %d)", x, y)
				return
			}
		}
	}
}

func TestRecolorPanic(t *testing.T) {
	// This should panic because the pixels darker than 40% luminance don't
	// match any of the luminance thresholds
	assert.Panics(t, func() {
		recolor(threeByThreeImg, []luminanceThreshold{
			{luminancePercent: 40, color: blueGreen},
			{luminancePercent: 50, color: charcoal},
			{luminancePercent: 75, color: violet},
		})
	})
}
