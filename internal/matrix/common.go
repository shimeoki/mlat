package matrix

import (
	"math/rand"
)

func FillRandom(m [][]float64, lower, upper float64) error {
	if lower >= upper {
		return NewError("fill random: lower >= upper")
	}

	for i := range m {
		for j := range m[i] {
			m[i][j] = lower + upper*rand.Float64()
		}
	}

	return nil
}

// Time complexity is O(n), where n is the number of rows.
//
// Returns (0, 0) if m is not valid.
func GetRowsCols(m [][]float64) (rows, cols int) {
	rows, cols = 0, 0

	for i := range m {
		rows++

		colsInRow := len(m[i])

		if cols == 0 {
			if colsInRow == 0 {
				return 0, 0
			}

			cols = colsInRow
			continue
		}

		if cols != colsInRow {
			return 0, 0
		}
	}

	return
}

func IsRectangle(m [][]float64) bool {
	if rows, _ := GetRowsCols(m); rows == 0 {
		return false
	} else {
		return true
	}
}

func IsSquare(m [][]float64) bool {
	if rows, cols := GetRowsCols(m); rows == 0 {
		return false
	} else if rows != cols {
		return false
	} else {
		return true
	}
}

// Allocates memory for matrix in single allocation.
// It is more efficient in terms of access speed.
//
// Returns empty matrix if rows or cols are equal to zero.
func Malloc(rows, cols int) ([][]float64, error) {
	if rows < 0 || cols < 0 {
		return nil, NewError("malloc: rows or cols are less than zero")
	}

	if rows == 0 || cols == 0 {
		return make([][]float64, 0), nil
	}

	matrix := make([][]float64, rows)
	memory := make([]float64, rows*cols)

	for i := range matrix {
		matrix[i] = memory[i*cols : (i+1)*cols]
	}

	return matrix, nil
}
