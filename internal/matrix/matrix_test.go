package matrix

import (
	"fmt"
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
		{ "arg 0 with blank", matrix1, 0, "" },
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

	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i+1, test.name),
			func(t *testing.T) {
				matrix, _ := NewMatrix(test.matrix, false)
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
		{ "arg -1", matrix1, -1, MatrixToString(matrix1, " ") },
		{ "arg 0 with blank", matrix1, 0, "" },
		{ "arg 0 with 1 0-elem", matrix2, 0, "0\n0\n" },
		{ "arg 0 with 0-row", matrix3, 0, "0 0 0\n0 0 0\n" },
		{ "arg 0 with three 0-rows", matrix4, 0, "0 0 0\n0 0 0\n0 0 0\n0 0 0\n"},
		{ "arg 0 with 1 elem", matrix9, 0, "0\n1\n"},
		{ "arg 0 with row", matrix5, 0, "0 0 0\n1 2 3\n"},
		{ "arg 0 with three rows", matrix6, 0, "0 0 0\n1 2 3\n4 5 6\n7 8 9\n"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i+1, test.name),
			func(t *testing.T) {
				matrix, _ := NewMatrix(test.matrix, false)
				matrix.AddRow(test.arg)
				ans := MatrixToString(matrix.Data, " ")
				if ans != test.want {
					t.Errorf("\ngot:\n%s\nwant:\n%s", ans, test.want)
				}
			})
	}
}

func TestExtendRows(t *testing.T) {

}

func TestExtendCols(t *testing.T) {

}

func TestResizeRows(t *testing.T) {

}

func TestResizeCols(t *testing.T) {

}
