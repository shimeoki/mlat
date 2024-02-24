package main

import (
	"time"

	"github.com/shimeoki/mlat/internal/matrix"
	"github.com/shimeoki/mlat/internal/ui"
)

func main() {
	matrix, _ := matrix.NewMatrix[int](10, 10, false)
	matrix.FillRandom(10)

	gui := ui.NewGUI[int]()

	go func() {
		time.Sleep(time.Duration(time.Second * 2))
		gui.Matrix = matrix
		gui.Table.Refresh()
	}()

	gui.Run()
}
