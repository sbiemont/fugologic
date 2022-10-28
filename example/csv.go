package example

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func writeCSV(filename string, title []string, values [][]float64) error {
	// Create file
	f, errCreate := os.Create(filename)
	if errCreate != nil {
		return errCreate
	}

	// Convert floats to strings
	// Convert "." to ","
	fltToStr := func(flts []float64) []string {
		result := make([]string, len(flts))
		for i, flt := range flts {
			result[i] = strings.Replace(fmt.Sprintf("%.3f", flt), ".", ",", 1)
		}
		return result
	}

	// Convert data into strings
	data := make([][]string, len(values)+1)
	data[0] = title
	for i, row := range values {
		data[i+1] = fltToStr(row)
	}

	// Write
	return csv.NewWriter(f).WriteAll(data)
}
