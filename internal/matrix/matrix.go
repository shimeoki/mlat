package matrix

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

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

func (m *Matrix) FillRandom(lower, upper float64) error {
	return FillRandom(m.Data, lower, upper)
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
	rows, cols := GetRowsCols(m)
	if rows == 0 {
		return false
	} else if rows != cols {
		return false
	} else {
		return true
	}
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

// Allocates memory for matrix in single allocation.
// It is more efficient in terms of access speed.
func Malloc(rows, cols int) ([][]float64, error) {
	if rows <= 0 || cols <= 0 {
		return nil, NewError("malloc: rows or cols equal or less than zero")
	}

	matrix := make([][]float64, rows)
	memory := make([]float64, rows*cols)

	for i := range matrix {
		matrix[i] = memory[i*cols : (i+1)*cols]
	}

	return matrix, nil
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

func calcDets(matrix [][]float64) []float64 {
	cols := len(matrix[0])

	chans := make([]chan float64, cols)
	dets := make([]float64, cols)

	for i := range cols {
		chans[i] = make(chan float64)

		var newMatrix [][]float64
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

func calcDet(matrix [][]float64) float64 {
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
		answer := .0

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

func (m *Matrix) DeleteRow(index int) error {
	if index < 0 || index >= m.Rows {
		return NewError("delete row: invalid row")
	}

	if m.Rows == 1 {
		m.Data = make([][]float64, 0)
		return nil
	}

	n, _ := Malloc(m.Rows-1, m.Cols)
	for i, row := range n {
		if i < index {
			copy(row, m.Data[i])
		} else {
			copy(row, m.Data[i+1])
		}
	}

	m.Data = n
	return nil
}

// Time complexity is O(n), where n is the number of rows.
//
// Does not modify original matrix.
func DeleteCol(m [][]float64, col int) ([][]float64, error) {
	rows, cols := GetRowsCols(m)
	if cols == 0 {
		return nil, NewError("delete col: invalid matrix")
	}

	if col < 0 || col >= cols {
		return nil, NewError("delete col: invalid col")
	}

	if cols == 1 {
		return make([][]float64, 0), nil
	}

	n, _ := Malloc(rows, cols-1)
	for i := range n {
		copy(n[i], slices.Concat(m[i][:col], m[i][col+1:]))
	}

	return n, nil
}

// Time complexity is O(n), where n is the number of rows.
//
// Does not modify original matrix.
func DeleteRowAndCol(m [][]float64, row, col int) ([][]float64, error) {
	rows, cols := GetRowsCols(m)
	if rows == 0 {
		return nil, NewError("delete row and col: invalid matrix")
	}

	if row < 0 || row >= rows || col < 0 || col >= cols {
		return nil, NewError("delete row and col: invalid row or col")
	}

	if rows == 1 && cols == 1 {
		return make([][]float64, 0), nil
	}

	n, _ := Malloc(rows-1, cols-1)
	for i := range n {
		var index int
		if i < row {
			index = i
		} else {
			index = i + 1
		}

		copy(n[i], slices.Concat(m[index][:col], m[index][col+1:]))
	}

	return n, nil
}

func (m *Matrix) String() string {
	rows := make([]string, m.Rows)

	for i, row := range m.Data {
		rows[i] = fmt.Sprint(row)
	}

	return strings.Join(rows, "\n")
}

// Time complexity is O(n), where n is the number of rows.
//
// Does not modify original matrix.
func ReplaceRow(m [][]float64, row []float64, index int) ([][]float64, error) {
	rows, cols := GetRowsCols(m)
	if rows == 0 {
		return nil, NewError("replace row: invalid matrix")
	}

	if index < 0 || index >= rows {
		return nil, NewError("replace row: invalid index")
	}

	n, _ := Malloc(rows, cols)
	copy(n[index], row)

	return n, nil
}

// Time complexity is O(n), where n is the number of rows.
//
// Does not modify original matrix.
func ReplaceCol(m [][]float64, col []float64, index int) ([][]float64, error) {
	rows, cols := GetRowsCols(m)
	if rows == 0 {
		return nil, NewError("replace col: invalid matrix")
	}

	if index < 0 || index >= rows {
		return nil, NewError("replace col: invalid index")
	}

	n, _ := Malloc(rows, cols)
	for i := range n {
		row := slices.Clone(n[i])
		row[index] = col[i]
		copy(n[i], row)
	}

	return n, nil
}

func replaceColInAugmented(matrix [][]float64, index int) (newMatrix [][]float64) {
	rows, cols := len(matrix), len(matrix[0])-1

	newMatrix, memory, _ := Malloc(rows, cols)
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

	adjugate, memory, _ := Malloc(m.Rows, m.Cols)
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
	transpose, memory, _ := Malloc(m.Cols, m.Rows)

	for i := range m.Cols {
		for range m.Rows {
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
	newMatrix, memory, _ := Malloc(rows, cols)
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
	matrix, memory, _ := Malloc(m.Rows, m.Cols)
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
	matrix, memory, _ := Malloc(m.Rows, m.Cols)
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
	matrix, memory, _ := Malloc(rows, m.Cols)
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
	matrix, memory, _ := Malloc(m.Rows, m.Cols)
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

	matrix, memory, _ := Malloc(rows, cols)
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
	n, err := Malloc(rows, cols)
	if err != nil {
		return err
	}

	m.Data = n
	m.Rows = rows
	m.Cols = cols

	return nil
}

func (m *Matrix) Reset() {
	m.WriteBlank(m.Rows, m.Cols)
}
