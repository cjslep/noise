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

var _ Noiser = &Simplex{}

// Simplex implements simplex noise generation in two dimensions.
type Simplex struct {
	rng  *rand.Rand
	hash []int
}

// NewSimplex returns a new source of simplex noise. Identical seeds generate
// identical noise for the same inputs.
func NewSimplex(seed int64) *Simplex {
	s := &Simplex{
		rng:  rand.New(rand.NewSource(seed)),
		hash: make([]int, 0, hashSize2D*2),
	}
	s.init()
	return s
}

// init creates the gradient hash lookup used for noise generation.
func (s *Simplex) init() {
	for i := 0; i < hashSize2D; i++ {
		s.hash = append(s.hash, s.rng.Intn(hashSize2D))
	}
	s.hash = append(s.hash, s.hash...)
}

// Noise creates two-dimensional simplex noise.
func (s *Simplex) Noise(x, y float64) float64 {
	commonFactorUnskew := (x + y) * coordTransformToUnskew(2)
	simplexX := intFloor(x + commonFactorUnskew)
	simplexY := intFloor(y + commonFactorUnskew)

	skewFactor := coordTransformToSkew(2)
	commonFactorSkew := float64(simplexX+simplexY) * skewFactor
	skewSimplexX := float64(simplexX) + commonFactorSkew
	skewSimplexY := float64(simplexY) + commonFactorSkew

	firstX := x - skewSimplexX
	firstY := y - skewSimplexY

	unitX := 0
	unitY := 0
	if firstX > firstY {
		unitX = 1 // Lower Simplex
	} else {
		unitY = 1 // Upper Simplex
	}

	middleX := firstX - float64(unitX) - skewFactor
	middleY := firstY - float64(unitY) - skewFactor
	lastX := firstX - 1 - 2*skewFactor
	lastY := firstY - 1 - 2*skewFactor

	simplexX = intMod(simplexX, hashSize2D)
	for simplexX < 0 {
		simplexX += hashSize2D
	}
	simplexY = intMod(simplexY, hashSize2D)
	for simplexY < 0 {
		simplexY += hashSize2D
	}

	grad0 := intMod(s.hash[simplexX+s.hash[simplexY]], len(gradient2D))
	grad1 := intMod(s.hash[simplexX+unitX+s.hash[simplexY+unitY]], len(gradient2D))
	grad2 := intMod(s.hash[simplexX+1+s.hash[simplexY+1]], len(gradient2D))

	t0 := 0.5 - firstX*firstX - firstY*firstY
	t1 := 0.5 - middleX*middleX - middleY*middleY
	t2 := 0.5 - lastX*lastX - lastY*lastY

	contrib0 := 0.0
	contrib1 := 0.0
	contrib2 := 0.0

	if t0 > 0 {
		contrib0 = t0 * t0 * t0 * t0 * gradient2D[grad0].DotFloat64(firstX, firstY)
	}
	if t1 > 0 {
		contrib1 = t1 * t1 * t1 * t1 * gradient2D[grad1].DotFloat64(middleX, middleY)
	}
	if t2 > 0 {
		contrib2 = t2 * t2 * t2 * t2 * gradient2D[grad2].DotFloat64(lastX, lastY)
	}
	return contrib0 + contrib1 + contrib2
}
