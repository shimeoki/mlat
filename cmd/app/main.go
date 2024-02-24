package main

import (
	"github.com/shimeoki/mlat/internal/matrix"
	"github.com/shimeoki/mlat/internal/ui"
)

func main() {
	gui := ui.NewGUI()

	mx, _ := matrix.NewBlankMatrix[int](10, 11, true)
	mx.FillRandom(10)
	matrix.Write("./test/04.txt", mx.Data)

	gui.Run()
}
