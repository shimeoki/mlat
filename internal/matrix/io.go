package matrix

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadSlow(path string) ([][]float64, error) {
	rows, cols, err := getRowsCols(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matrix, memory, _ := Malloc(rows, cols)
	scanner := bufio.NewScanner(file)
	for i := range rows {
		matrix[i] = memory[(i * cols) : (i+1)*cols]
		scanner.Scan()
		fields := strings.Fields(scanner.Text())
		for j, field := range fields {
			matrix[i][j], _ = strconv.ParseFloat(field, 64)
		}
	}

	return matrix, nil
}

func Write(path string, matrix [][]float64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	file.WriteString(MatrixToString(matrix, " "))

	return nil
}

func getLinesCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	count := 0
	separator := []byte{'\n'}

	for {
		n, err := file.Read(buffer)
		count += bytes.Count(buffer[:n], separator)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func getRowsCols(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	rows, cols := 0, 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if cols == 0 {
			cols = len(fields)
		} else if cols != len(fields) {
			return 0, 0, errors.New("file is not a matrix")
		}

		rows++
	}

	return rows, cols, nil
}

func ArrayToString(array []float64, separator string) string {
	return strings.Trim(
		strings.Replace(fmt.Sprint(array), " ", separator, -1), "[]",
	)
}

func MatrixToString(matrix [][]float64, separator string) string {
	if matrix == nil {
		return ""
	} else if len(matrix) == 0 {
		return ""
	}

	var sb strings.Builder

	for _, row := range matrix {
		sb.WriteString(ArrayToString(row, " "))
		sb.WriteByte('\n')
	}

	return sb.String()
}
