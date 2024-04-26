package matrix

import "fmt"

type MatrixError struct {
	msg string
}

func (e *MatrixError) Error() string {
	return fmt.Sprintf("matrix: %s", e.msg)
}
