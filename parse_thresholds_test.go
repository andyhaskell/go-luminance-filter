package main

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseThresholdPairValid(t *testing.T) {
	th, err := parseThresholdPair("0", "80C9AC")
	require.NoError(t, err)
	assert.InDelta(t, 0, th.luminancePercent, 0.01)
	assertSameColor(t, blueGreen, th.color)

	// try parsing a valid threshold pair with punctuations
	th, err = parseThresholdPair("100%%", "#80C9AC")
	require.NoError(t, err)
	assert.InDelta(t, 100, th.luminancePercent, 0.01)
	assertSameColor(t, blueGreen, th.color)

	// try parsing a color whose hex value would also be valid base 10
	th, err = parseThresholdPair("100%%", "#112358")
	require.NoError(t, err)
	assert.InDelta(t, 100, th.luminancePercent, 0.01)
	assertSameColor(t, color.RGBA{R: 0x11, G: 0x23, B: 0x58}, th.color)

	// try parsing a color whose hex value has leading zeros
	th, err = parseThresholdPair("100%%", "#0000FF")
	require.NoError(t, err)
	assert.InDelta(t, 100, th.luminancePercent, 0.01)
	assertSameColor(t, color.RGBA{R: 0, G: 0, B: 255}, th.color)
}

func TestParseThresholdPairInvalid(t *testing.T) {
	// luminance out of range
	_, err := parseThresholdPair("-1", "80C9AC")
	require.Error(t, err)

	_, err = parseThresholdPair("101", "80C9AC")
	require.Error(t, err)

	// invalid punctuation
	_, err = parseThresholdPair("50", "?80C9AC")
	require.Error(t, err)

	_, err = parseThresholdPair("50?", "80C9AC")
	require.Error(t, err)

	// invalid hexadecimal color
	_, err = parseThresholdPair("50", "#GGGGGG")
	require.Error(t, err)

	// negative hexadecimal color
	_, err = parseThresholdPair("50", "-FFFFFF")
	require.Error(t, err)

	// hexadecimal color that's too large
	_, err = parseThresholdPair("50", "FFFFFFF")
	require.Error(t, err)

	// hexadecimal color that's too large but has leading zeros
	_, err = parseThresholdPair("50", "00FFFFF")
	require.Error(t, err)

	// hexadecimal color that's too small
	_, err = parseThresholdPair("50", "00FFF")
	require.Error(t, err)
}

func TestParseValidThresholds(t *testing.T) {
	// two valid luminance thresholds
	thresholds, err := parseThresholds("0,80C9AC,50,221F20")
	require.NoError(t, err)
	require.Len(t, thresholds, 2)

	assert.InDelta(t, 0, thresholds[0].luminancePercent, 0.01)
	assertSameColor(t, blueGreen, thresholds[0].color)
	assert.InDelta(t, 50, thresholds[1].luminancePercent, 0.01)
	assertSameColor(t, charcoal, thresholds[1].color)

	// three valid luminance thresholds
	thresholds, err = parseThresholds("0%,#80C9AC,50%,#221F20,75%,#FF00FF")
	require.NoError(t, err)
	require.Len(t, thresholds, 3)

	assert.InDelta(t, 0, thresholds[0].luminancePercent, 0.01)
	assertSameColor(t, blueGreen, thresholds[0].color)
	assert.InDelta(t, 50, thresholds[1].luminancePercent, 0.01)
	assertSameColor(t, charcoal, thresholds[1].color)
	assert.InDelta(t, 75, thresholds[2].luminancePercent, 0.01)
	assertSameColor(t, color.RGBA{R: 255, G: 0, B: 255}, thresholds[2].color)

	// luminance thresholds are sorted after being parsed
	thresholds, err = parseThresholds("50%,#221F20,0%,#80C9AC,75%,#FF00FF")
	require.NoError(t, err)
	require.Len(t, thresholds, 3)

	assert.InDelta(t, 0, thresholds[0].luminancePercent, 0.01)
	assertSameColor(t, blueGreen, thresholds[0].color)
	assert.InDelta(t, 50, thresholds[1].luminancePercent, 0.01)
	assertSameColor(t, charcoal, thresholds[1].color)
	assert.InDelta(t, 75, thresholds[2].luminancePercent, 0.01)
	assertSameColor(t, color.RGBA{R: 255, G: 0, B: 255}, thresholds[2].color)
}

func TestParseInvalidThresholds(t *testing.T) {
	// no thresholds
	_, err := parseThresholds("")
	assert.Error(t, err)

	// only one threshold
	_, err = parseThresholds("0,80C9AC")
	assert.Error(t, err)

	// odd number of strings
	_, err = parseThresholds("0,80C9AC,50,221F20,75")
	assert.Error(t, err)

	// invalid luminance
	_, err = parseThresholds("0,80C9AC,50,221F20,135,FF00FF")
	assert.Error(t, err)

	// invalid hex color
	_, err = parseThresholds("0,80C9AC,50,221F2G")
	assert.Error(t, err)

	// no zero luminance percent
	_, err = parseThresholds("50,80C9AC,75,221F20")
	assert.Error(t, err)

	// duplicate luminance percents
	_, err = parseThresholds("0,FFFFFF,50,80C9AC,50,221F20")
	assert.Error(t, err)
}
