# gosfft
![Build status](https://travis-ci.org/davidkleiven/gosfft.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/davidkleiven/gosfft/badge.svg?branch=master)](https://coveralls.io/github/davidkleiven/gosfft?branch=master)

Simple FFT (SFFT) is a simple FFT library that is based on Gonum's FFT routine. It implements a simple interface for 1D, 2D and 3D transforms.

# Examples

Below is a selection of examples whoen

## [1D fourier transform](examples/fft1d/main.go)

Fourier transform of a square pulse

![Signal 1D](figs/signal1D.png) ![Fourier Transform](figs/fft1D.png)

## [2D fourier transform](examples/fft2d/main.go)

Fourier transform of a square

![Signal 2D](figs/img.png) ![Fourier Transform](figs/fft2D.png)
