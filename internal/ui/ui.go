package ui

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

type MultiplyTab struct {
	Common binding.Int
	Rows   binding.Int
	Cols   binding.Int

	MatrixA         *cmatrix.Matrix
	TableA          *widget.Table
	TableAContainer *fyne.Container

	MatrixB         *cmatrix.Matrix
	TableB          *widget.Table
	TableBContainer *fyne.Container

	MatrixResult         *cmatrix.Matrix
	TableResult          *widget.Table
	TableResultContainer *fyne.Container

	ActionsCommon          *widget.Entry
	ActionsCommonContainer *fyne.Container

	ActionsRows          *widget.Entry
	ActionsRowsContainer *fyne.Container

	ActionsCols          *widget.Entry
	ActionsColsContainer *fyne.Container

	ActionsOptions *fyne.Container

	ActionsImportA          *widget.Button
	ActionsImportAContainer *fyne.Container

	ActionsImportB          *widget.Button
	ActionsImportBContainer *fyne.Container

	ActionsCalculate          *widget.Button
	ActionsCalculateContainer *fyne.Container

	ActionsContainer *fyne.Container

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
			if *matrix == nil {
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
	multiplyTab := gui.newMultiplyTab()
	gui.Tabs = container.NewAppTabs(
		container.NewTabItem("Determinant", determinantTab.MainContainer),
		container.NewTabItem("Multiply", multiplyTab.MainContainer),
	)

	gui.Window.SetContent(gui.Tabs)

	return gui
}

func (p *GUI) Run() {
	p.Window.ShowAndRun()
}

// TODO change behaviour
func (p *GUI) newDeterminantTab() *DeterminantTab {
	tab := &DeterminantTab{}
	tab.GUI = p

	tab.Matrix = nil

	tab.Table = createTable(&tab.Matrix)
	tab.TableContainer = container.NewPadded(tab.Table)
	tab.createOptions()
	tab.createActions()
	tab.MainContainer = container.NewBorder(
		nil, tab.ActionsContainer, tab.OptionsContainer, nil, tab.TableContainer,
	)

	return tab
}

// TODO change behaviour
func (p *GUI) newMultiplyTab() *MultiplyTab {
	tab := &MultiplyTab{}
	tab.GUI = p

	tab.Common = binding.NewInt()
	tab.Common.Set(1)
	tab.Rows = binding.NewInt()
	tab.Rows.Set(1)
	tab.Cols = binding.NewInt()
	tab.Cols.Set(1)

	tab.MatrixA, _ = cmatrix.NewBlankMatrix(1, 1, false)
	tab.MatrixB, _ = cmatrix.NewBlankMatrix(1, 1, false)
	tab.MatrixResult, _ = cmatrix.NewBlankMatrix(1, 1, false)

	tab.TableA = createTable(&tab.MatrixA)
	tab.TableB = createTable(&tab.MatrixB)
	tab.TableResult = createTable(&tab.MatrixResult)

	tab.TableAContainer = container.NewPadded(container.NewBorder(
		container.NewCenter(widget.NewLabelWithStyle(
			"Matrix A", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})), nil, nil, nil, tab.TableA))
	tab.TableBContainer = container.NewPadded(container.NewBorder(
		container.NewCenter(widget.NewLabelWithStyle(
			"Matrix B", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})), nil, nil, nil, tab.TableB))
	tab.TableResultContainer = container.NewPadded(container.NewBorder(
		container.NewCenter(widget.NewLabelWithStyle(
			"Result", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})), nil, nil, nil, tab.TableResult))

	validator := func(s string) error {
		value, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		if value <= 0 || value >= 100 {
			return errors.New("error: value is not in range (0, 100]")
		}

		return nil
	}

	tab.ActionsCommon = widget.NewEntryWithData(
		binding.IntToString(tab.Common),
	)
	tab.ActionsCommon.OnChanged = func(s string) {
		if tab.ActionsCommon.Validate() != nil {
			return
		}

		common, _ := strconv.Atoi(s)
		tab.MatrixA.ResizeCols(common)
		tab.MatrixB.ResizeRows(common)
		tab.TableA.Refresh()
		tab.TableB.Refresh()
	}
	tab.ActionsCommon.Validator = validator
	tab.ActionsCommonContainer = container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, nil, container.NewCenter(
			widget.NewLabelWithStyle("Common: ", fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))),
		container.NewPadded(tab.ActionsCommon))

	tab.ActionsRows = widget.NewEntryWithData(
		binding.IntToString(tab.Rows),
	)
	tab.ActionsRows.OnChanged = func(s string) {
		if tab.ActionsRows.Validate() != nil {
			return
		}

		rows, _ := strconv.Atoi(s)
		tab.MatrixA.ResizeRows(rows)
		tab.MatrixResult.ResizeRows(rows)
		tab.TableA.Refresh()
		tab.TableResult.Refresh()
	}
	tab.ActionsRows.Validator = validator
	tab.ActionsRowsContainer = container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, nil, container.NewCenter(
			widget.NewLabelWithStyle("Rows: ", fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true}))),
		container.NewPadded(tab.ActionsRows))

	tab.ActionsCols = widget.NewEntryWithData(
		binding.IntToString(tab.Cols),
	)
	tab.ActionsCols.OnChanged = func(s string) {
		if tab.ActionsCols.Validate() != nil {
			return
		}

		cols, _ := strconv.Atoi(s)
		tab.MatrixB.ResizeCols(cols)
		tab.MatrixResult.ResizeCols(cols)
		tab.TableB.Refresh()
		tab.TableResult.Refresh()
	}
	tab.ActionsCols.Validator = validator
	tab.ActionsColsContainer = container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, nil, container.NewCenter(
			widget.NewLabelWithStyle("Cols: ", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}))),
		container.NewPadded(tab.ActionsCols))

	tab.ActionsOptions = container.NewHBox(container.NewVBox(
		tab.ActionsCommonContainer,
		tab.ActionsRowsContainer,
		tab.ActionsColsContainer,
	))

	tab.ActionsContainer = container.NewPadded(
		tab.ActionsOptions,
	)

	tab.MainContainer = container.NewAdaptiveGrid(
		2,
		tab.TableAContainer,
		tab.TableBContainer,
		tab.TableResultContainer,
		tab.ActionsContainer,
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

func (p *GUI) createImportButton(text string, matrix **cmatrix.Matrix) (button *widget.Button) {
	button = widget.NewButtonWithIcon(
		text,
		theme.UploadIcon(),
		func() {
			dialog.NewFileOpen(
				func(uri fyne.URIReadCloser, err error) {
					if uri == nil || err != nil {
						return
					}

					mx, err := cmatrix.ReadSlow(uri.URI().Path())
					if err != nil {
						return
					}

					*matrix, _ = cmatrix.NewMatrix(mx, false)
				},
				p.Window,
			).Show()
		},
	)
	return
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
