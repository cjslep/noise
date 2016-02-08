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
	"math"
	"math/rand"
)

var _ Noiser = &PerlinCatmullRom{}

// PerlinCatmullRom creates Perlin noise using centripetal Catmull-Rom spline
// interpolation. This is slow.
//
// Warning: this contains visual gridline artifacts.
type PerlinCatmullRom struct {
	rng             *rand.Rand
	splineCacheSize int
	hash            []int
}

// Constructs a new source of noise using a centripetal Catmull-Rom spline
// interpolation, with splines having the given number of points. The seed is
// used to ensure idential PerlinCatmullRoms will return the same noise values
// for the same inputs.
//
// Warning: this contains visual gridline artifacts.
func NewPerlinCatmullRom(splineCacheSize int, seed int64) *PerlinCatmullRom {
	s := &PerlinCatmullRom{
		rng:             rand.New(rand.NewSource(seed)),
		splineCacheSize: splineCacheSize,
		hash:            make([]int, 0, hashSize2D*2),
	}
	s.init()
	return s
}

// init creates the Perlin noise gradient hash.
func (s *PerlinCatmullRom) init() {
	for i := 0; i < hashSize2D; i++ {
		s.hash = append(s.hash, s.rng.Intn(hashSize2D))
	}
	s.hash = append(s.hash, s.hash...)
}

// hashIdx handles wrapping out of bounds indices when doing hash lookups.
func (s *PerlinCatmullRom) hashIdx(idx int) int {
	return s.hash[intMod(idx, len(s.hash))]
}

