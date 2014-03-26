// Package world provide a Go interface to WORLD - a high-quality speech analysis, modification and synthesis system.
package world

// Functions in this file ara simple ports to original cpp functions.
// Go-like interfaces will be found in go_interface.go

// #cgo pkg-config: world
// #include <world/dio.h>
// #include <world/stonemask.h>
// #include <world/platinum.h>
// #include <world/synthesis.h>
// #include <world/star.h>
// #include <world/aperiodicity.h>
// #include <world/synthesisfromaperiodicity.h>
// #include <world/tandem_ap.h>
// #include <world/synthesis_ap.h>
// #include <stdlib.h>
import "C"

const (
	byteSizeOfFloat64 = 8
)

type DioOption struct {
	F0Floor          float64
	F0Ceil           float64
	ChannelsInOctave float64
	FramePeriod      float64 // [ms]
	Speed            int     // (1,2, ..., 12)
}

func NewDioOption() DioOption {
	opt := DioOption{
		F0Floor:          80.0,
		F0Ceil:           800.0,
		FramePeriod:      5,
		ChannelsInOctave: 4.0,
		Speed:            2,
	}
	return opt
}

func Dio(x []float64, fs int, option DioOption) ([]float64, []float64) {
	numFrames := GetSamplesForDIO(fs, len(x), option.FramePeriod)
	timeAxis := make([]float64, numFrames)
	f0 := make([]float64, numFrames)

	// Create C interface of DioOption
	var opt C.DioOption
	opt.frame_period = C.double(option.FramePeriod)
	opt.f0_floor = C.double(option.F0Floor)
	opt.f0_ceil = C.double(option.F0Ceil)
	opt.channels_in_octave = C.double(option.ChannelsInOctave)
	opt.speed = C.int(option.Speed)

	// Perform DIO
	C.Dio((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		opt,
		(*C.double)(&timeAxis[0]),
		(*C.double)(&f0[0]))

	return timeAxis, f0
}

func GetSamplesForDIO(fs, x_length int, framePeriod float64) int {
	return int(C.GetSamplesForDIO(C.int(fs), C.int(x_length), C.double(framePeriod)))
}

func StoneMask(x []float64, fs int, timeAxis, f0 []float64) []float64 {
	refinedF0 := make([]float64, len(f0))

	// Perform StoneMask
	C.StoneMask((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		(*C.double)(&timeAxis[0]),
		(*C.double)(&f0[0]),
		C.int(len(f0)),
		(*C.double)(&refinedF0[0]))

	return refinedF0
}

func Star(x []float64, fs int, timeAxis, f0 []float64) [][]float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))
	numFreqBins := C.size_t(FFTSize/2 + 1)

	spectrogram := make([][]float64, len(f0))
	for i := range spectrogram {
		spectrogram[i] = make([]float64, numFreqBins)
	}

	spectrogramUsedInC := Make2DCArrayAlternative(spectrogram)

	// Perform star
	C.Star((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		(*C.double)(&timeAxis[0]),
		(*C.double)(&f0[0]),
		C.int(len(f0)),
		(**C.double)(&spectrogramUsedInC[0]))

	return spectrogram
}

func GetFFTSizeForStar(fs int) int {
	return int(C.GetFFTSizeForStar(C.int(fs)))
}

func Platinum(x []float64, fs int, timeAxis, f0 []float64, spectrogram [][]float64) [][]float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))
	numFreqBins := C.size_t(FFTSize + 1)

	residualSpectrogram := make([][]float64, len(f0))
	for i := range residualSpectrogram {
		residualSpectrogram[i] = make([]float64, numFreqBins)
	}

	spectrogramUsedInC := Make2DCArrayAlternative(spectrogram)
	residualSpectrogramUsedInC := Make2DCArrayAlternative(residualSpectrogram)

	C.Platinum((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		(*C.double)(&timeAxis[0]),
		(*C.double)(&f0[0]),
		C.int(len(f0)),
		(**C.double)(&spectrogramUsedInC[0]),
		C.int(FFTSize),
		(**C.double)(&residualSpectrogramUsedInC[0]))

	return residualSpectrogram
}

