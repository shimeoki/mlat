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
	mx, err := matrix.ReadSlow("test/01.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	matrix.Write("test/01-copy.txt", mx)
	mx, err = matrix.ReadSlow("test/01-copy.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(mx)
}