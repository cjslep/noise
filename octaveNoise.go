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

var _ Noiser = &OctaveNoise{}

// OctaveNoise uses other Noisers to create more noises composed on one another
// using constant gain and lacunarity.
type OctaveNoise struct {
	persistence float64
	octaves     []Noiser
}

// NewOctaveNoise creates an octave-based noise with the given persistence.
// A typical persistence value is around one-half.
func NewOctaveNoise(persistence float64) *OctaveNoise {
	return &OctaveNoise{
		persistence: persistence,
		octaves:     make([]Noiser, 0, 2),
	}
}

// AddOctave adds the Noiser. Persistence is applied in the order they are
// added. Noisers added later will have a higher sampling frequency but a lower
// amplitude if the persistence was less than one.
func (o *OctaveNoise) AddOctave(n Noiser) {
	o.octaves = append(o.octaves, n)
}

// Noise generates noise for the given input.
func (o *OctaveNoise) Noise(x, y float64) float64 {
	frequency := 1.0
	amplitude := 1.0
	result := 0.0
	cumulativeAmp := 0.0
	for _, octave := range o.octaves {
		result += octave.Noise(x*frequency, y*frequency) * amplitude
		frequency *= 2
		cumulativeAmp += amplitude
		amplitude *= o.persistence
	}
	return result
}
