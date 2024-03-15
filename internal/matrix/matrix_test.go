package matrix

import (
	"os"
	"testing"
)

var matrix1 [][]float64
var matrix2 [][]float64
var matrix3 [][]float64
var matrix4 [][]float64
var matrix5 [][]float64
var matrix6 [][]float64
var matrix7 [][]float64
var matrix8 [][]float64
var matrix9 [][]float64

func TestMain(m *testing.M) {
	matrix1 = [][]float64{
		{},
	}
	matrix2 = [][]float64{
		{0},
	}
	matrix3 = [][]float64{
		{0, 0, 0},
	}
	matrix4 = [][]float64{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	matrix5 = [][]float64{
		{1, 2, 3},
	}
	matrix6 = [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	matrix7 = [][]float64{
		{ 0 },
		{ 0 },
		{ 0 },
	}
	matrix8 = [][]float64{
		{ 1 },
		{ 2 },
		{ 3 },
	}
	matrix9 = [][]float64{
		{ 1 },
	}

	tests := m.Run()

	os.Exit(tests)
}

func TestAddRow(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix6, -1, MatrixToString(matrix6, " ") },
		{ "arg 0 with blank", matrix1, 0, "\n" },
		{ "arg 0 with 1 0-elem", matrix2, 0, "0\n0\n" },
		{ "arg 0 with 0-row", matrix3, 0, "0 0 0\n0 0 0\n" },
		{ "arg 0 with three 0-rows", matrix4, 0, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 0 with 1 elem", matrix9, 0, "0\n1\n" },
		{ "arg 0 with row", matrix5, 0, "0 0 0\n1 2 3\n" },
		{ "arg 0 with three rows", matrix6, 0, "0 0 0\n1 2 3\n4 5 6\n7 8 9\n" },
		{ "arg 0 with 0-column", matrix7, 0, "0\n0\n0\n0\n" },
		{ "arg 0 with column", matrix8, 0, "0\n1\n2\n3\n" },
		{ "arg 1 with three rows", matrix6, 1, "1 2 3\n0 0 0\n4 5 6\n7 8 9\n" },
		{ "arg 2 with three 0-rows", matrix4, 2, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 2 with three rows", matrix6, 2, "1 2 3\n4 5 6\n0 0 0\n7 8 9\n" },
		{ "arg 2 with 0-column", matrix7, 2, "0\n0\n0\n0\n" },
		{ "arg 2 with column", matrix8, 2, "1\n2\n0\n3\n" },
		{ "arg 3 with row", matrix5, 3, "1 2 3\n" },
		{ "arg 3 with three rows", matrix6, 3, "1 2 3\n4 5 6\n7 8 9\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.AddRow(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestAddCol(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix6, -1, MatrixToString(matrix6, " ") },
		{ "arg 0 with blank", matrix1, 0, "\n" },
		{ "arg 0 with 1 0-elem", matrix2, 0, "0 0\n" },
		{ "arg 0 with 0-row", matrix3, 0, "0 0 0 0\n" },
		{ "arg 0 with three 0-rows", matrix4, 0, "0 0 0 0\n0 0 0 0\n0 0 0 0\n"},
		{ "arg 0 with 1 elem", matrix9, 0, "0 1\n"},
		{ "arg 0 with row", matrix5, 0, "0 1 2 3\n"},
		{ "arg 0 with three rows", matrix6, 0, "0 1 2 3\n0 4 5 6\n0 7 8 9\n"},
		{ "arg 1 with three rows", matrix6, 1, "1 0 2 3\n4 0 5 6\n7 0 8 9\n"},
		{ "arg 2 with three rows", matrix6, 2, "1 2 0 3\n4 5 0 6\n7 8 0 9\n"},
		{ "arg 3 with three rows", matrix6, 3, "1 2 3\n4 5 6\n7 8 9\n"},
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.AddCol(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestExtendRows(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix2, -1, "0\n" },
		{ "arg 0", matrix2, 0, "0\n" },
		{ "arg 1 with 1 0-elem", matrix2, 1, "0\n0\n" },
		{ "arg 1 with 0-row", matrix3, 1, "0 0 0\n0 0 0\n" },
		{ "arg 1 with three 0-rows", matrix4, 1, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 1 with row", matrix5, 1, "1 2 3\n0 0 0\n" },
		{ "arg 1 with three rows", matrix6, 1, "1 2 3\n4 5 6\n7 8 9\n0 0 0\n" },
		{ "arg 2 with row", matrix5, 2, "1 2 3\n0 0 0\n0 0 0\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.ExtendRows(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestExtendCols(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix2, -1, "0\n" },
		{ "arg 0", matrix2, 0, "0\n" },
		{ "arg 1 with 1 0-elem", matrix2, 1, "0 0\n" },
		{ "arg 1 with 0-row", matrix3, 1, "0 0 0 0\n" },
		{ "arg 1 with three 0-rows", matrix4, 1, "0 0 0 0\n0 0 0 0\n0 0 0 0\n" },
		{ "arg 1 with row", matrix5, 1, "1 2 3 0\n" },
		{ "arg 1 with three rows", matrix6, 1, "1 2 3 0\n4 5 6 0\n7 8 9 0\n" },
		{ "arg 2 with 0-row", matrix3, 2, "0 0 0 0 0\n" },
		{ "arg 2 with row", matrix5, 2, "1 2 3 0 0\n" },
		{ "arg 3 with three rows", matrix6, 3, "1 2 3 0 0 0\n4 5 6 0 0 0\n7 8 9 0 0 0\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.ExtendCols(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestResizeRows(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix2, -1, "0\n" },
		{ "arg 0", matrix2, 0, "0\n" },
		{ "arg 1 with 1 0-elem", matrix2, 1, "0\n" },
		{ "arg 2 with 1 0-elem", matrix2, 2, "0\n0\n" },
		{ "arg 3 with 1 0-elem", matrix2, 3, "0\n0\n0\n" },
		{ "arg 1 with 0-row", matrix3, 1, "0 0 0\n" },
		{ "arg 2 with 0-row", matrix3, 2, "0 0 0\n0 0 0\n" },
		{ "arg 3 with 0-row", matrix3, 3, "0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 1 with three 0-rows", matrix4, 1, "0 0 0\n" },
		{ "arg 2 with three 0-rows", matrix4, 2, "0 0 0\n0 0 0\n" },
		{ "arg 3 with three 0-rows", matrix4, 3, "0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 4 with three 0-rows", matrix4, 4, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 5 with three 0-rows", matrix4, 5, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 1 with row", matrix5, 1, "1 2 3\n" },
		{ "arg 2 with row", matrix5, 2, "1 2 3\n0 0 0\n" },
		{ "arg 3 with row", matrix5, 3, "1 2 3\n0 0 0\n0 0 0\n" },
		{ "arg 1 with three rows", matrix6, 1, "1 2 3\n" },
		{ "arg 2 with three rows", matrix6, 2, "1 2 3\n4 5 6\n" },
		{ "arg 3 with three rows", matrix6, 3, "1 2 3\n4 5 6\n7 8 9\n" },
		{ "arg 4 with three rows", matrix6, 4, "1 2 3\n4 5 6\n7 8 9\n0 0 0\n" },
		{ "arg 5 with three rows", matrix6, 5, "1 2 3\n4 5 6\n7 8 9\n0 0 0\n0 0 0\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.ResizeRows(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestResizeCols(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg    int
		want   string
	}{
		{ "arg -1", matrix2, -1, "0\n" },
		{ "arg 0", matrix2, 0, "0\n" },
		{ "arg 1 with 1 0-elem", matrix2, 1, "0\n" },
		{ "arg 2 with 1 0-elem", matrix2, 2, "0 0\n" },
		{ "arg 3 with 1 0-elem", matrix2, 3, "0 0 0\n" },
		{ "arg 1 with 0-row", matrix3, 1, "0\n" },
		{ "arg 2 with 0-row", matrix3, 2, "0 0\n" },
		{ "arg 3 with 0-row", matrix3, 3, "0 0 0\n" },
		{ "arg 1 with three 0-rows", matrix4, 1, "0\n0\n0\n" },
		{ "arg 2 with three 0-rows", matrix4, 2, "0 0\n0 0\n0 0\n" },
		{ "arg 3 with three 0-rows", matrix4, 3, "0 0 0\n0 0 0\n0 0 0\n" },
		{ "arg 4 with three 0-rows", matrix4, 4, "0 0 0 0\n0 0 0 0\n0 0 0 0\n" },
		{ "arg 5 with three 0-rows", matrix4, 5, "0 0 0 0 0\n0 0 0 0 0\n0 0 0 0 0\n" },
		{ "arg 1 with row", matrix5, 1, "1\n" },
		{ "arg 2 with row", matrix5, 2, "1 2\n" },
		{ "arg 3 with row", matrix5, 3, "1 2 3\n" },
		{ "arg 4 with row", matrix5, 4, "1 2 3 0\n" },
		{ "arg 5 with row", matrix5, 5, "1 2 3 0 0\n" },
		{ "arg 1 with three rows", matrix6, 1, "1\n4\n7\n" },
		{ "arg 2 with three rows", matrix6, 2, "1 2\n4 5\n7 8\n" },
		{ "arg 3 with three rows", matrix6, 3, "1 2 3\n4 5 6\n7 8 9\n" },
		{ "arg 4 with three rows", matrix6, 4, "1 2 3 0\n4 5 6 0\n7 8 9 0\n" },
		{ "arg 5 with three rows", matrix6, 5, "1 2 3 0 0\n4 5 6 0 0\n7 8 9 0 0\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.ResizeCols(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestResize(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
		arg1    int
		arg2 int
		want   string
	}{
		{ "arg -1, -1", matrix2, -1, -1, "0\n" },
		{ "arg 0, 0", matrix2, 0, 0, "0\n" },
		{ "arg 1, 1 with 0-elem", matrix2, 1, 1, "0\n" },
		{ "arg 2, 2 with 0-elem", matrix2, 2, 2, "0 0\n0 0\n" },
		{ "arg 1, 2 with 0-elem", matrix2, 1, 2, "0 0\n" },
		{ "arg 2, 1 with 0-elem", matrix2, 2, 1, "0\n0\n" },
		{ "arg 1, 3 with 0-row", matrix3, 1, 3, "0 0 0\n" },
		{ "arg 2, 3 with 0-row", matrix3, 2, 3, "0 0 0\n0 0 0\n" },
		{ "arg 1, 2 with 0-row", matrix3, 1, 2, "0 0\n" },
		{ "arg 1, 1 with 0-row", matrix3, 1, 1, "0\n" },
		{ "arg 2, 1 with 0-row", matrix3, 2, 1, "0\n0\n" },
		{ "arg 4, 4 with 0-row", matrix3, 4, 4, "0 0 0 0\n0 0 0 0\n0 0 0 0\n0 0 0 0\n" },
		{ "arg 1, 4 with 0-row", matrix3, 1, 4, "0 0 0 0\n" },
		{ "arg 1, 1 with three rows", matrix6, 1, 1, "1\n" },
		{ "arg 2, 2 with three rows", matrix6, 2, 2, "1 2\n4 5\n" },
		{ "arg 3, 3 with three rows", matrix6, 3, 3, "1 2 3\n4 5 6\n7 8 9\n" },
		{ "arg 4, 4 with three rows", matrix6, 4, 4, "1 2 3 0\n4 5 6 0\n7 8 9 0\n0 0 0 0\n" },
		{ "arg 1, 3 with three rows", matrix6, 1, 3, "1 2 3\n" },
		{ "arg 3, 1 with three rows", matrix6, 3, 1, "1\n4\n7\n" },
	}

	for _, test := range tests {
		t.Run(test.name,
			func(t *testing.T) {
				matrix, err := NewMatrix(test.matrix, false)
				if err != nil {
					return
				}
				matrix.Resize(test.arg1, test.arg2)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}