func Synthesis(f0 []float64, spectrogram, residualSpectrogram [][]float64, framePeriod float64, fs, length int) []float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))

	spectrogramUsedInC := Make2DCArrayAlternative(spectrogram)
	residualSpectrogramUsedInC := Make2DCArrayAlternative(residualSpectrogram)

	synthesized := make([]float64, length)
	C.Synthesis((*C.double)(&f0[0]),
		C.int(len(f0)),
		(**C.double)(&spectrogramUsedInC[0]),
		(**C.double)(&residualSpectrogramUsedInC[0]),
		C.int(FFTSize),
		C.double(framePeriod),
		C.int(fs),
		C.int(len(synthesized)),
		(*C.double)(&synthesized[0]))

	return synthesized
}

func GetNumberOfBands(fs int) int {
	return int(C.GetNumberOfBands(C.int(fs)))
}

func AperiodicityRatio(x []float64, fs int, f0, timeAxis []float64) [][]float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))
	numBins := C.size_t(FFTSize/2 + 1)

	aperiodicity := make([][]float64, len(f0))
	for i := range aperiodicity {
		aperiodicity[i] = make([]float64, numBins)
	}
	aperiodicityUsedInC := Make2DCArrayAlternative(aperiodicity)

	// Peform aperiodicity analysis
	C.AperiodicityRatio((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		(*C.double)(&f0[0]),
		C.int(len(f0)),
		(*C.double)(&timeAxis[0]),
		C.int(FFTSize),
		(**C.double)(&aperiodicityUsedInC[0]))

	return aperiodicity
}

// deprecated (will be removed)
func AperiodicityRatioOld(x []float64, fs int, f0 []float64, framePeriod float64) ([][]float64, float64) {
	numBands := GetNumberOfBands(fs)

	aperiodicity := make([][]float64, len(f0))
	for i := range aperiodicity {
		aperiodicity[i] = make([]float64, numBands)
	}
	aperiodicityUsedInC := Make2DCArrayAlternative(aperiodicity)

	// Peform aperiodicity analysis
	targetF0 := C.AperiodicityRatioOld((*C.double)(&x[0]),
		C.int(len(x)),
		C.int(fs),
		(*C.double)(&f0[0]),
		C.int(len(f0)),
		C.double(framePeriod),
		(**C.double)(&aperiodicityUsedInC[0]))

	return aperiodicity, float64(targetF0)
}

func SynthesisFromAperiodicity(f0 []float64, spectrogram, aperiodicity [][]float64, framePeriod float64, fs, length int) []float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))

	spectrogramUsedInC := Make2DCArrayAlternative(spectrogram)
	aperiodicityUsedInC := Make2DCArrayAlternative(aperiodicity)

	synthesized := make([]float64, length)
	C.SynthesisFromAperiodicity((*C.double)(&f0[0]),
		C.int(len(f0)),
		(**C.double)(&spectrogramUsedInC[0]),
		(**C.double)(&aperiodicityUsedInC[0]),
		C.int(FFTSize),
		C.double(framePeriod),
		C.int(fs),
		C.int(len(synthesized)),
		(*C.double)(&synthesized[0]))

	return synthesized
}

// deprecated (will be removed)
func SynthesisFromAperiodicityOld(f0 []float64, spectrogram, aperiodicity [][]float64, targetF0, framePeriod float64, fs, length int) []float64 {
	FFTSize := C.size_t(C.GetFFTSizeForStar(C.int(fs)))
	numBands := GetNumberOfBands(fs)

	spectrogramUsedInC := Make2DCArrayAlternative(spectrogram)
	aperiodicityUsedInC := Make2DCArrayAlternative(aperiodicity)

	synthesized := make([]float64, length)
	C.SynthesisFromAperiodicityOld((*C.double)(&f0[0]),
		C.int(len(f0)),
		(**C.double)(&spectrogramUsedInC[0]),
		C.int(FFTSize),
		(**C.double)(&aperiodicityUsedInC[0]),
		C.int(numBands),
		C.double(targetF0),
		C.double(framePeriod),
		C.int(fs),
		C.int(len(synthesized)),
		(*C.double)(&synthesized[0]))

	return synthesized
}
