package matrix

import (
	"errors"
	"fmt"
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

func (m *Matrix) FillRandom(lower, upper float64) error {
	return FillRandom(m.Data, lower, upper)
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

func (m *Matrix) Det() float64 {
	if !m.Square {
		return .0
	}

	switch m.Rows {
	case 1:
		return m.Data[0][0]
	case 2:
		return m.Data[0][0]*m.Data[1][1] -
			m.Data[0][1]*m.Data[1][0]
	case 3:
		return m.Data[0][0]*m.Data[1][1]*m.Data[2][2] +
			m.Data[0][1]*m.Data[1][2]*m.Data[2][0] +
			m.Data[0][2]*m.Data[1][0]*m.Data[2][1] -
			m.Data[0][2]*m.Data[1][1]*m.Data[2][0] -
			m.Data[0][1]*m.Data[1][0]*m.Data[2][2] -
			m.Data[0][0]*m.Data[1][2]*m.Data[2][1]
	default:
		answer := .0

		for i, row := range m.Data {
			value := row[0]
			if i%2 != 0 {
				value = -value
			}

			n, _ := m.NewDeleteRowAndCol(i, 0)
			answer += value * n.Det()
		}

		return answer
	}
}

func (m *Matrix) NewDeleteRow(index int) (*Matrix, error) {
	if index < 0 || index >= m.Rows {
		return nil, NewError("delete row: invalid row")
	}

	if m.Rows == 1 {
		return NewBlankMatrix(0, 0, m.Augmented)
	}

	n, _ := Malloc(m.Rows-1, m.Cols)
	for i, row := range n {
		if i < index {
			copy(row, m.Data[i])
		} else {
			copy(row, m.Data[i+1])
		}
	}

	return NewMatrix(n, m.Augmented)
}

func (m *Matrix) DeleteRow(index int) error {
	n, err := m.NewDeleteRow(index)
	if err != nil {
		return err
	}

	m.Data = n.Data
	return nil
}

func (m *Matrix) NewDeleteCol(index int) (*Matrix, error) {
	if index < 0 || index >= m.Cols {
		return nil, NewError("delete col: invalid col")
	}

	if m.Cols == 1 {
		return NewBlankMatrix(0, 0, m.Augmented)
	}

	n, _ := Malloc(m.Rows, m.Cols-1)
	for i, row := range n {
		copy(row, slices.Concat(m.Data[i][:index], m.Data[i][index+1:]))
	}

	return NewMatrix(n, m.Augmented)
}

func (m *Matrix) DeleteCol(index int) error {
	n, err := m.NewDeleteCol(index)
	if err != nil {
		return err
	}

	m.Data = n.Data
	return nil
}

func (m *Matrix) NewDeleteRowAndCol(row, col int) (*Matrix, error) {
	if row < 0 || row >= m.Rows || col < 0 || col >= m.Cols {
		return nil, NewError("delete row and col: invalid row or col")
	}

	if m.Rows == 1 && m.Cols == 1 {
		return NewBlankMatrix(0, 0, m.Augmented)
	}

	n, _ := Malloc(m.Rows-1, m.Cols-1)
	for i := range n {
		index := i
		if i >= row {
			index++
		}

		copy(n[i], slices.Concat(m.Data[index][:col], m.Data[index][col+1:]))
	}

	return NewMatrix(n, m.Augmented)
}

func (m *Matrix) DeleteRowAndCol(row, col int) error {
	n, err := m.NewDeleteRowAndCol(row, col)
	if err != nil {
		return err
	}

	m.Data = n.Data
	return nil
}

func (m *Matrix) String() string {
	rows := make([]string, m.Rows)

	for i, row := range m.Data {
		rows[i] = fmt.Sprint(row)
	}

	return strings.Join(rows, "\n")
}

func (m *Matrix) ReplaceRow(row []float64, index int) error {
	if index < 0 || index >= m.Rows {
		return NewError("replace row: invalid index")
	}

	if len(row) != m.Cols {
		return NewError("replace col: invalid row")
	}

	copy(m.Data[index], row)
	return nil
}

func (m *Matrix) ReplaceCol(col []float64, index int) error {
	if index < 0 || index >= m.Cols {
		return NewError("replace col: invalid index")
	}

	if len(col) != m.Rows {
		return NewError("replace col: invalid col")
	}

	for i, value := range col {
		m.Data[i][index] = value
	}

	return nil
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

func (m *Matrix) NewCofactor() (*Matrix, error) {
	if !m.Square {
		return nil, NewError("get adjugate: matrix is not a square")
	}

	cofactor, _ := Malloc(m.Rows, m.Cols)
	for i := range m.Rows {
		for j := range m.Cols {
			n, _ := m.NewDeleteRowAndCol(i, j)
			value := n.Det()
			if (i+j)%2 != 0 {
				value = -value
			}

			cofactor[i][j] = value
		}
	}

	return NewMatrix(cofactor, false)
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
