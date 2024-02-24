package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/shimeoki/mlat/internal/matrix"
)

type GUI struct {
	TableContainer *fyne.Container
	Table          *widget.Table
	Matrix         *matrix.Matrix[float64]

	OptionsContainer *fyne.Container
	OptionsLabel     *canvas.Text
	OptionsAugmented *widget.Check
	OptionsRows      *widget.Entry
	OptionsCols      *widget.Entry
	OptionsSolution  *widget.Select

	ActionsContainer    *fyne.Container
	ActionsImport       *widget.Button
	ActionsImportDialog *dialog.FileDialog
	ActionsExport       *widget.Button
	ActionsExportDialog *dialog.FileDialog
	ActionsCalculate    *widget.Button
	ActionsCopy         *widget.Button
	ActionsAnswer       *widget.Entry
	ActionsStatus       *widget.ProgressBarInfinite

	MainContainer *fyne.Container

	Window fyne.Window
	App    fyne.App
}

func (p *GUI) createTable() {
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
			if id.Col+1 == p.Matrix.Shape[1] && p.Matrix.Augmented {
				text = ""
			} else {
				text = fmt.Sprint(id.Col + 1)
			}
		} else {
			text = fmt.Sprint(id.Row + 1)
		}

		object.(*widget.Label).SetText(text)
	}

	p.Table = table
	p.TableContainer = container.NewPadded(p.Table)
}

func NewGUI() *GUI {
	gui := &GUI{}
	gui.Matrix = nil

	gui.App = app.New()
	gui.Window = gui.App.NewWindow("mlat")

	gui.createTable()
	gui.createOptions()
	gui.createActions()
	gui.MainContainer = container.NewBorder(
		nil, gui.ActionsContainer, gui.OptionsContainer, nil, gui.TableContainer,
	)
	gui.Window.SetContent(gui.MainContainer)

	return gui
}

func (p *GUI) Run() {
	p.Window.ShowAndRun()
}

func (p *GUI) createOptions() {
	p.OptionsContainer = container.NewVBox()

	p.OptionsLabel = canvas.NewText("Options", theme.ForegroundColor())
	p.OptionsLabel.TextStyle.Bold = true
	p.OptionsLabel.TextSize = 24
	p.OptionsContainer.Add(p.OptionsLabel)

	p.OptionsAugmented = widget.NewCheck(
		"Augmented", 
		func(state bool) {
			if p.Matrix == nil {
				return
			}

			p.Matrix, _ = matrix.NewMatrix(p.Matrix.Data, state)
			p.Table.Refresh()
		},
	)
	// p.OptionsAugmented.Disable()
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

func (p *GUI) createActions() {
	p.ActionsContainer = container.NewGridWithRows(1)

	p.ActionsImportDialog = dialog.NewFileOpen(
		func(uri fyne.URIReadCloser, err error) {
			if uri == nil || err != nil {
				return
			}

			mx, err := matrix.ReadSlow(uri.URI().Path())
			if err != nil {
				return
			}

			p.Matrix, _ = matrix.NewMatrix(mx, false)
			p.OptionsRows.SetText(fmt.Sprint(p.Matrix.Shape[0]))
			p.OptionsCols.SetText(fmt.Sprint(p.Matrix.Shape[1]))
			p.Table.Refresh()
		},
		p.Window,
	)
	p.ActionsImport = widget.NewButtonWithIcon(
		"Import Matrix",
		theme.UploadIcon(),
		func() {
			p.ActionsImportDialog.Show()
		},
	)

	p.ActionsExportDialog = dialog.NewFileSave(
		func(uri fyne.URIWriteCloser, err error) {
			if uri == nil || err != nil {
				return
			}

			matrix.Write(uri.URI().Path(), p.Matrix.Data)
		},
		p.Window,
	)
	p.ActionsExport = widget.NewButtonWithIcon(
		"Export Matrix",
		theme.DownloadIcon(),
		func() {
			p.ActionsExportDialog.Show()
		},
	)

	p.ActionsCalculate = widget.NewButtonWithIcon(
		"Calculate",
		theme.GridIcon(),
		func() {
			p.ActionsStatus.Start()
			defer p.ActionsStatus.Stop()
			answer, err := p.Matrix.Calculate()
			if err != nil {
				p.ActionsStatus.Stop()
				dialog.ShowInformation(
					"Error!",
					"Matrix is not a square.",
					p.Window,
				)
				return
			}

			if len(answer) == 1 {
				p.ActionsAnswer.SetText(fmt.Sprintf("%f", answer[0]))
			} else {
				p.ActionsAnswer.SetText(
					matrix.ArrayToString(p.Matrix.GetRoots(), " "),
				)
			}
		},
	)

	p.ActionsAnswer = widget.NewEntry()
	p.ActionsAnswer.SetPlaceHolder("Answer...")
	p.ActionsAnswer.Disable()

	p.ActionsCopy = widget.NewButtonWithIcon(
		"",
		theme.ContentCopyIcon(),
		func() {
			p.Window.Clipboard().SetContent(
				p.ActionsAnswer.Text,
			)
			p.ActionsCopy.Icon = theme.ConfirmIcon()
			p.ActionsCopy.Refresh()

			time.Sleep(1 * time.Second)
			p.ActionsCopy.Icon = theme.ContentCopyIcon()
			p.ActionsCopy.Refresh()
		},
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
