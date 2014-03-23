package world

import (
	"math/rand"
	"testing"
)

var defaultDioOption = DioOption{
	F0Floor:          80.0,
	F0Ceil:           800.0,
	FramePeriod:      5,
	ChannelsInOctave: 4.0,
	Speed:            2,
}

func createRandomSignal(length int) []float64 {
	x := make([]float64, length)
	for i := range x {
		x[i] = 32767 * rand.Float64()
	}
	return x
}

func BenchmarkDio(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate) // 10 sec. data

	for i := 0; i < b.N; i++ {
		w.Dio(x, defaultDioOption)
	}
}

func BenchmarkStoneMask(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.StoneMask(x, timeAxis, f0)
	}
}

func BenchmarkStar(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.Star(x, timeAxis, f0)
	}
}

func BenchmarkPlatinum(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.Star(x, timeAxis, f0)
	for i := 0; i < b.N; i++ {
		w.Platinum(x, timeAxis, f0, spectrogram)
	}
}

func BenchmarkSynthesis(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.Star(x, timeAxis, f0)
	residual := w.Platinum(x, timeAxis, f0, spectrogram)
	for i := 0; i < b.N; i++ {
		w.Synthesis(f0, spectrogram, residual, len(x))
	}
}

func BenchmarkAperiodicityRatio(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.AperiodicityRatio(x, f0, timeAxis)
	}
}

func BenchmarkSynthesisFromAperiodicity(b *testing.B) {
	sampleRate := 44100
	w := New(sampleRate, defaultDioOption.FramePeriod)
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.Star(x, timeAxis, f0)
	aperiodicity := w.AperiodicityRatio(x, f0, timeAxis)
	for i := 0; i < b.N; i++ {
		w.SynthesisFromAperiodicity(f0, spectrogram, aperiodicity, len(x))
	}
}
