package world

import "C"

func Make2DCArrayAlternative(matrix [][]float64) []*C.double {
	alternative := make([]*C.double, len(matrix))
	for i := range alternative {
		// DO NOT free because the source slice is managed by Go
		alternative[i] = (*C.double)(&matrix[i][0])
	}
	return alternative
}
