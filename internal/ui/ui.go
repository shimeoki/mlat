package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	
	"github.com/shimeoki/mlat/internal/matrix"
)

type GUI[number matrix.Number] struct {
	TableContainer *fyne.Container
	Table          *widget.Table
	Matrix         *matrix.Matrix[number]

	OptionsContainer *fyne.Container
	OptionsLabel     *canvas.Text
	OptionsAugmented *widget.Check
	OptionsRows      *widget.Entry
	OptionsCols      *widget.Entry
	OptionsSolution  *widget.Select

	ActionsContainer *fyne.Container
	ActionsImport    *widget.Button
	ActionsExport    *widget.Button
	ActionsCalculate *widget.Button
	ActionsCopy      *widget.Button
	ActionsAnswer    *widget.Entry
	ActionsStatus    *widget.ProgressBarInfinite

	MainContainer *fyne.Container

	Window fyne.Window
	App    fyne.App
}

func (p *GUI[number]) createTable() {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			if p.Matrix == nil {
				return 0, 0
			}
			return p.Matrix.Shape[0], p.Matrix.Shape[1]
		},
		func() fyne.CanvasObject {
			object := widget.NewEntry()
			object.Resize(fyne.NewSize(40, 20))
			return object
		},
		func(cell widget.TableCellID, object fyne.CanvasObject) {
			var text string
			if p.Matrix == nil {
				text = "nil"
			} else {
				text = fmt.Sprintf("%v", p.Matrix.Data[cell.Row][cell.Col])
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

	p.Table = table
	p.TableContainer = container.NewPadded(p.Table)
}

func NewGUI[number matrix.Number]() *GUI[number] {
	ui := &GUI[number]{}
	ui.Matrix = nil

	ui.App = app.New()
	ui.Window = ui.App.NewWindow("mlat")

	ui.createTable()
	ui.createOptions()
	ui.createActions()
	ui.MainContainer = container.NewBorder(
		nil, ui.ActionsContainer, ui.OptionsContainer, nil, ui.TableContainer,
	)
	ui.Window.SetContent(ui.MainContainer)

	return ui
}

func (p *GUI[_]) Run() {
	p.Window.ShowAndRun()
}

func (p *GUI[_]) createOptions() {
	p.OptionsContainer = container.NewVBox()

	p.OptionsLabel = canvas.NewText("Options", theme.ForegroundColor())
	p.OptionsLabel.TextStyle.Bold = true
	p.OptionsLabel.TextSize = 24
	p.OptionsContainer.Add(p.OptionsLabel)

	p.OptionsAugmented = widget.NewCheck("Augmented", func(bool) {})
	p.OptionsAugmented.Disable()
	p.OptionsContainer.Add(p.OptionsAugmented)

	validator := func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	p.OptionsRows = widget.NewEntry()
	p.OptionsRows.SetPlaceHolder("Rows...")
	p.OptionsRows.Validator = validator
	p.OptionsRows.Disable()
	p.OptionsContainer.Add(p.OptionsRows)

	p.OptionsCols = widget.NewEntry()
	p.OptionsCols.SetPlaceHolder("Columns...")
	p.OptionsCols.Validator = validator
	p.OptionsCols.Disable()
	p.OptionsContainer.Add(p.OptionsCols)

	p.OptionsSolution = widget.NewSelect(
		[]string{"Calculate determinant(s)"},
		func(string) {},
	)
	p.OptionsSolution.SetSelectedIndex(0)
	p.OptionsSolution.Disable()
	p.OptionsContainer.Add(p.OptionsSolution)

	p.OptionsContainer = container.NewPadded(p.OptionsContainer)
}

func (p *GUI[_]) createActions() {
	p.ActionsContainer = container.NewGridWithRows(1)

	p.ActionsImport = widget.NewButtonWithIcon(
		"Import Matrix",
		theme.UploadIcon(),
		func() {},
	)

	p.ActionsExport = widget.NewButtonWithIcon(
		"Export Matrix",
		theme.DownloadIcon(),
		func() {},
	)

	p.ActionsCalculate = widget.NewButtonWithIcon(
		"Calculate",
		theme.GridIcon(),
		func() {},
	)

	p.ActionsAnswer = widget.NewEntry()
	p.ActionsAnswer.SetPlaceHolder("Answer...")
	p.ActionsAnswer.Disable()

	p.ActionsCopy = widget.NewButtonWithIcon(
		"",
		theme.ContentCopyIcon(),
		func() {},
	)

	p.ActionsStatus = widget.NewProgressBarInfinite()
	p.ActionsStatus.Stop()

	p.ActionsContainer.Add(p.ActionsImport)
	p.ActionsContainer.Add(p.ActionsExport)
	p.ActionsContainer.Add(p.ActionsCalculate)
	p.ActionsContainer.Add(p.ActionsCopy)
	p.ActionsContainer.Add(p.ActionsAnswer)
	p.ActionsContainer.Add(p.ActionsStatus)
	p.ActionsContainer = container.NewPadded(p.ActionsContainer)
}
