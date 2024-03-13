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

	cmatrix "github.com/shimeoki/mlat/internal/matrix"
)

type GUI struct {
	Tabs *container.AppTabs

	Window fyne.Window
	App    fyne.App
}

type DeterminantTab struct {
	TableContainer *fyne.Container
	Table          *widget.Table
	Matrix         *cmatrix.Matrix

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

	GUI *GUI
}

func createTable(matrix **cmatrix.Matrix) (table *widget.Table) {
	table = widget.NewTableWithHeaders(
		func() (int, int) {
			if *matrix == nil {
				return 0, 0
			}
			return (*matrix).Rows, (*matrix).Cols
		},
		func() fyne.CanvasObject {
			cell := widget.NewEntry()
			cell.Resize(fyne.NewSize(40, 20))
			return cell
		},
		func(cellID widget.TableCellID, cell fyne.CanvasObject) {
			var text string
			if matrix == nil {
				text = "nil"
			} else {
				text = fmt.Sprintf("%v", (*matrix).Data[cellID.Row][cellID.Col])
			}
			cell.(*widget.Entry).SetText(text)
		},
	)
	table.CreateHeader = func() fyne.CanvasObject {
		header := widget.NewLabel("header")
		header.TextStyle.Bold = true
		header.Alignment = fyne.TextAlignCenter
		return header
	}
	table.UpdateHeader = func(cellID widget.TableCellID, cell fyne.CanvasObject) {
		var text string

		if cellID.Row == -1 && cellID.Col == -1 {
			text = ""
		} else if cellID.Row == -1 {
			if cellID.Col+1 == (*matrix).Cols && (*matrix).Augmented {
				text = ""
			} else {
				text = fmt.Sprint(cellID.Col + 1)
			}
		} else {
			text = fmt.Sprint(cellID.Row + 1)
		}

		cell.(*widget.Label).SetText(text)
	}

	return
}

func NewGUI() *GUI {
	gui := &GUI{}

	gui.App = app.New()
	gui.Window = gui.App.NewWindow("mlat")

	determinantTab := gui.newDeterminantTab()
	gui.Tabs = container.NewAppTabs(
		container.NewTabItem("Determinant", determinantTab.MainContainer),
	)

	gui.Window.SetContent(gui.Tabs)

	return gui
}

func (p *GUI) Run() {
	p.Window.ShowAndRun()
}

func (p *GUI) newDeterminantTab() *DeterminantTab {
	tab := &DeterminantTab{}
	tab.Matrix = nil

	tab.GUI = p

	tab.Table = createTable(&tab.Matrix)
	tab.TableContainer = container.NewPadded(tab.Table)
	tab.createOptions()
	tab.createActions()
	tab.MainContainer = container.NewBorder(
		nil, tab.ActionsContainer, tab.OptionsContainer, nil, tab.TableContainer,
	)

	return tab
}

func (p *DeterminantTab) createOptions() {
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

			p.Matrix, _ = cmatrix.NewMatrix(p.Matrix.Data, state)
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

func (p *DeterminantTab) createActions() {
	p.ActionsContainer = container.NewGridWithRows(1)

	p.ActionsImportDialog = dialog.NewFileOpen(
		func(uri fyne.URIReadCloser, err error) {
			if uri == nil || err != nil {
				return
			}

			mx, err := cmatrix.ReadSlow(uri.URI().Path())
			if err != nil {
				return
			}

			p.Matrix, _ = cmatrix.NewMatrix(mx, false)
			p.OptionsRows.SetText(fmt.Sprint(p.Matrix.Rows))
			p.OptionsCols.SetText(fmt.Sprint(p.Matrix.Cols))
			p.Table.Refresh()
		},
		p.GUI.Window,
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

			cmatrix.Write(uri.URI().Path(), p.Matrix.Data)
		},
		p.GUI.Window,
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
					p.GUI.Window,
				)
				return
			}

			if len(answer) == 1 {
				p.ActionsAnswer.SetText(fmt.Sprintf("%f", answer[0]))
			} else {
				p.ActionsAnswer.SetText(
					cmatrix.ArrayToString(p.Matrix.GetRoots(), " "),
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
			p.GUI.Window.Clipboard().SetContent(
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
