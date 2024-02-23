package ui

import (
	"fmt"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shimeoki/mlat/internal/matrix"
)

func newTable[number matrix.Number](mx *matrix.Matrix[number]) fyne.Widget {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return mx.Shape[0], mx.Shape[1]
		},
		func() fyne.CanvasObject {
			object := widget.NewEntry()
			object.Resize(fyne.NewSize(40, 20))
			return object
		},
		func(cell widget.TableCellID, object fyne.CanvasObject) {
			object.(*widget.Entry).SetText(fmt.Sprintf("%v", mx.Data[cell.Row][cell.Col]))
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
			text = fmt.Sprint(id.Col+1)
		} else {
			text = fmt.Sprint(id.Row+1)
		}

		object.(*widget.Label).SetText(text)
	}


	return table
}

func RunUI[number matrix.Number](mx *matrix.Matrix[number]) {
	a := app.New()
	w := a.NewWindow("mlat")

	// table := NewTable(mx)
	// table.Resize(fyne.NewSize(440, 440))

	// tableContainer := container.NewPadded()

	// tableContainer.Resize(fyne.NewSize(440, 440))

	// mainContainer := container.NewWithoutLayout(tableContainer)
	// mainContainer.Resize(fyne.NewSize(1000, 1000))

	w.SetContent(getContent(mx))
	w.ShowAndRun()
}

func getOptions() *fyne.Container {
	options := container.NewVBox()

	label := canvas.NewText("Options", theme.ForegroundColor())
	label.TextStyle.Bold = true
	label.TextSize = 24
	options.Add(label)

	augmented := widget.NewCheck("Augmented", func(bool) {})
	options.Add(augmented)

	validator := func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	rows := widget.NewEntry()
	rows.SetPlaceHolder("Enter Rows...")
	rows.Validator = validator
	options.Add(rows)

	cols := widget.NewEntry()
	cols.SetPlaceHolder("Enter Columns...")
	cols.Validator = validator
	options.Add(cols)

	solution := widget.NewSelect(
		[]string{"default"},
		func(string) {},
	)
	options.Add(solution)

	options.Add(layout.NewSpacer())

	return container.NewPadded(options)
}

func getTable[number matrix.Number](mx *matrix.Matrix[number]) *fyne.Container {
	table := newTable(mx)
	return container.NewPadded(table)
}

func getCenter[number matrix.Number](mx *matrix.Matrix[number]) *fyne.Container {
	options := getOptions()
	table := getTable(mx)
	center := container.NewBorder(nil, nil, options, nil, table)
	return container.NewPadded(center)
}

func getBottomRow() *fyne.Container {
	row := container.NewHBox()

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

	row.Add(layout.NewSpacer())

	answerEntry := widget.NewEntry()
	answerEntry.SetPlaceHolder("Answer")
	answerEntry.Disable()

	answerCopy := widget.NewButtonWithIcon(
		"",
		theme.ContentCopyIcon(),
		func() {},
	)

	answerStatus := widget.NewProgressBarInfinite()
	answerStatus.Stop()

	row.Add(answerCopy)
	row.Add(answerEntry)
	row.Add(answerStatus)

	return container.NewPadded(row)
}

func getContent[number matrix.Number](mx *matrix.Matrix[number]) *fyne.Container {
	center := getCenter(mx)
	bottom := getBottomRow()
	return container.NewBorder(nil, bottom, nil, nil, center)
}
