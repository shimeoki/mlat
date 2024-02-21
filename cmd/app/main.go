package main

import (
	"fmt"
	// "time"
	"github.com/shimeoki/mlat/internal/matrix"
)

func main() {
	testRead()
	// matrix, _ := matrix.NewMatrix[int](13, 14, true)
	// matrix.FillRandom(10)

	// fmt.Println(matrix)
	// fmt.Println()

	// start := time.Now()
	// matrix.Calculate()
	// fmt.Println(matrix.GetRoots())
	// fmt.Printf("Time spent: %s", time.Since(start))
	// fmt.Println(matrix.GetRoots())
}

func testRead() {
	matrix, err := matrix.ReadSlow("test/01.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(matrix)
}