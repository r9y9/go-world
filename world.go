// Package world provide a Go interface to WORLD - a high-quality speech
// analysis, modification and synthesis system.
package world

// World represents WORLD - high-quality speech analysis, modification
// and synthesis system.
type World struct {
	Fs          int
	FramePeriod float64 // [ms]
}

// New returns a new world instanece with sample rate (fs) and framePeriod [ms].
func New(fs int, framePeriod float64) *World {
	w := &World{
		Fs:          fs,
		FramePeriod: framePeriod,
	}
	return w
}

// NewDioOption returns a default DioOpton.
func (w *World) NewDioOption() DioOption {
	d := NewDioOption()
	d.FramePeriod = w.FramePeriod
	return d
}

func (w *World) Dio(x []float64, option DioOption) ([]float64, []float64) {
	return Dio(x, w.Fs, option)
}

func (w *World) StoneMask(x []float64, timeAxis, f0 []float64) []float64 {
	return StoneMask(x, w.Fs, timeAxis, f0)
}

// Star is deprecated. Use CheapTrick instead.
func (w *World) Star(x []float64, timeAxis, f0 []float64) [][]float64 {
	return CheapTrick(x, w.Fs, timeAxis, f0)
}

func (w *World) CheapTrick(x []float64, timeAxis, f0 []float64) [][]float64 {
	return CheapTrick(x, w.Fs, timeAxis, f0)
}

func (w *World) D4C(x, timeAxis, f0 []float64) [][]float64 {
	return D4C(x, w.Fs, timeAxis, f0)
}

func (w *World) Synthesis(f0 []float64,
	spectrogram, aperiodicity [][]float64, length int) []float64 {
	return Synthesis(f0, spectrogram, aperiodicity, w.FramePeriod,
		w.Fs, length)
}
