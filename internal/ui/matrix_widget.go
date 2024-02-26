package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	matrix "github.com/shimeoki/mlat/internal/cmatrix"
)

type MatrixWidget[number matrix.Number] struct {
	BaseWidget widget.BaseWidget
	Matrix     *matrix.CustomMatrix[number]
	visible    bool
	size       fyne.Size
	cellSize   float32
	position   fyne.Position
}

func NewMatrixWidget[number matrix.Number](matrix *matrix.CustomMatrix[number]) *MatrixWidget[number] {
	matrixWidget := &MatrixWidget[number]{}
	matrixWidget.Matrix = matrix

	matrixWidget.BaseWidget.ExtendBaseWidget(matrixWidget)
	return matrixWidget
}

func (p *MatrixWidget[number]) Show() {
	p.visible = true
}

func (p *MatrixWidget[number]) Hide() {
	p.visible = false
}

func (p *MatrixWidget[number]) Visible() bool {
	return p.visible
}

func (p *MatrixWidget[number]) MinSize() fyne.Size {
	rows, cols := p.Matrix.Shape[0], p.Matrix.Shape[1]
	return fyne.NewSize(p.cellSize*float32(cols), p.cellSize*float32(rows))
}

func (p *MatrixWidget[number]) Move(position fyne.Position) {

}

func (p *MatrixWidget[number]) Position() fyne.Position {
	return fyne.NewPos(0, 0)
}

func (p *MatrixWidget[number]) Refresh() {

}

func (p *MatrixWidget[number]) Resize(size fyne.Size) {

}

func (p *MatrixWidget[number]) Size() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (p *MatrixWidget[number]) CreateRenderer() fyne.WidgetRenderer {
	return &matrixRenderer[number]{matrixWidget: p}
}

type matrixRenderer[number matrix.Number] struct {
	matrixWidget *MatrixWidget[number]
}

func (p *matrixRenderer[number]) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (p *matrixRenderer[number]) Layout(size fyne.Size) {

}

func (p *matrixRenderer[number]) Refresh() {

}

func (p *matrixRenderer[number]) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}

func (p *matrixRenderer[number]) Destroy() {

}
