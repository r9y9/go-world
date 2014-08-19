package world

import (
	"math"
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

	spectrogram := w.Star(x, timeAxis, f0)
	if isNaN2D(spectrogram) {
		t.Errorf("NaN detected in spectrogram computed by Star")
	}

	spectrogram = w.CheapTrick(x, timeAxis, f0)
	if isNaN2D(spectrogram) {
		t.Errorf("NaN detected in spectrogram computed by CheapTrick")
	}

	residual := w.Platinum(x, timeAxis, f0, spectrogram)
	if isNaN2D(residual) {
		t.Errorf("NaN detected in residual computed by Platinum")
	}

	y := w.Synthesis(f0, spectrogram, residual, len(x))

	if isNaN1D(y) {
		t.Errorf("NaN detected in the synthesized signal")
	}
}

func TestWorldAperiodicityNaN(t *testing.T) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.CheapTrick(x, timeAxis, f0)
	aperiodicity := w.AperiodicityRatio(x, f0, timeAxis)
	if isNaN2D(aperiodicity) {
		t.Errorf("NaN detected in aperiodicity")
	}

	y := w.SynthesisFromAperiodicity(f0, spectrogram, aperiodicity, len(x))

	if isNaN1D(y) {
		t.Errorf("NaN detected in the synthesized signal by aperiodicity")
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

func BenchmarkStar(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.Star(x, timeAxis, f0)
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

func BenchmarkPlatinum(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.Star(x, timeAxis, f0)
	for i := 0; i < b.N; i++ {
		w.Platinum(x, timeAxis, f0, spectrogram)
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
	spectrogram := w.Star(x, timeAxis, f0)
	residual := w.Platinum(x, timeAxis, f0, spectrogram)
	for i := 0; i < b.N; i++ {
		w.Synthesis(f0, spectrogram, residual, len(x))
	}
}

func BenchmarkAperiodicityRatio(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	for i := 0; i < b.N; i++ {
		w.AperiodicityRatio(x, f0, timeAxis)
	}
}

func BenchmarkSynthesisFromAperiodicity(b *testing.B) {
	sampleRate := 44100
	w := &World{
		Fs:          sampleRate,
		FramePeriod: defaultDioOption.FramePeriod,
	}
	x := createRandomSignal(10 * sampleRate)

	timeAxis, f0 := w.Dio(x, defaultDioOption)
	spectrogram := w.Star(x, timeAxis, f0)
	aperiodicity := w.AperiodicityRatio(x, f0, timeAxis)
	for i := 0; i < b.N; i++ {
		w.SynthesisFromAperiodicity(f0, spectrogram, aperiodicity, len(x))
	}
}
