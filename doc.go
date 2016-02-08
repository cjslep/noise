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

/*
	Package noise is a small library for generating noise. It contains
	implementations for Perlin and Simplex noise in two dimensions. All
	noises use a seed and a lookup table to create consistent outputs for
	the same seed, even from different noise generators. All noises
	implement the Noiser interface.

	Octave noise is also provided, although must be composed of other
	noises to produce the resulting smoothed noise.

	This package provides two different kinds of Perlin noise: A very fast
	4-point interpolation using a simple fading function that has a zero
	second derivative at the end of the interpolation ranges ensuring a
	smoother result, and a slow 16-point interpolation using five
	Catmull-Rom splines. The former generates noise that can sometimes
	betray the regular rectangular pattern of the gradient used to
	interpolate, while the latter can elminate these artifacts.
	Furthermore, the noise sampled at coordinate points will all be zero.
	It is knowwn that the Catmull-Rom spline interpolation method contains
	gridlike artifacts.

		// Perlin noise with 4-point interpolation
		perlinGenerator := noise.NewPerlin(1)
		val := perlinGenerator.Noise(0.5, 0.5)
		// Perlin noise with 16-point interpolation
		catmullRomPerlinGenerator := noise.NewPerlinCatmullRom(1)
		val := catmullRomPerlinGenerator.Noise(0.5, 0.5)

	Simplex noise uses simplexes to efficiently interpolate noise instead
	of a regular rectangular grid. This can result in a different skewed
	repetitive pattern along the simplexes used for interpolation.

		// Simplex noise with simplex interpolation
		simplexGenerator := noise.NewSimplex(1)
		val := simplexGenerator.Noise(0.5, 0.5)

	Octave noise combines several different noises with increasing
	persistence which diminishes the amplitude of subsequently-added noise
	and widens its sampling frequency. The gain and lacunarity are both
	held constant. Octave noise composed of the same constituent noise with
	the same seeds with a persistence of one-half is sometimes referred to
	as pink noise, fractional noise, or fractal noise.

		// Pink noise with two octaves.
		pinkNoiseGenerator := noise.NewOctaveNoise(0.5)
		pinkNoiseGenerator.AddOctave(noise.NewPerlin(1))
		pinkNoiseGenerator.AddOctave(noise.NewPerlin(1))
		val := pinkNoiseGenerator.Noise(0.5, 0.5)

	A utility function is provided to help write out Noisers to greyscale
	PNG images. It uses goroutines to parallelize sampling due to the
	slowness of some noise generation methods.
*/
package noise