// Noise creates Perlin noise using Catmull-Rom spline interpolations, which is
// considerably slower than the simple variant of Perlin noise.
func (s *PerlinCatmullRom) Noise(x, y float64) float64 {
	x0 := intFloor(x)
	y0 := intFloor(y)

	relX := x - float64(x0)
	relY := y - float64(y0)

	// Shift (0,0) to be lower left of 4x4 point grid
	x0 = intMod(x0-1, hashSize2D)
	for x0 < 0 {
		x0 += hashSize2D
	}
	y0 = intMod(y0-1, hashSize2D)
	for y0 < 0 {
		y0 += hashSize2D
	}

	grad00 := intMod(s.hashIdx(x0+s.hashIdx(y0)), len(gradient2D))
	grad10 := intMod(s.hashIdx(x0+1+s.hashIdx(y0)), len(gradient2D))
	grad20 := intMod(s.hashIdx(x0+2+s.hashIdx(y0)), len(gradient2D))
	grad30 := intMod(s.hashIdx(x0+3+s.hashIdx(y0)), len(gradient2D))
	grad01 := intMod(s.hashIdx(x0+s.hashIdx(y0+1)), len(gradient2D))
	grad11 := intMod(s.hashIdx(x0+1+s.hashIdx(y0+1)), len(gradient2D))
	grad21 := intMod(s.hashIdx(x0+2+s.hashIdx(y0+1)), len(gradient2D))
	grad31 := intMod(s.hashIdx(x0+3+s.hashIdx(y0+1)), len(gradient2D))
	grad02 := intMod(s.hashIdx(x0+s.hashIdx(y0+2)), len(gradient2D))
	grad12 := intMod(s.hashIdx(x0+1+s.hashIdx(y0+2)), len(gradient2D))
	grad22 := intMod(s.hashIdx(x0+2+s.hashIdx(y0+2)), len(gradient2D))
	grad32 := intMod(s.hashIdx(x0+3+s.hashIdx(y0+2)), len(gradient2D))
	grad03 := intMod(s.hashIdx(x0+s.hashIdx(y0+3)), len(gradient2D))
	grad13 := intMod(s.hashIdx(x0+1+s.hashIdx(y0+3)), len(gradient2D))
	grad23 := intMod(s.hashIdx(x0+2+s.hashIdx(y0+3)), len(gradient2D))
	grad33 := intMod(s.hashIdx(x0+3+s.hashIdx(y0+3)), len(gradient2D))

	pt00 := point2D{relX + 1, relY + 1}
	pt10 := point2D{relX, relY + 1}
	pt20 := point2D{relX - 1, relY + 1}
	pt30 := point2D{relX - 2, relY + 1}
	pt01 := point2D{relX + 1, relY}
	pt11 := point2D{relX, relY}
	pt21 := point2D{relX - 1, relY}
	pt31 := point2D{relX - 2, relY}
	pt02 := point2D{relX + 1, relY - 1}
	pt12 := point2D{relX, relY - 1}
	pt22 := point2D{relX - 1, relY - 1}
	pt32 := point2D{relX - 2, relY - 1}
	pt03 := point2D{relX + 1, relY - 2}
	pt13 := point2D{relX, relY - 2}
	pt23 := point2D{relX - 1, relY - 2}
	pt33 := point2D{relX - 2, relY - 2}

	noise00 := gradient2D[grad00].Dot(pt00)
	noise10 := gradient2D[grad10].Dot(pt10)
	noise20 := gradient2D[grad20].Dot(pt20)
	noise30 := gradient2D[grad30].Dot(pt30)
	noise01 := gradient2D[grad01].Dot(pt01)
	noise11 := gradient2D[grad11].Dot(pt11)
	noise21 := gradient2D[grad21].Dot(pt21)
	noise31 := gradient2D[grad31].Dot(pt31)
	noise02 := gradient2D[grad02].Dot(pt02)
	noise12 := gradient2D[grad12].Dot(pt12)
	noise22 := gradient2D[grad22].Dot(pt22)
	noise32 := gradient2D[grad32].Dot(pt32)
	noise03 := gradient2D[grad03].Dot(pt03)
	noise13 := gradient2D[grad13].Dot(pt13)
	noise23 := gradient2D[grad23].Dot(pt23)
	noise33 := gradient2D[grad33].Dot(pt33)

	c := newCentripetalCached(s.splineCacheSize, point2D{float64(x0), noise00}, point2D{float64(x0) + 1, noise10}, point2D{float64(x0) + 2, noise20}, point2D{float64(x0) + 3, noise30})
	noiseX0 := c.InterpolateX(float64(x0) + 1 + relX)
	noiseX0 = math.Max(-1, math.Min(1, noiseX0))
	c = newCentripetalCached(s.splineCacheSize, point2D{float64(x0), noise01}, point2D{float64(x0) + 1, noise11}, point2D{float64(x0) + 2, noise21}, point2D{float64(x0) + 3, noise31})
	noiseX1 := c.InterpolateX(float64(x0) + 1 + relX)
	noiseX1 = math.Max(-1, math.Min(1, noiseX1))
	c = newCentripetalCached(s.splineCacheSize, point2D{float64(x0), noise02}, point2D{float64(x0) + 1, noise12}, point2D{float64(x0) + 2, noise22}, point2D{float64(x0) + 3, noise32})
	noiseX2 := c.InterpolateX(float64(x0) + 1 + relX)
	noiseX2 = math.Max(-1, math.Min(1, noiseX2))
	c = newCentripetalCached(s.splineCacheSize, point2D{float64(x0), noise03}, point2D{float64(x0) + 1, noise13}, point2D{float64(x0) + 2, noise23}, point2D{float64(x0) + 3, noise33})
	noiseX3 := c.InterpolateX(float64(x0) + 1 + relX)
	noiseX3 = math.Max(-1, math.Min(1, noiseX3))

	c = newCentripetalCached(s.splineCacheSize, point2D{float64(y0), noiseX0}, point2D{float64(y0) + 1, noiseX1}, point2D{float64(y0) + 2, noiseX2}, point2D{float64(y0) + 3, noiseX3})
	noise := c.InterpolateX(float64(y0) + 1 + relY)
	return noise
}
