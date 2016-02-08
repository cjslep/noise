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

// hashSize2D is the size of the gradient lookup table used in noise
// generation.
const hashSize2D = 2048 * 4

// unitCircleDelta is used to generate a uniformly-spaced set of vectors that
// lie on the unit circle.
var unitCircleDelta float64 = math.Pi / 6

// gradient2D is a set of evenly-spaced vectors that lie on the unit circle.
// Lookup tables randomly map points in space to these as gradient vectors.
var gradient2D []point2D = []point2D{
	{1, 0},
	{math.Cos(unitCircleDelta), math.Sin(unitCircleDelta)},
	{math.Cos(unitCircleDelta * 2), math.Sin(unitCircleDelta * 2)},
	{math.Cos(unitCircleDelta * 3), math.Sin(unitCircleDelta * 3)},
	{math.Cos(unitCircleDelta * 4), math.Sin(unitCircleDelta * 4)},
	{math.Cos(unitCircleDelta * 5), math.Sin(unitCircleDelta * 5)},
	{math.Cos(unitCircleDelta * 6), math.Sin(unitCircleDelta * 6)},
	{math.Cos(unitCircleDelta * 7), math.Sin(unitCircleDelta * 7)},
	{math.Cos(unitCircleDelta * 8), math.Sin(unitCircleDelta * 8)},
	{math.Cos(unitCircleDelta * 9), math.Sin(unitCircleDelta * 9)},
	{math.Cos(unitCircleDelta * 10), math.Sin(unitCircleDelta * 10)},
	{math.Cos(unitCircleDelta * 11), math.Sin(unitCircleDelta * 11)},
}

var origin = point2D{0, 0}

// point2D represents a two-dimensional point.
type point2D struct {
	X, Y float64
}

// DotFloat64 performs an inner product.
func (p point2D) DotFloat64(x, y float64) float64 {
	return p.X*x + p.Y*y
}

// Dot performs an inner product.
func (p point2D) Dot(o point2D) float64 {
	return p.DotFloat64(o.X, o.Y)
}

// Scale applies a scalar to both x and y values of this point.
func (p point2D) Scale(s float64) point2D {
	return point2D{p.X * s, p.Y * s}
}

// Add adds a point to this one.
func (p point2D) Add(o point2D) point2D {
	return point2D{p.X + o.X, p.Y + o.Y}
}

// LinearInterpolation performs linear interpolation between two points.
func (p point2D) LinearInterpolation(o point2D, x float64) float64 {
	return p.Y + (o.Y-p.Y)*(x-p.X)/(o.X-p.X)
}

// Mag returns the magnitude of this point as a vector from (0,0).
func (p point2D) Mag() float64 {
	return distance0(p.X, p.Y)
}

// intFloor is a helper for converting a float to an int after flooring.
func intFloor(x float64) int {
	return int(math.Floor(x))
}

// fader is a second-derivative-continuous fading function for Perlin noise.
func fader(t float64) float64 {
	return t * t * t * (10 + t*(-15+t*6))
}

// linearInterpolation performs linear interpolation using a fractional t value
// in the range 0 <= t <= 1.
func linearInterpolation(x0, x1, t float64) float64 {
	return (1-t)*x0 + t*x1
}

// coordTransformToUnskew calculates the skew value for simplex noise when
// transforming to unskewed coordinates.
func coordTransformToUnskew(dims int) float64 {
	return (math.Sqrt(float64(dims+1)) - 1) / float64(dims)
}

// coordTransformToSkew calculates the skew value for simplex noise when
// transforming to skewed coordinates.
func coordTransformToSkew(dims int) float64 {
	return ((1 / math.Sqrt(float64(dims+1))) - 1) / float64(dims)
}

// intMod returns the mod value between two integers as an integer.
func intMod(a, b int) int {
	return int(math.Mod(float64(a), float64(b)))
}

// distance determines the distance between two points.
func distance(x0, y0, x1, y1 float64) float64 {
	dx := x1 - x0
	dy := y1 - y0
	return math.Sqrt(dx*dx + dy*dy)
}

// distance0 determines the distance from the origin.
func distance0(x, y float64) float64 {
	return distance(0, 0, x, y)
}
