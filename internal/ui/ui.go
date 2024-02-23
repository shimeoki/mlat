package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	"strconv"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shimeoki/mlat/internal/matrix"
)

type UI[number matrix.Number] struct {
	Matrix *matrix.Matrix[number]
	Table fyne.Widget
}

func createTable[number matrix.Number](mx *matrix.Matrix[number]) fyne.Widget {
	// if mx == nil {
	// 	mx, _ = matrix.NewMatrix[number](1, 1, false)
	// }

	table := widget.NewTableWithHeaders(
		func() (int, int) {
			if mx == nil {
				return 0, 0
			}
			return mx.Shape[0], mx.Shape[1]
		},
		func() fyne.CanvasObject {
			object := widget.NewEntry()
			object.Resize(fyne.NewSize(40, 20))
			return object
		},
		func(cell widget.TableCellID, object fyne.CanvasObject) {
			var text string
			if mx == nil {
				text = "nil"
			} else {
				text = fmt.Sprintf("%v", mx.Data[cell.Row][cell.Col])
			}
			object.(*widget.Entry).SetText(text)
		},
	)

	table.CreateHeader = func() fyne.CanvasObject {
		object := widget.NewLabel("header")
		object.TextStyle.Bold = true
		object.Alignment = fyne.TextAlignCenter
		return object
	}
	table.UpdateHeader = func(id widget.TableCellID, object fyne.CanvasObject) {
		var text string

		if id.Row == -1 && id.Col == -1 {
			text = ""
		} else if id.Row == -1 {
			text = fmt.Sprint(id.Col + 1)
		} else {
			text = fmt.Sprint(id.Row + 1)
		}

		object.(*widget.Label).SetText(text)
	}

	return table
}

func MakeUI[number matrix.Number]() *UI[number] {
	a := app.New()
	w := a.NewWindow("mlat")
	
	ui := &UI[number]{}
	ui.Matrix = nil
	ui.Table = createTable[number](ui.Matrix)

	tableContaner := container.NewPadded(ui.Table)
	optionsContainer := createOptions()
	actionsContainer := createActions()
	mainContainer := container.NewBorder(
		nil, actionsContainer, optionsContainer, nil, tableContaner,
	)
	w.SetContent(mainContainer)
	w.ShowAndRun()

	return ui
}

func createOptions() *fyne.Container {
	options := container.NewVBox()

	label := canvas.NewText("Options", theme.ForegroundColor())
	label.TextStyle.Bold = true
	label.TextSize = 24
	options.Add(label)

	augmented := widget.NewCheck("Augmented", func(bool) {})
	augmented.Disable()
	options.Add(augmented)

	validator := func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	rows := widget.NewEntry()
	rows.SetPlaceHolder("Rows...")
	rows.Validator = validator
	rows.Disable()
	options.Add(rows)

	cols := widget.NewEntry()
	cols.SetPlaceHolder("Columns...")
	cols.Validator = validator
	cols.Disable()
	options.Add(cols)

	solution := widget.NewSelect(
		[]string{"Calculate determinant(s)"},
		func(string) {},
	)
	solution.SetSelectedIndex(0)
	solution.Disable()
	options.Add(solution)

	// options.Add(layout.NewSpacer())

	return container.NewPadded(options)
}

func createActions() *fyne.Container {
	row := container.NewGridWithRows(1)

	importButton := widget.NewButtonWithIcon(
		"Import Matrix",
		theme.UploadIcon(),
		func() {},
	)
	row.Add(importButton)

	exportButton := widget.NewButtonWithIcon(
		"Export Matrix",
		theme.DownloadIcon(),
		func() {},
	)
	row.Add(exportButton)

	// row.Add(layout.NewSpacer())

	answerCalculate := widget.NewButtonWithIcon(
		"Calculate",
		theme.GridIcon(),
		func() {},
	)

	answerEntry := widget.NewEntry()
	answerEntry.SetPlaceHolder("Answer...")
	answerEntry.Disable()

	answerCopy := widget.NewButtonWithIcon(
		"",
		theme.ContentCopyIcon(),
		func() {},
	)

	answerStatus := widget.NewProgressBarInfinite()
	answerStatus.Stop()

	row.Add(answerCalculate)
	row.Add(answerCopy)
	row.Add(answerEntry)
	row.Add(answerStatus)

	return container.NewPadded(row)
}
