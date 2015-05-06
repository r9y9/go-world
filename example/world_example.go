package main

import (
	"flag"
	"fmt"
	"github.com/r9y9/go-dsp/wav"
	"github.com/r9y9/go-world"
	"log"
	"os"
	"strings"
	"time"
)

var defaultDioOption = world.DioOption{
	F0Floor:          80.0,
	F0Ceil:           800.0,
	FramePeriod:      5,
	ChannelsInOctave: 4.0,
	Speed:            2,
}

func worldExampleAp(input []float64, sampleRate int) []float64 {
	w := world.New(sampleRate, defaultDioOption.FramePeriod)

	// 1. Fundamental frequency
	timeAxis, f0 := w.Dio(input, defaultDioOption)

	// 2. Refine F0 estimation result
	f0 = w.StoneMask(input, timeAxis, f0)

	// 3. Spectral envelope
	spectrogram := w.CheapTrick(input, timeAxis, f0)

	// 4. Aperiodicity spectrum
	apiriodicity := w.AperiodicityRatio(input, f0, timeAxis)

	// 5. Synthesis from aperiodicity
	return w.SynthesisFromAperiodicity(f0, spectrogram, apiriodicity, len(input))
}

func worldExample(input []float64, sampleRate int) []float64 {
	w := world.New(sampleRate, defaultDioOption.FramePeriod)

	// 1. Fundamental frequency
	timeAxis, f0 := w.Dio(input, defaultDioOption)

	// 2. Refine F0 estimation result
	f0 = w.StoneMask(input, timeAxis, f0)

	// 3. Spectral envelope
	spectrogram := w.CheapTrick(input, timeAxis, f0)

	// 4. Excitation spectrum
	residual := w.Platinum(input, timeAxis, f0, spectrogram)

	// 5. Synthesis
	return w.Synthesis(f0, spectrogram, residual, len(input))
}

func main() {
	ifilename := flag.String("i", "input.wav", "Input filename")
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

	// Synthesis
	outfile := strings.Replace(*ifilename, ".wav", "_synthesized.wav", -1)
	start := time.Now()
	synthesized := worldExample(input, sampleRate)
	fmt.Println("Elapsed time in re-synthesis:", time.Now().Sub(start))

	// Write to file
	werr = wav.WriteMono(outfile, synthesized, w.SampleRate)
	if werr != nil {
		log.Fatal(werr)
	}
	fmt.Println(outfile, "is created.")

	// Synthesis from aperiodicity
	outfile = strings.Replace(*ifilename, ".wav", "_ap_synthesized.wav", -1)
	start = time.Now()
	synthesized = worldExampleAp(input, sampleRate)
	fmt.Println("Elapsed time in re-synthesis from aperiodicity:",
		time.Now().Sub(start))

	// Write to file
	werr = wav.WriteMono(outfile, synthesized, w.SampleRate)
	if werr != nil {
		log.Fatal(werr)
	}
	fmt.Println(outfile, "is created.")
}
