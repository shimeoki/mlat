package matrix

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

type Number interface {
	int | float64
}

type Matrix struct {
	Data         [][]float64
	Rows         int
	Cols         int
	Augmented    bool
	Square       bool
	Determinants []float64
	Roots        []float64
}

func NewBlankMatrix(rows, cols int, augmented bool) (*Matrix, error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("shape is invalid")
	}

	// create matrix with contiguous memory allocation
	matrix, memory, _ := Malloc[float64](rows, cols)
	for i := 0; i < rows; i++ {
		matrix[i] = memory[(i * cols):((i + 1) * cols)]
	}

	// return pointer to new matrix
	return &Matrix{
		matrix,
		rows,
		cols,
		augmented,
		isSquare(augmented, rows, cols),
		makeDets(augmented, cols),
		makeRoots(augmented, cols),
	}, nil
}

func NewMatrix(matrix [][]float64, augmented bool) (*Matrix, error) {
	if matrix == nil || matrix[0] == nil {
		return nil, errors.New("error: matrix is nil")
	}

	rows, cols := len(matrix), len(matrix[0])
	if rows == 0 || cols == 0 {
		return nil, errors.New("error: matrix has empty rows or cols")
	}

	cmatrix, memory, _ := Malloc[float64](rows, cols)
	for i := 0; i < rows; i++ {
		if len(matrix[i]) != cols {
			return nil, errors.New("error: matrix is invalid")
		}
		cmatrix[i] = memory[(i * cols):((i + 1) * cols)]
		copy(cmatrix[i], matrix[i])
	}

	return &Matrix{
		cmatrix,
		rows,
		cols,
		augmented,
		isSquare(augmented, rows, cols),
		makeDets(augmented, cols),
		makeRoots(augmented, cols),
	}, nil
}

func (p *Matrix) FillRandom(upper int) error {
	if upper <= 0 {
		return errors.New("upper limit is invalid")
	}

	for i := 0; i < p.Rows; i++ {
		for j := 0; j < p.Cols; j++ {
			p.Data[i][j] = rand.Float64() * float64(upper)
		}
	}

	return nil
}

func isSquare(augmented bool, rows, cols int) bool {
	if augmented {
		return rows == (cols - 1)
	} else {
		return rows == cols
	}
}

func makeDets(augmented bool, cols int) []float64 {
	if augmented {
		return make([]float64, cols)
	} else {
		return make([]float64, 1)
	}
}

func makeRoots(augmented bool, cols int) []float64 {
	if augmented {
		return make([]float64, cols-1)
	} else {
		return nil
	}
}

// allocates memory closely for contingious memory allocation afterwards
//
//	for i := range matrix {
//		matrix[i] = memory[(i * cols):((i + 1) * cols)]
//	}
func Malloc[number Number](rows, cols int) ([][]number, []number, error) {
	if rows <= 0 || cols <= 0 {
		return nil, nil, errors.New("invalid row or column value")
	}

	return make([][]number, rows), make([]number, rows*cols), nil
}

// if matrix is not square, returns nil and error
//
// matrix is square if matrix is not augmented and rows equals cols
// or matrix is augmented and rows equals cols - 1
//
// calculates determinants and roots if matrix is augmented. returns determinants, not roots
//
// if matrix is not augmented, returns [1]number.
// otherwise returns [cols]number
func (p *Matrix) Calculate() ([]float64, error) {
	if !p.Square {
		return nil, errors.New("matrix is not a square")
	}

	if !p.Augmented {
		p.Determinants[0] = calcDet(p.Data)
	} else {
		p.Determinants = calcDets(p.Data)

		if p.Determinants[0] != 0 {
			for i := range p.Cols - 1 {
				p.Roots[i] = float64(p.Determinants[i+1]) / float64(p.Determinants[0])
			}
		}
	}

	return p.Determinants, nil
}

func (p *Matrix) GetRoots() []float64 {
	if p.Determinants[0] == 0 || p.Determinants == nil || !p.Augmented {
		return nil
	}

	for i := range p.Cols - 1 {
		p.Roots[i] = float64(p.Determinants[i+1]) / float64(p.Determinants[0])
	}

	return p.Roots
}

