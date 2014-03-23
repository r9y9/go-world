# GO-WORLD
-------------

Go port to WORLD - a high-quality speech analysis, modification and synthesis system written in C++. The version of WORLD used in this port is 0.1.3.

Check [original site](http://ml.cs.yamanashi.ac.jp/world/) for details about the WORLD. 

## Install WORLD

     git clone git@github.com:r9y9/world.git && cd world
     ./waf configure && ./waf
     sudo ./waf install

or download the original code and make & install it (not tested).

## Install GO-WORLD

    go get github.com/r9y9/go-world

## How to use

Import the package

    import "github.com/r9y9/go-world"

and create a world instance with sample rate [hz] and frame period [ms].

    w := world.New(sampleRate, framePeriod) // e.g. (44100, 5)

and then, do whatever you want with WORLD.

### Dio

    timeAxis, f0 := w.Dio(input, w.NewDioOption()) // default option is used

### StoneMask

    refinedF0 := w.StoneMask(input, timeAxis, f0)

### Star

    spectrogram := w.Star(input, timeAxis, f0)

### Platinum

    residual := w.Platinum(input, timeAxis, f0, spectrogram)

### Synthesis

    synthesized := w.Synthesis(f0, spectrogram, residual, len(input))

...check go codes to know more about GO-WORLD.

## Example

see [example/world_example.go](example/world_example.go)

## Docmentation

- [Godoc](http://godoc.org/github.com/r9y9/go-world)
- [GoWalker](https://gowalker.org/github.com/r9y9/go-world)

## License

Modified-BSD