package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	matrix "github.com/shimeoki/mlat/internal/cmatrix"
)

type MatrixWidget struct {
	BaseWidget widget.BaseWidget
	Matrix     *matrix.CustomMatrix
	visible    bool
	size       fyne.Size
	cellSize   float32
	position   fyne.Position
}

func NewMatrixWidget(matrix *matrix.CustomMatrix) *MatrixWidget {
	matrixWidget := &MatrixWidget{}
	matrixWidget.Matrix = matrix

	matrixWidget.BaseWidget.ExtendBaseWidget(matrixWidget)
	return matrixWidget
}

func (p *MatrixWidget) Show() {
	p.visible = true
}

func (p *MatrixWidget) Hide() {
	p.visible = false
}

func (p *MatrixWidget) Visible() bool {
	return p.visible
}

func (p *MatrixWidget) MinSize() fyne.Size {
	rows, cols := p.Matrix.Shape[0], p.Matrix.Shape[1]
	return fyne.NewSize(p.cellSize*float32(cols), p.cellSize*float32(rows))
}

func (p *MatrixWidget) Move(position fyne.Position) {

}

func (p *MatrixWidget) Position() fyne.Position {
	return fyne.NewPos(0, 0)
}

func (p *MatrixWidget) Refresh() {

}

func (p *MatrixWidget) Resize(size fyne.Size) {

}

func (p *MatrixWidget) Size() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (p *MatrixWidget) CreateRenderer() fyne.WidgetRenderer {
	return &matrixRenderer{matrixWidget: p}
}

type matrixRenderer struct {
	matrixWidget *MatrixWidget
}

func (p *matrixRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (p *matrixRenderer) Layout(size fyne.Size) {

}

func (p *matrixRenderer) Refresh() {

}

func (p *matrixRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}

func (p *matrixRenderer) Destroy() {

}