func calcDets[number Number](matrix [][]number) []number {
	cols := len(matrix[0])

	chans := make([]chan number, cols)
	dets := make([]number, cols)

	for i := range cols {
		chans[i] = make(chan number)

		var newMatrix [][]number
		if i == 0 {
			newMatrix = deleteCol(matrix, cols-1)
		} else {
			newMatrix = replaceColInAugmented(matrix, i-1)
		}

		go func(i int) {
			chans[i] <- calcDet(newMatrix)
		}(i)
	}

	for i := range chans {
		dets[i] = <-chans[i]
	}

	return dets
}

func calcDet[number Number](matrix [][]number) number {
	switch dimension := len(matrix); dimension {
	case 1:
		return matrix[0][0]
	case 2:
		return matrix[0][0]*matrix[1][1] -
			matrix[0][1]*matrix[1][0]
	case 3:
		return matrix[0][0]*matrix[1][1]*matrix[2][2] +
			matrix[0][1]*matrix[1][2]*matrix[2][0] +
			matrix[0][2]*matrix[1][0]*matrix[2][1] -
			matrix[0][2]*matrix[1][1]*matrix[2][0] -
			matrix[0][1]*matrix[1][0]*matrix[2][2] -
			matrix[0][0]*matrix[1][2]*matrix[2][1]
	default:
		answer := number(0)

		for i := range dimension {
			value := matrix[i][0]
			if i%2 != 0 {
				value = -value
			}

			rewritedMatrix := deleteRowAndCol(matrix, i, 0)
			answer += value * calcDet(rewritedMatrix)
		}

		return answer
	}
}

func deleteRowAndCol[number Number](matrix [][]number, row, col int) (newMatrix [][]number) {
	rows, cols := len(matrix)-1, len(matrix[0])-1

	newMatrix, memory, _ := Malloc[number](rows, cols)
	for i := range newMatrix {
		newMatrix[i] = memory[(i * cols):((i + 1) * cols)]
		if i < row {
			copy(newMatrix[i], slices.Concat(matrix[i][:col], matrix[i][col+1:]))
		} else {
			copy(newMatrix[i], slices.Concat(matrix[i+1][:col], matrix[i+1][col+1:]))
		}
	}

	return
}

func (p *Matrix) String() string {
	rows := make([]string, p.Rows)

	for i, row := range p.Data {
		rows[i] = fmt.Sprint(row)
	}

	return strings.Join(rows, "\n")
}

func deleteCol[number Number](matrix [][]number, col int) (newMatrix [][]number) {
	rows, cols := len(matrix), len(matrix[0])-1

	newMatrix, memory, _ := Malloc[number](rows, cols)
	for i := range matrix {
		newMatrix[i] = memory[(i * cols):((i + 1) * cols)]
		copy(newMatrix[i], slices.Concat(matrix[i][:col], matrix[i][col+1:]))
	}

	return
}

func replaceColInAugmented[number Number](matrix [][]number, index int) (newMatrix [][]number) {
	rows, cols := len(matrix), len(matrix[0])-1

	newMatrix, memory, _ := Malloc[number](rows, cols)
	for i := range matrix {
		newMatrix[i] = memory[(i * cols) : (i+1)*cols]
		copy(newMatrix[i], slices.Concat(matrix[i][:index], matrix[i][cols:], matrix[i][index+1:cols]))
	}

	return
}

func (p *Matrix) GetAdjugate() (newMatrix [][]float64) {
	if !p.Square {
		return nil
	}

	adjugate, memory, _ := Malloc[float64](p.Rows, p.Cols)
	for i := range p.Rows {
		adjugate[i] = memory[(i * p.Cols) : (i+1)*p.Cols]
		for j := range p.Cols {
			cellValue := p.Data[i][j]
			if (i+j)%2 != 0 {
				cellValue = -cellValue
			}
			matrix := deleteRowAndCol(p.Data, i, j)
			adjugate[i][j] = cellValue * calcDet(matrix)
		}
	}

	return adjugate
}

