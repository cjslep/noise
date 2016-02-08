# Noise library

## Installation

```
go get github.com/cjslep/noise
```

## About

A small and simple library for generating noise. This library provides Perlin
and simplex two dimensional noise generators. It also provides a form of Perlin
noise with Catmull-Rom spline interpolation, although it displays visual
artifacts in the form of faint gridlines.

The library also provides functionality for noise composed of octaves using a
persistence value. The octaves have constant gain and lacunarity.

The following images are generated when running `go test`:

Perlin Noise:
![Perlin Noise](perlin_test.png)

Simplex Noise:
![Simplex Noise](simplex_test.png)

Perlin Noise with Catmull-Rom Spline Interpolation:
![Perlin Noise with Catmull-Rom Spline Interpolation](perlin_spline_test.png)

Perlin Pink Octave Noise:
![Perlin Pink Octave Noise](octave_perlin_test.png)

Simplex Pink Octave Noise:
![Simplex Pink Octave Noise](octave_simplex_test.png)

Perlin Pink Octave Noise with Catmull-Rom Spline Interpolation:
![Perlin Pink Octave Noise with Catmull-Rom Spline Interpolation](octave_perlin_spline_test.png)

## How To Use

All noise generators implement the `Noiser` interface. Please see the
[documentation](https://godoc.org/github.com/cjslep/noise) for details.

## License

GPLv3
