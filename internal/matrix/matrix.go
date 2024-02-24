package matrix

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
)

type Number interface {
	int | float64
}

type Matrix[number Number] struct {
	Data         [][]number
	Type         reflect.Type
	Shape        [2]int
	Augmented    bool
	Square       bool
	Determinants []number
	Roots        []float64
}

func NewBlankMatrix[number Number](rows, cols int, augmented bool) (*Matrix[number], error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("shape is invalid")
	}

	// create matrix with contiguous memory allocation
	matrix, memory, _ := Malloc[number](rows, cols)
	for i := 0; i < rows; i++ {
		matrix[i] = memory[(i * cols):((i + 1) * cols)]
	}

	// return pointer to new matrix
	shape := [2]int{rows, cols}
	var t [0]number
	return &Matrix[number]{
		matrix,
		reflect.TypeOf(t).Elem(),
		shape,
		augmented,
		isSquare(augmented, shape),
		makeDets[number](augmented, shape[1]),
		makeRoots(augmented, shape[1]),
	}, nil
}

func NewMatrix[number Number](matrix [][]number, augmented bool) (*Matrix[number], error) {
	if matrix == nil {
		return nil, nil
	}

	shape := [2]int{len(matrix), len(matrix[0])}
	var t [0]number
	return &Matrix[number]{
		matrix,
		reflect.TypeOf(t).Elem(),
		shape,
		augmented,
		isSquare(augmented, shape),
		makeDets[number](augmented, shape[1]),
		makeRoots(augmented, shape[1]),
	}, nil
}

func (p *Matrix[number]) FillRandom(upper int) error {
	if upper <= 0 {
		return errors.New("upper limit is invalid")
	}

	for i := 0; i < p.Shape[0]; i++ {
		for j := 0; j < p.Shape[1]; j++ {
			switch p.Type {
			case reflect.TypeOf(0):
				p.Data[i][j] = number(rand.Intn(upper))
			case reflect.TypeOf(0.0):
				p.Data[i][j] = number(rand.Float64() * float64(upper))
			}
		}
	}

	return nil
}

func isSquare(augmented bool, shape [2]int) bool {
	if augmented {
		return shape[0] == (shape[1] - 1)
	} else {
		return shape[0] == shape[1]
	}
}

func makeDets[number Number](augmented bool, cols int) []number {
	if augmented {
		return make([]number, cols)
	} else {
		return make([]number, 1)
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
func (p *Matrix[number]) Calculate() ([]number, error) {
	if !p.Square {
		return nil, errors.New("matrix is not a square")
	}

	if !p.Augmented {
		p.Determinants[0] = calcDet(p.Data)
	} else {
		p.Determinants = calcDets(p.Data)

		if p.Determinants[0] != 0 {
			for i := range p.Shape[1] - 1 {
				p.Roots[i] = float64(p.Determinants[i+1]) / float64(p.Determinants[0])
			}
		}
	}

	return p.Determinants, nil
}

func (p *Matrix[_]) GetRoots() []float64 {
	if p.Determinants[0] == 0 || p.Determinants == nil || !p.Augmented {
		return nil
	}

	for i := range p.Shape[1] - 1 {
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

func (p *Matrix[number]) String() string {
	rows := make([]string, p.Shape[0])

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