func (p *Matrix) Multiply(matrix *Matrix) (newMatrix [][]float64) {
	if p.Cols != matrix.Rows {
		return nil
	}

	rows, cols := p.Rows, matrix.Cols
	newMatrix, memory, _ := Malloc[float64](rows, cols)
	for i := range rows {
		newMatrix[i] = memory[i*cols : (i+1)*cols]
		for j := range cols {
			cellValue := 0.0
			for k := range p.Cols {
				cellValue += p.Data[i][k] * matrix.Data[k][j]
			}
			newMatrix[i][j] = cellValue
		}
	}

	return newMatrix
}

func (p *Matrix) AddRow(index int) {
	if index < 0 || index >= p.Rows {
		return
	}

	p.Rows++
	matrix, memory, _ := Malloc[float64](p.Rows, p.Cols)
	for i := range matrix {
		matrix[i] = memory[(i * p.Cols):((i + 1) * p.Cols)]
		if i < index {
			copy(matrix[i], p.Data[i])
		} else if i > index {
			copy(matrix[i], p.Data[i-1])
		}
	}

	p.Data = matrix
}

func (p *Matrix) AddCol(index int) {
	if index < 0 || index >= p.Cols {
		return
	}

	p.Cols++
	matrix, memory, _ := Malloc[float64](p.Rows, p.Cols)
	for i := range matrix {
		matrix[i] = memory[(i * p.Cols) : (i+1)*p.Cols]
		copy(
			matrix[i],
			slices.Concat(
				p.Data[i][:index],
				make([]float64, 1),
				p.Data[i][index:]),
		)
	}

	p.Data = matrix
}

func (p *Matrix) ExtendRows(rows int) {
	if rows <= 0 {
		return
	}

	p.Rows += rows
	matrix, memory, _ := Malloc[float64](rows, p.Cols)
	for i := range matrix {
		matrix[i] = memory[i*p.Cols : (i+1)*p.Cols]
	}

	p.Data = append(p.Data, matrix...)
}

func (p *Matrix) ExtendCols(cols int) {
	if cols <= 0 {
		return
	}

	p.Cols += cols
	matrix, memory, _ := Malloc[float64](p.Rows, p.Cols)
	for i := range matrix {
		matrix[i] = memory[i*p.Cols : (i+1)*p.Cols]
		copy(matrix[i], p.Data[i])
	}

	p.Data = matrix
}

func (p *Matrix) Extend(rows, cols int) {
	p.ExtendRows(rows)
	p.ExtendCols(cols)
}

func (p *Matrix) ResizeRows(rows int) {
	if rows <= 0 || rows == p.Rows {
		return
	}

	if rows > p.Rows {
		p.ExtendRows(rows - p.Rows)
	} else {
		p.Rows = rows
		p.Data = p.Data[:rows]
	}
}

func (p *Matrix) ResizeCols(cols int) {
	if cols <= 0 || cols == p.Cols {
		return
	}

	if cols > p.Cols {
		p.ExtendCols(cols - p.Cols)
		return
	}

	p.Cols = cols
	for i := range p.Data {
		p.Data[i] = p.Data[i][:cols]
	}
}

func (p *Matrix) Resize(rows, cols int) {
	p.ResizeRows(rows)
	p.ResizeCols(cols)
}

func (p *Matrix) Write(data [][]float64) error {
	if data == nil || data[0] == nil {
		return errors.New("error: data is nil")
	}

	rows, cols := len(data), len(data[0])
	matrix, memory, _ := Malloc[float64](rows, cols)
	for i := 0; i < rows; i++ {
		if len(data[i]) != cols {
			return errors.New("error: data is not rectangle-shaped")
		}

		matrix[i] = memory[(i * cols):((i + 1) * cols)]
		copy(matrix[i], data[i])
	}

	p.Data = matrix
	p.Rows = rows
	p.Cols = cols

	return nil
}
