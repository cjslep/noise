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

// catmullRomCached caches points along the Catmull-Rom spline curve for later
// interpolation. This allows using a spline without requiring a parametric
// parameter. This implementation does not handle loops, vertical splines, nor
// ones that could not otherwise be considered a function.
type catmullRomCached struct {
	cache []point2D
}

// newCachedCatmullRom creates a cache for a generic Catmull-Rom spline.
func newCachedCatmullRom(nPoints int, alpha float64, p0, p1, p2, p3 point2D) *catmullRomCached {
	c := &catmullRomCached{
		make([]point2D, 0, nPoints),
	}
	c.init(newCatmullRom(alpha, p0, p1, p2, p3), nPoints)
	return c
}

// newCentripetalCached creates a cache for a centripetal Catmull-Rom spline.
func newCentripetalCached(nPoints int, p0, p1, p2, p3 point2D) *catmullRomCached {
	return newCachedCatmullRom(nPoints, 0.5, p0, p1, p2, p3)
}

// newUniformCached creates a cache for a uniform Catmull-Rom spline.
func newUniformCached(nPoints int, p0, p1, p2, p3 point2D) *catmullRomCached {
	return newCachedCatmullRom(nPoints, 0, p0, p1, p2, p3)
}

// newChordalCached creates a cache for a chordal Catmull-Rom spline.
func newChordalCached(nPoints int, p0, p1, p2, p3 point2D) *catmullRomCached {
	return newCachedCatmullRom(nPoints, 1, p0, p1, p2, p3)
}

// init creates the cache for a Catmull-Rom spline.
func (c *catmullRomCached) init(curve *catmullRom, nPoints int) {
	if nPoints <= 2 {
		nPoints = 3
	}
	delta := (curve.UpperT() - curve.LowerT()) / float64(nPoints-1)
	for t := curve.LowerT(); t <= curve.UpperT(); t += delta {
		c.cache = append(c.cache, curve.At(t))
	}
}

// InterpolateX uses a binary search to estimate the y-value for the given
// point. It is not guaranteed to exist if the x-value is outside p1 or p2's
// x-values.
func (c *catmullRomCached) InterpolateX(x float64) float64 {
	if len(c.cache) == 0 {
		return 0
	}
	min := 0
	max := len(c.cache) - 1
	i := (max + min) / 2
	for min <= max && i > 0 && i < len(c.cache)-1 {
		if c.cache[i].X <= x && c.cache[i+1].X > x {
			break
		} else if c.cache[i].X > x {
			max = i - 1
		} else if c.cache[i+1].X <= x {
			min = i + 1
		}
		i = (max + min) / 2
	}
	if i >= len(c.cache)-1 {
		return c.cache[len(c.cache)-1].Y
	} else if i < 0 {
		return c.cache[0].Y
	}
	return c.cache[i].LinearInterpolation(c.cache[i+1], x)
}
