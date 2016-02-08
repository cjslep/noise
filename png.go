/*
	This file is part of noise.

	noise is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	noise is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with noise.  If not, see <http://www.gnu.org/licenses/>.
*/

package noise

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
)

// WriteGreyImagePng handles writing a Noiser out to a PNG image. It samples
// the noise space in increments specified by noiseSampleDelta. The image is
// generated using multiple goroutines to handle slower methods. It also
// normalizes the image to provide full greyscale values.
//
// The minimum X and minimum Y determine the lower-left point in noise space to
// begin sampling. The number of samples to do in the X and Y directions also
// determines the width and height of the resulting image. Finally, the
// noise sampling delta determines the distance between two sampled points in
// noise space. Setting this parameter to one will yield a black image for
// Perlin Noise.
func WriteGreyImagePng(w io.Writer, noiser Noiser, minX, minY, numberSamplesX, numberSamplesY int, noiseSampleDelta float64) error {
	if numberSamplesX <= 0 || numberSamplesY <= 0 {
		return errors.New("invalid dimensions")
	}

	grays := make([][]float64, numberSamplesY)
	noiseY := float64(minY)
	done := make(chan struct{ MinValue, MaxValue float64 })
	for y := 0; y < numberSamplesY; y++ {
		go func(y int, noiseY float64) {
			grays[y] = make([]float64, numberSamplesX)
			noiseX := float64(minX)
			minValue := 0.0
			maxValue := 0.0
			for x := 0; x < numberSamplesX; x++ {
				grays[y][x] = noiser.Noise(noiseX, noiseY)
				if x == 0 && y == 0 {
					minValue = grays[y][x]
					maxValue = minValue
				} else if grays[y][x] > maxValue {
					maxValue = grays[y][x]
				} else if grays[y][x] < minValue {
					minValue = grays[y][x]
				}
				noiseX += noiseSampleDelta
			}
			done <- struct{ MinValue, MaxValue float64 }{minValue, maxValue}
		}(y, noiseY)
		noiseY += noiseSampleDelta
	}
	minValue := 0.0
	maxValue := 0.0
	for y := 0; y < numberSamplesY; y++ {
		mm := <-done
		if y == 0 {
			minValue = mm.MinValue
			maxValue = mm.MaxValue
		} else if mm.MinValue < minValue {
			minValue = mm.MinValue
		} else if mm.MaxValue > maxValue {
			maxValue = mm.MaxValue
		}
	}
	img := image.NewGray16(image.Rect(0, 0, numberSamplesX, numberSamplesY))

	for y := 0; y < numberSamplesY; y++ {
		for x := 0; x < numberSamplesX; x++ {
			frac := (grays[y][x] - minValue) / (maxValue - minValue)
			uIntVal := uint16(linearInterpolation(0, 65535, frac))
			shade := color.Gray16{uIntVal}
			img.SetGray16(x, numberSamplesY-y-1, shade)
		}
	}

	return png.Encode(w, img)
}
