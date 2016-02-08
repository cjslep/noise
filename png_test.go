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
	"bytes"
	"io/ioutil"
	"testing"
)

const (
	seed            = 42
	startCorner     = -100
	imageDimension  = 200
	imageSize       = imageDimension * imageDimension
	splineCacheSize = 2
	sampleStep      = 0.137
)

func testWithNoiser(t *testing.T, generator Noiser, filename string) {
	buffer := bytes.NewBuffer(make([]byte, 0, imageSize))
	if err := WriteGreyImagePng(buffer, generator, startCorner, startCorner, imageDimension, imageDimension, sampleStep); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filename, buffer.Bytes(), 0666); err != nil {
		t.Fatal(err)
	}
}

func TestPerlin(t *testing.T) {
	testWithNoiser(t, NewPerlin(seed), "perlin_test.png")
}

func TestPerlinCatmullRom(t *testing.T) {
	testWithNoiser(t, NewPerlinCatmullRom(splineCacheSize, seed), "perlin_spline_test.png")
}

func TestSimplex(t *testing.T) {
	testWithNoiser(t, NewSimplex(seed), "simplex_test.png")
}

func TestPinkPerlinOctave(t *testing.T) {
	octaveGenerator := NewOctaveNoise(0.5)
	for i := 0; i < 8; i++ {
		octaveGenerator.AddOctave(NewPerlin(seed))
	}
	testWithNoiser(t, octaveGenerator, "octave_perlin_test.png")
}

func TestPinkPerlinCatmullRomOctave(t *testing.T) {
	octaveGenerator := NewOctaveNoise(0.5)
	for i := 0; i < 8; i++ {
		octaveGenerator.AddOctave(NewPerlinCatmullRom(splineCacheSize, seed))
	}
	testWithNoiser(t, octaveGenerator, "octave_perlin_spline_test.png")
}

func TestPinkSimplexOctave(t *testing.T) {
	octaveGenerator := NewOctaveNoise(0.5)
	for i := 0; i < 8; i++ {
		octaveGenerator.AddOctave(NewSimplex(seed))
	}
	testWithNoiser(t, octaveGenerator, "octave_simplex_test.png")
}
