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
	"math/rand"
)

var _ Noiser = &Perlin{}

// Perlin implements simple Perlin noise using a fading function whose second
// derivative is zero at the interpolation boundaries. This results in a
// smoother visualization.
type Perlin struct {
	rng  *rand.Rand
	hash []int
}

// NewPerlin constructs a new Perlin noise with the given seed. Multiple
// instances constructed from the same seed will return the same noise values
// for the same inputs.
func NewPerlin(seed int64) *Perlin {
	s := &Perlin{
		rng:  rand.New(rand.NewSource(seed)),
		hash: make([]int, 0, hashSize2D*2),
	}
	s.init()
	return s
}

// Constructs the internal hash used to generate the perlin noise.
func (s *Perlin) init() {
	for i := 0; i < hashSize2D; i++ {
		s.hash = append(s.hash, s.rng.Intn(hashSize2D))
	}
	s.hash = append(s.hash, s.hash...)
}

// Noise generates simple Perlin noise.
func (s *Perlin) Noise(x, y float64) float64 {
	x0 := intFloor(x)
	y0 := intFloor(y)

	relX := x - float64(x0)
	relY := y - float64(y0)

	x0 = intMod(x0, hashSize2D)
	for x0 < 0 {
		x0 += hashSize2D
	}
	y0 = intMod(y0, hashSize2D)
	for y0 < 0 {
		y0 += hashSize2D
	}

	grad00 := intMod(s.hash[x0+s.hash[y0]], len(gradient2D))
	grad10 := intMod(s.hash[x0+1+s.hash[y0]], len(gradient2D))
	grad01 := intMod(s.hash[x0+s.hash[y0+1]], len(gradient2D))
	grad11 := intMod(s.hash[x0+1+s.hash[y0+1]], len(gradient2D))

	noise00 := gradient2D[grad00].DotFloat64(relX, relY)
	noise10 := gradient2D[grad10].DotFloat64(relX-1, relY)
	noise01 := gradient2D[grad01].DotFloat64(relX, relY-1)
	noise11 := gradient2D[grad11].DotFloat64(relX-1, relY-1)

	fadeX := fader(relX)
	fadeY := fader(relY)

	noiseX0 := linearInterpolation(noise00, noise10, fadeX)
	noiseX1 := linearInterpolation(noise01, noise11, fadeX)
	return linearInterpolation(noiseX0, noiseX1, fadeY)
}
