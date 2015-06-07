# GO-WORLD

[![Build Status](https://travis-ci.org/r9y9/go-world.svg?branch=master)](https://travis-ci.org/r9y9/go-world)
[![GoDoc](https://godoc.org/github.com/r9y9/go-world?status.svg)](https://godoc.org/github.com/r9y9/go-world)

GO-WORLD is a Go port to WORLD - a high-quality speech analysis, modification and synthesis system. WORLD provides a way to decompose a speech signal into:

- Fundamental frequency (F0)
- spectral envelope
- aperiodicity

and re-synthesize a speech signal from these paramters. See [here](http://ml.cs.yamanashi.ac.jp/world/english/index.html) for the original WORLD.

## Supported Platforms

- Linux
- Mac OS X

Note that the original WORLD works in windows as well. In order to use WORLD in windows, you have to build WORLD yourself since currently we don't have a installation script.

## Installation

### Binary dependency

First you need to install WORLD as a shared library:

```bash
git clone https://github.com/r9y9/WORLD.git && cd world
git checkout v0.2.1-2
./waf configure && ./waf
sudo ./waf install
```

### GO-WORLD

```bash
go get github.com/r9y9/go-world
```

complete!

## Usage

Import the package

```go
import "github.com/r9y9/go-world"
```

and create a world instance with sample rate [hz] and frame period [ms].

```go
w := world.New(sampleRate, framePeriod) // e.g. (44100, 5)
```

and then you can do whatever you want with WORLD.

### F0 estimation and refinement

#### Dio

```go
timeAxis, f0 := w.Dio(input, w.NewDioOption()) // default option is used
```

![](https://raw.githubusercontent.com/r9y9/WORLD.jl/master/examples/f0_by_dio.png)

#### StoneMask

```go
refinedF0 := w.StoneMask(input, timeAxis, f0)
```

![](https://raw.githubusercontent.com/r9y9/WORLD.jl/master/examples/f0_refinement.png)

### Spectral envelope estimation by CheapTrick

```go
spectrogram := w.CheapTrick(input, timeAxis, f0)
```

![](https://raw.githubusercontent.com/r9y9/WORLD.jl/master/examples/envelope_by_cheaptrick.png)

### Aperiodicity ratio estimation by D4C

```go
apiriodicity := w.D4C(input, timeAxis, f0)
```

![](https://raw.githubusercontent.com/r9y9/WORLD.jl/master/examples/aperiodicity_by_d4c.png)

### Synthesis

```go
w.Synthesis(f0, spectrogram, apiriodicity, len(input))
```

![](https://raw.githubusercontent.com/r9y9/WORLD.jl/master/examples/synthesis.png)

## Example

see [example/world_example.go](example/world_example.go)

## License

Modified-BSD
