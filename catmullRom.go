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
)

// catmullRom represents a generic Catmull-Rom spline.
type catmullRom struct {
	alpha          float64
	p0, p1, p2, p3 point2D
	t1, t2, t3     float64
}

// newCatmullRom creates a generic Catmull-Rom spline. If the alpha value is
// 0.5, then the spline is a centripetal one. If alpha is 1.0, then the spline
// is a chordal one. Finally, if alpha is 0, the spline is a uniform one.
func newCatmullRom(alpha float64, p0, p1, p2, p3 point2D) *catmullRom {
	c := &catmullRom{
		alpha: alpha,
		p0:    p0,
		p1:    p1,
		p2:    p2,
		p3:    p3,
	}
	c.t1 = c.nextT(c.p0, c.p1)
	c.t2 = c.nextT(c.p1, c.p2) + c.t1
	c.t3 = c.nextT(c.p2, c.p3) + c.t2
	return c
}

// newCentripetal creates a centripetal Catmull-Rom spline.
func newCentripetal(p0, p1, p2, p3 point2D) *catmullRom {
	return newCatmullRom(0.5, p0, p1, p2, p3)
}

// newUniform creates a uniform Catmull-Rom spline.
func newUniform(p0, p1, p2, p3 point2D) *catmullRom {
	return newCatmullRom(0, p0, p1, p2, p3)
}

// newChordal creates a uniform Catmull-Rom spline.
func newChordal(p0, p1, p2, p3 point2D) *catmullRom {
	return newCatmullRom(1, p0, p1, p2, p3)
}

// LowerT returns the lower parametric value for the spline curve.
func (c *catmullRom) LowerT() float64 {
	return c.t1
}

// UpperT returns the upper parametric value for the spline curve.
func (c *catmullRom) UpperT() float64 {
	return c.t2
}

// At returns the spline's point at the parametric value. Valid values for t
// are between LowerT and UpperT inclusive. Values of t outside this range
// result in undefined behavior.
func (c *catmullRom) At(t float64) point2D {
	a1Lower, a1Upper := c.lowerUpper(0, c.t1, t)
	a2Lower, a2Upper := c.lowerUpper(c.t1, c.t2, t)
	a3Lower, a3Upper := c.lowerUpper(c.t2, c.t3, t)

	a1 := c.p0.Scale(a1Lower).Add(c.p1.Scale(a1Upper))
	a2 := c.p1.Scale(a2Lower).Add(c.p2.Scale(a2Upper))
	a3 := c.p2.Scale(a3Lower).Add(c.p3.Scale(a3Upper))

	b1Lower, b1Upper := c.lowerUpper(0, c.t2, t)
	b2Lower, b2Upper := c.lowerUpper(c.t1, c.t3, t)

	b1 := a1.Scale(b1Lower).Add(a2.Scale(b1Upper))
	b2 := a2.Scale(b2Lower).Add(a3.Scale(b2Upper))

	cLower, cUpper := c.lowerUpper(c.t1, c.t2, t)

	return b1.Scale(cLower).Add(b2.Scale(cUpper))
}

// lowerUpper calculates the coefficients necessary for generating a
// Catmull-Rom spline from parametric values.
func (c *catmullRom) lowerUpper(tLower, tUpper, t float64) (lower float64, upper float64) {
	lower = (tUpper - t) / (tUpper - tLower)
	upper = (t - tLower) / (tUpper - tLower)
	return
}

// nextT calculates the difference between the next parametric values for the
// points.
func (c *catmullRom) nextT(p0, p1 point2D) float64 {
	dx := p1.X - p0.X
	dx2 := dx * dx
	dy := p1.Y - p0.Y
	dy2 := dy * dy
	return math.Pow(math.Sqrt(dx2+dy2), c.alpha)
}
