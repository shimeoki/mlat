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
	m := &Matrix{}

	err := m.WriteBlank(rows, cols)
	if err != nil {
		return nil, err
	}

	m.Augmented = augmented
	m.Square = isSquare(m.Augmented, m.Rows, m.Cols)
	m.Determinants = makeDets(m.Augmented, m.Cols)
	m.Roots = makeRoots(m.Augmented, m.Cols)

	return m, nil
}

func NewMatrix(data [][]float64, augmented bool) (*Matrix, error) {
	m := &Matrix{}

	err := m.Write(data)
	if err != nil {
		return nil, err
	}

	m.Augmented = augmented
	m.Square = isSquare(m.Augmented, m.Rows, m.Cols)
	m.Determinants = makeDets(m.Augmented, m.Cols)
	m.Roots = makeRoots(m.Augmented, m.Cols)

	return m, nil
}

func (m *Matrix) FillRandom(upper int) error {
	if upper <= 0 {
		return errors.New("upper limit is invalid")
	}

	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			m.Data[i][j] = rand.Float64() * float64(upper)
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
func (m *Matrix) Calculate() ([]float64, error) {
	if !m.Square {
		return nil, errors.New("matrix is not a square")
	}

	if !m.Augmented {
		m.Determinants[0] = calcDet(m.Data)
	} else {
		m.Determinants = calcDets(m.Data)

		if m.Determinants[0] != 0 {
			for i := range m.Cols - 1 {
				m.Roots[i] = float64(m.Determinants[i+1]) / float64(m.Determinants[0])
			}
		}
	}

	return m.Determinants, nil
}

func (m *Matrix) GetRoots() []float64 {
	if m.Determinants == nil {
		return nil
	}

	if m.Determinants[0] == 0 || !m.Augmented {
		return nil
	}

	for i := range m.Cols - 1 {
		m.Roots[i] = float64(m.Determinants[i+1]) / float64(m.Determinants[0])
	}

	return m.Roots
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

func (m *Matrix) String() string {
	rows := make([]string, m.Rows)

	for i, row := range m.Data {
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

func (m *Matrix) GetAdjugate() (newMatrix [][]float64) {
	if !m.Square {
		return nil
	}

	adjugate, memory, _ := Malloc[float64](m.Rows, m.Cols)
	for i := range m.Rows {
		adjugate[i] = memory[(i * m.Cols) : (i+1)*m.Cols]
		for j := range m.Cols {
			// sign := m.Data[i][j]
			sign := 1.0
			if (i+j)%2 != 0 {
				sign = -sign
			}
			matrix := deleteRowAndCol(m.Data, i, j)
			adjugate[i][j] = calcDet(matrix) * sign
		}
	}

	return adjugate
}

func (m *Matrix) GetTranspose() (newMatrix [][]float64) {
	transpose, memory, _ := Malloc[float64](m.Cols, m.Rows)

	for i := range m.Cols {
		for _ = range m.Rows {
			transpose[i] = memory[(i * m.Rows) : (i+1)*m.Rows]
		}
	}

	for i := range m.Rows {
		for j := range m.Cols {
			transpose[j][i] = m.Data[i][j]
		}
	}

	return transpose
}

func (m *Matrix) GetInverse() (newMatrix [][]float64) {
	if !m.Square {
		return nil
	}

	if m.Determinants[0] == 0.0 {
		m.Calculate()
	}

	if m.Determinants[0] == 0.0 {
		return nil
	}

	adjugate, _ := NewMatrix(m.GetAdjugate(), m.Augmented)
	inverse := adjugate.GetTranspose()

	for i := range len(inverse) {
		for j := range len(inverse[0]) {
			inverse[i][j] /= m.Determinants[0]
		}
	}

	return inverse
}

func (m *Matrix) Multiply(matrix *Matrix) (newMatrix [][]float64) {
	if m.Cols != matrix.Rows {
		return nil
	}

	rows, cols := m.Rows, matrix.Cols
	newMatrix, memory, _ := Malloc[float64](rows, cols)
	for i := range rows {
		newMatrix[i] = memory[i*cols : (i+1)*cols]
		for j := range cols {
			cellValue := 0.0
			for k := range m.Cols {
				cellValue += m.Data[i][k] * matrix.Data[k][j]
			}
			newMatrix[i][j] = cellValue
		}
	}

	return newMatrix
}

func (m *Matrix) AddRow(index int) {
	if index < 0 || index >= m.Rows {
		return
	}

	m.Rows++
	matrix, memory, _ := Malloc[float64](m.Rows, m.Cols)
	for i := range matrix {
		matrix[i] = memory[(i * m.Cols):((i + 1) * m.Cols)]
		if i < index {
			copy(matrix[i], m.Data[i])
		} else if i > index {
			copy(matrix[i], m.Data[i-1])
		}
	}

	m.Data = matrix
}

func (m *Matrix) AddCol(index int) {
	if index < 0 || index >= m.Cols {
		return
	}

	m.Cols++
	matrix, memory, _ := Malloc[float64](m.Rows, m.Cols)
	for i := range matrix {
		matrix[i] = memory[(i * m.Cols) : (i+1)*m.Cols]
		copy(
			matrix[i],
			slices.Concat(
				m.Data[i][:index],
				make([]float64, 1),
				m.Data[i][index:]),
		)
	}

	m.Data = matrix
}

func (m *Matrix) ExtendRows(rows int) {
	if rows <= 0 {
		return
	}

	m.Rows += rows
	matrix, memory, _ := Malloc[float64](rows, m.Cols)
	for i := range matrix {
		matrix[i] = memory[i*m.Cols : (i+1)*m.Cols]
	}

	m.Data = append(m.Data, matrix...)
}

func (m *Matrix) ExtendCols(cols int) {
	if cols <= 0 {
		return
	}

	m.Cols += cols
	matrix, memory, _ := Malloc[float64](m.Rows, m.Cols)
	for i := range matrix {
		matrix[i] = memory[i*m.Cols : (i+1)*m.Cols]
		copy(matrix[i], m.Data[i])
	}

	m.Data = matrix
}

func (m *Matrix) Extend(rows, cols int) {
	m.ExtendRows(rows)
	m.ExtendCols(cols)
}

func (m *Matrix) ResizeRows(rows int) {
	if rows <= 0 || rows == m.Rows {
		return
	}

	if rows > m.Rows {
		m.ExtendRows(rows - m.Rows)
	} else {
		m.Rows = rows
		m.Data = m.Data[:rows]
	}
}

func (m *Matrix) ResizeCols(cols int) {
	if cols <= 0 || cols == m.Cols {
		return
	}

	if cols > m.Cols {
		m.ExtendCols(cols - m.Cols)
		return
	}

	m.Cols = cols
	for i := range m.Data {
		m.Data[i] = m.Data[i][:cols]
	}
}

func (m *Matrix) Resize(rows, cols int) {
	m.ResizeRows(rows)
	m.ResizeCols(cols)
}

func (m *Matrix) Write(data [][]float64) error {
	if data == nil || data[0] == nil {
		return errors.New("error: data is nil")
	}

	rows, cols := len(data), len(data[0])
	if rows == 0 || cols == 0 {
		return errors.New("error: data is empty")
	}

	matrix, memory, _ := Malloc[float64](rows, cols)
	for i := 0; i < rows; i++ {
		if len(data[i]) != cols {
			return errors.New("error: data is not rectangle-shaped")
		}

		matrix[i] = memory[(i * cols):((i + 1) * cols)]
		copy(matrix[i], data[i])
	}

	m.Data = matrix
	m.Rows = rows
	m.Cols = cols

	return nil
}

func (m *Matrix) WriteBlank(rows, cols int) error {
	if rows <= 0 || cols <= 0 {
		return errors.New("error: rows or cols equal or less than zero")
	}

	matrix, memory, _ := Malloc[float64](rows, cols)
	for i := 0; i < rows; i++ {
		matrix[i] = memory[(i * cols):((i + 1) * cols)]
	}

	m.Data = matrix
	m.Rows = rows
	m.Cols = cols

	return nil
}

func (p *Matrix) Reset() {
	p.WriteBlank(p.Rows, p.Cols)
}
