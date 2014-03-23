package main

import (
	"flag"
	"fmt"
	"github.com/r9y9/go-dsp/wav"
	"github.com/r9y9/go-world"
	"log"
	"os"
	"time"
)

var defaultDioOption = world.DioOption{
	F0Floor:          80.0,
	F0Ceil:           800.0,
	FramePeriod:      5,
	ChannelsInOctave: 4.0,
	Speed:            2,
}

func worldExample(input []float64, sampleRate int) []float64 {
	w := world.New(sampleRate, defaultDioOption.FramePeriod)

	// 1. Fundamental frequency
	timeAxis, f0 := w.Dio(input, defaultDioOption)

	// 2. Refine F0 estimation result
	f0 = w.StoneMask(input, timeAxis, f0)

	// 3. Spectral envelope
	spectrogram := w.Star(input, timeAxis, f0)

	// 4. Excitation spectrum
	residual := w.Platinum(input, timeAxis, f0, spectrogram)

	// 5. Synthesis
	return w.Synthesis(f0, spectrogram, residual, len(input))
}

// v0.1.3
func worldExampleAp(input []float64, sampleRate int) []float64 {
	w := world.New(sampleRate, defaultDioOption.FramePeriod)

	// 1. Fundamental frequency
	timeAxis, f0 := w.Dio(input, defaultDioOption)

	// 2. Spectral envelope
	spectrogram := w.Star(input, timeAxis, f0)

	// 3. Apiriodiciy
	apiriodicity := w.AperiodicityRatio(input, f0, timeAxis)

	// 4. Synthesis
	return w.SynthesisFromAperiodicity(f0, spectrogram, apiriodicity, len(input))
}

// v0.1.2 (will be removed)
func worldExampleApOld(input []float64, sampleRate int) []float64 {
	w := world.New(sampleRate, defaultDioOption.FramePeriod)

	// 1. Fundamental frequency
	timeAxis, f0 := w.Dio(input, defaultDioOption)

	// 2. Spectral envelope
	spectrogram := w.Star(input, timeAxis, f0)

	// 3. Apiriodiciy
	apiriodicity, targetF0 := w.AperiodicityRatioOld(input, f0)

	// 4. Synthesis
	return w.SynthesisFromAperiodicityOld(f0, spectrogram, apiriodicity, targetF0, len(input))
}

func main() {
	ifilename := flag.String("i", "input.wav", "Input filename")
	ofilename := flag.String("o", "output.wav", "Output filename")
	flag.Parse()

	file, err := os.Open(*ifilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read wav data
	w, werr := wav.ReadWav(file)
	if werr != nil {
		log.Fatal(werr)
	}

	input := w.GetMonoData()
	sampleRate := int(w.SampleRate)

	// WORLD examples
	start := time.Now()
	synthesized := worldExample(input, sampleRate)
	//synthesized := worldExampleAp(input, sampleRate)
	//synthesized := worldExampleApOld(input, sampleRate)

	// Output elapsed timme
	fmt.Println("Finished. Elapsed time:", time.Now().Sub(start))

	// Write to file
	werr = wav.WriteMono(*ofilename, synthesized, w.SampleRate)
	if werr != nil {
		log.Fatal(werr)
	}
	fmt.Println(*ofilename, "is created.")
}
