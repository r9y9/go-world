package world

import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func CArrayToGoSlice(array *C.double, length C.int) []float64 {
	slice := make([]float64, int(length))
	b := C.GoBytes(unsafe.Pointer(array), C.int(8*length))
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, slice)
	if err != nil {
		panic(err)
	}
	return slice
}

func Make2DCArrayAlternative(matrix [][]float64) []*C.double {
	alternative := make([]*C.double, len(matrix))
	for i := range alternative {
		// DO NOT free because the source slice is managed by Go
		alternative[i] = (*C.double)(&matrix[i][0])
	}
	return alternative
}
