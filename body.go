package Tablam

import (
	"errors"
	// "fmt"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Tablam is a gtk grid with a THeader and an slice of TRow
type Tablam struct {
	*gtk.Grid
	header *THeader
	erows  []*TRow
}

func NewTablam(h *THeader, rows, cols, width int) (*Tablam, error) {
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(RowSep)
	grid.AddEvents(int(gdk.SCROLL_MASK))

	if h == nil {
		gData = newTData(cols, width)
		gTheme = newTheme()
		gTheme.setFontBold('c', true)

		gData.widths = make([]int, cols)
		for i := 0; i < cols; i++ {
			gData.widths[i] = width
		}

	} else {
		if len(h.titles) != cols {
			return nil, errors.New(
				"New Tablam: Different number of columns and (header) aligns")
		}
		copy(gData.widths, h.widths)
	}

	var erows []*TRow

	for i := 0; i < rows; i++ {
		erow := newTRow(nil, i)
		erows = append(erows, erow)
		grid.Attach(erow, 0, i, 1, 1)
	}

	if h != nil {
		h.widthsRecalc()
	}

	return &Tablam{grid, h, erows}, nil
}

func (t *Tablam) FillWithData(d [][]string, aligns []MyAlign) error {
	if len(d) > 0 && len(d[0]) != len(aligns) {
		return errors.New(
			"Fill Data: Different number of columns and aligns.")
	}
	if len(d) == 0 || len(t.erows) == 0 {
		return nil
	}

	gData.setAligns(aligns)
	gData.createDataRows(d)

	for i := range gData.widths {
		gData.recalcColumnWidth(i, t.header)
	}
	gData.recalcXtitles()
	gData.setEmpties()

	if SortingColumn >= 0 {
		sortData()
	}

	if IndexColumn >= 0 {
		indexData()
	}

	t.refreshData(0)

	if t.header != nil {
		t.header.widthsRecalc()
	}

	return nil
}

func (t Tablam) DataRows() int {
	return len(gData.drows)
}

func (t Tablam) TotalRows() int {
	return len(t.erows)
}

func (t *Tablam) ShowOrHideColumn(c int) {
	if t.header.hidden[c] {
		t.header.hidden[c] = false
		t.header.eventBoxes[c].Show()
		for _, erow := range t.erows {
			erow.tboxes[c].Show()
		}
	} else {
		t.header.hidden[c] = true
		t.header.eventBoxes[c].Hide()
		for _, erow := range t.erows {
			erow.tboxes[c].Hide()
		}
	}
}

func (t *Tablam) CursorUp() {
	dc := false
	if Position >= 0 {
		t.clearCursor()
	}
	if Position > 0 {
		Position--
		dc = true
	} else {
		if len(t.erows) > 0 {
			Position = len(t.erows) - 1
			dc = true
		}
	}
	if dc {
		t.drawCursor()
	}
}

func (t *Tablam) CursorDown() {
	if len(t.erows) == 0 {
		return
	}
	if Position >= 0 {
		t.clearCursor()
	}
	if Position < len(t.erows)-1 {
		Position++
	} else {
		Position = 0
	}
	t.drawCursor()
}

func (t *Tablam) RemoveCursor() {
	t.clearCursor()
	Position = -1
}

func (t *Tablam) UpdateCursor(withSel bool) {
	if withSel {
		oldPosition = Position
		Position = selectedRow
	}

	if oldPosition >= 0 && oldPosition != Position {
		op := oldPosition
		t.erows[op].refreshNormalMarkup(op)
	}

	if Position >= 0 && Position < len(t.erows) {
		p := Position
		t.erows[p].refreshCursorMarkup(p)
	}
}

func (t *Tablam) clearCursor() {
	if Position >= 0 {
		r := Position
		t.erows[r].refreshNormalMarkup(r)
	}
}

func (t *Tablam) clearSelect() {
	if Position >= 0 {
		r := Position
		t.erows[r].refreshNormalMarkup(r)
	}
}

func (t *Tablam) drawCursor() {
	if Position >= 0 {
		r := Position
		t.erows[r].refreshCursorMarkup(r)
	}
}

func (t Tablam) ActiveRowData() []string {
	vo := VerticalOffset
	if Position >= 0 && Position < len(gData.drows)-vo {
		return gData.drows[Position+vo].titles
	}
	return nil
}

func (t Tablam) selectRow(row int) int {
	vo := VerticalOffset
	if row >= 0 && row < len(gData.drows)-vo {
		alreadySel := false
		for i, sel := range Selection {
			if sel == row {
				alreadySel = true
				Selection = append(Selection[:i], Selection[i+1:]...)
				break
			}
		}
		if !alreadySel {
			Selection = append(Selection, row)
			return row
		}
	}
	return -1
}

func (t Tablam) SelectCursorRow() int {
	return t.selectRow(Position)
}

func (t Tablam) SelectARow() {
	vo := VerticalOffset
	row := selectedRow
	if t.selectRow(row) < 0 {
		if row != Position {
			for col, tbox := range t.erows[row].tboxes {
				gTheme.normalMarkup(tbox.label, &gData.drows[row+vo].xtitles[col])
			}
		}
	} else {
		if row != Position {
			for col, tbox := range t.erows[row].tboxes {
				gTheme.selectMarkup(tbox.label, &gData.drows[row+vo].xtitles[col])
			}
		}
	}
}

func (t Tablam) ClearSelected() {
	vo := VerticalOffset
	for _, sel := range Selection {
		if sel != Position {
			for col, tbox := range t.erows[sel].tboxes {
				gTheme.normalMarkup(tbox.label, &gData.drows[sel+vo].xtitles[col])
			}
		}
	}
	Selection = []int{}
}

func (t Tablam) FontSize() int {
	return gTheme.themeFontSize()
}

func (t *Tablam) DecFontSize() {
	if gTheme.decFontSize() {
		if t.header != nil {
			t.header.refreshMarkups()
		}
		t.refreshData(0)
	}
}

func (t *Tablam) IncFontSize() {
	gTheme.incFontSize()
	if t.header != nil {
		t.header.refreshMarkups()
	}
	t.refreshData(0)
}

func (t Tablam) HeaderFontBold() bool {
	if gTheme.fontBold('h') {
		return true
	}
	return false
}
func (t Tablam) NormalFontBold() bool {
	if gTheme.fontBold('n') {
		return true
	}
	return false
}
func (t Tablam) CursorFontBold() bool {
	if gTheme.fontBold('c') {
		return true
	}
	return false
}
func (t Tablam) SelectFontBold() bool {
	if gTheme.fontBold('s') {
		return true
	}
	return false
}

func (t Tablam) SetHeaderFontBold(b bool) {
	gTheme.setFontBold('h', b)
	if t.header != nil {
		t.header.refreshMarkups()
	}
}

func (t Tablam) SetNormalFontBold(b bool) {
	gTheme.setFontBold('n', b)
	t.refreshData(0)
}

func (t Tablam) SetCursorFontBold(b bool) {
	gTheme.setFontBold('c', b)
	t.refreshData(0)
}

func (t Tablam) SetSelectFontBold(b bool) {
	gTheme.setFontBold('s', b)
	t.refreshData(0)
}

func (t *Tablam) SetCursorColors(fg, bg string) {
	gTheme.setColors('c', fg, bg)
	t.drawCursor()
}

func (t *Tablam) SetNormalColors(fg, bg string) {
	gTheme.setColors('n', fg, bg)
	t.refreshData(0)
}

func (t Tablam) ResetPosition() {
	oldPosition = -1
	Position = -1
}

func (t *Tablam) refreshData(n int) {
	// for r, erow := range t.erows {
	for r := n; r < len(t.erows); r++ {
		if r == Position {
			t.erows[r].refreshCursorMarkup(r)
		} else {
			t.erows[r].refreshNormalMarkup(r)
		}
	}
}

func (t *Tablam) UpdateWithOffset(n int) {
	Position -= n
	t.refreshData(0)
}

func (t *Tablam) UpdateCell(r, c int, title string) error {
	if r >= len(gData.drows) {
		return errors.New("Update Cell: Row out of range.")
	}
	if c >= len(gData.widths) {
		return errors.New("Update Cell: Column out of range.")
	}

	old := utf8.RuneCountInString(gData.drows[r].titles[c])
	new := utf8.RuneCountInString(title)
	gData.drows[r].titles[c] = title
	gData.recalcXtitle(r, c, &title)

	if old < gData.widths[c] && new <= gData.widths[c] {
		t.erows[r].tboxes[c].refreshMarkup(r, c)
	} else {
		gData.recalcColumnWidth(c, t.header)
		gData.recalcColumnXtitles(c)
		t.refreshColumn(c)
		if t.header != nil {
			t.header.widthsRecalc()
		}
	}
	return nil
}

func (t *Tablam) refreshColumn(c int) {
	vo := VerticalOffset
	for r, erow := range t.erows {
		if r < len(gData.drows)-vo {
			if r == Position {
				gTheme.cursorMarkup(erow.tboxes[c].label, &gData.drows[r].xtitles[c])
			} else {
				gTheme.normalMarkup(erow.tboxes[c].label, &gData.drows[r].xtitles[c])
			}
		} else {
			if r == Position {
				gTheme.cursorMarkup(erow.tboxes[c].label, &gData.empties[c])
			} else {
				gTheme.normalMarkup(erow.tboxes[c].label, &gData.empties[c])
			}
		}
	}
}

func (t *Tablam) UpdateRow(r int, data []string) {
	for c, str := range data {
		t.UpdateCell(r, c, str)
	}
}

func (t *Tablam) addNewRow(n int, row []string) {
	if n < 0 {
		gData.drows = append(gData.drows, newDRow(n, row, gData.widths, gData.aligns))
	} else {
		gData.drows = append(gData.drows[:n+1], gData.drows[n:]...)
		gData.drows[n] = newDRow(n, row, gData.widths, gData.aligns)
	}

	for i := VerticalOffset; i < len(gData.drows); i++ {
		if i < len(t.erows) {
			t.UpdateRow(i, gData.drows[i].titles)
		}
	}
}

// func (t *Tablam) addEmptyRow() {
// 	last := len(t.erows)
// 	t.erows = append(t.erows, newTRow(nil, len(t.erows)))
// 	t.Attach(t.erows[last], 0, last, 1, 1)
// 	t.ShowAll()
// }

func (t *Tablam) AddDataRow(row []string) {
	t.addNewRow(-1, row)
}

func (t *Tablam) InsertARow(n int, row []string) {
	if len(gData.drows) > 0 && n < len(gData.drows) {
		t.addNewRow(n, row)
	}
}

func (t *Tablam) RemoveARow(r int) {
	if r >= 0 && r < len(gData.drows) {
		gData.drows = append(gData.drows[:r], gData.drows[r+1:]...)

		if !PageMode {
			t.erows[r].remove()
			t.erows = append(t.erows[:r], t.erows[r+1:]...)
			t.RemoveRow(r)
		}

		if len(gData.drows) == 0 {
			Position = -1
		} else if Position == len(gData.drows) {
			Position--
			t.drawCursor()
		}

		t.refreshData(r)
	}
}

func (t *Tablam) SortData(n int) {
	if n != SortingColumn {
		SortingColumn = n
	}
	sortData()
	t.refreshData(0)
}
