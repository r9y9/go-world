package world

import (
	"math"
	"math/rand"
	"testing"
)

var defaultDioOption = DioOption{
	F0Floor:          71.0,
	F0Ceil:           700.0,
	FramePeriod:      5,
	ChannelsInOctave: 2.0,
	Speed:            1,
}

func createRandomSignal(length int) []float64 {
	x := make([]float64, length)
	for i := range x {
		x[i] = 32767 * rand.Float64()
	}
	return x
}

func isNaN1D(vec []float64) bool {
	for _, val := range vec {
		if math.IsNaN(val) {
			return true
		}
	}
	return false
}

func isNaN2D(mat [][]float64) bool {
	for _, vec := range mat {
		if isNaN1D(vec) {
			return true
		}
	}
	return false
}

func TestWorldNaN(t *testing.T) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	if isNaN1D(timeAxis) {
		t.Errorf("NaN detected in time axis computed by Dio")
	}
	if isNaN1D(f0) {
		t.Errorf("NaN detected in f0 computed by Dio")
	}

	f0 = w.StoneMask(x, timeAxis, f0)
	if isNaN1D(f0) {
		t.Errorf("NaN detected in f0 computed by StoneMask")
	}

	spectrogram := w.CheapTrick(x, timeAxis, f0)
	if isNaN2D(spectrogram) {
		t.Errorf("NaN detected in spectrogram computed by CheapTrick")
	}

	aperiodicity := w.D4C(x, timeAxis, f0)
	if isNaN2D(aperiodicity) {
		t.Errorf("NaN detected in aperiodicity")
	}

	y := w.Synthesis(f0, spectrogram, aperiodicity, len(x))
	if isNaN1D(y) {
		t.Errorf("NaN detected in the synthesized signal")
	}
}

func BenchmarkDio(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate) // 10 sec. data

	for i := 0; i < b.N; i++ {
		w.Dio(x, defaultDioOption)
	}
}

func BenchmarkStoneMask(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.StoneMask(x, timeAxis, f0)
	}
}

func BenchmarkCheapTrick(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.CheapTrick(x, timeAxis, f0)
	}
}

func BenchmarkD4C(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.D4C(x, timeAxis, f0)
	}
}

func BenchmarkSynthesis(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.CheapTrick(x, timeAxis, f0)
	aperiodicity := w.D4C(x, timeAxis, f0)
	for i := 0; i < b.N; i++ {
		w.Synthesis(f0, spectrogram, aperiodicity, len(x))
	}
}
