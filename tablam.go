package Tablam

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var RowSep uint = 1
var ColumnSep uint = 1
var LeftRightMargin = 1

var identMaker = 1000
var oldPosition = -1
var Position = -1
var VerticalOffset = 0

var IndexColumn = -1
var SortingColumn = -1
var ReverseSorting bool
var PageMode = false
var ActiveHeaderCell = 0

var Selection []int
var selectedRow = -1

type MyAlign uint8

const (
	AlignLeft MyAlign = iota + 1
	AlignCenter
	AlignRight
	AlignNone
)

var gTheme Theme

type DRow struct {
	index   int
	ident   int
	titles  []string
	xtitles []string
}

func newDRow(index int, titles []string, widths []int, aligns []MyAlign) DRow {
	ident := identMaker
	identMaker++
	var xtitles []string

	for i, title := range titles {
		grow := widths[i] - utf8.RuneCountInString(title)
		xtitle := generateX(&title, aligns[i], grow)
		xtitles = append(xtitles, xtitle)
	}

	return DRow{
		index,
		ident,
		titles,
		xtitles,
	}
}

func (d *DRow) recalcXtitle(title *string, c, width int, align MyAlign) {
	grow := width - utf8.RuneCountInString(*title)
	d.xtitles[c] = generateX(title, align, grow)
	d.titles[c] = *title
}

// TData holds the original data
type TData struct {
	drows   []DRow
	aligns  []MyAlign
	widths  []int
	empties []string
}

func newTData(cols, width int) TData {
	var widths []int
	for i := 0; i < cols; i++ {
		widths = append(widths, width)
	}
	return TData{
		drows:   nil,
		aligns:  make([]MyAlign, cols),
		widths:  widths,
		empties: make([]string, cols),
	}
}

func (d *TData) setAligns(aligns []MyAlign) {
	for i := range aligns {
		d.aligns[i] = aligns[i]
	}
}

func (d *TData) setWidths(widths []int) {
	for i := range widths {
		d.widths[i] = widths[i]
	}
}

func (d *TData) setDrows(data [][]string) {
	for i, drow := range d.drows {
		drow.titles = data[i]
	}
}

func (d *TData) setEmpties() {
	empty := ""
	for c, width := range d.widths {
		d.empties[c] = generateX(&empty, AlignNone, width)
	}
}

func (d *TData) createDataRows(data [][]string) {
	d.drows = []DRow{}
	for i, titles := range data {
		drow := newDRow(i, titles, d.widths, d.aligns)
		d.drows = append(d.drows, drow)
	}
}

func (d *TData) recalcEmpty(c int) {
	empty := ""
	d.empties[c] = generateX(&empty, AlignNone, d.widths[c])
}

func (d *TData) recalcColumnWidth(c int, header *THeader) bool {
	old := d.widths[c]
	new := 0
	changed := false
	if header != nil {
		new = header.widths[c]
	}
	for _, row := range d.drows {
		if utf8.RuneCountInString(row.titles[c]) > new {
			new = utf8.RuneCountInString(row.titles[c])
		}
	}
	if new != old {
		d.widths[c] = new
		changed = true
		d.recalcEmpty(c)
	}
	return changed
}

func (d *TData) recalcColumnXtitles(c int) {
	for r, row := range d.drows {
		grow := d.widths[c] - utf8.RuneCountInString(row.titles[c])
		d.drows[r].xtitles[c] = generateX(&row.titles[c], d.aligns[c], grow)
	}

	// for _, row := range d.drows {
	// 	row.tboxes[c].updateWidth(gData.widths[c], gData.aligns[c])
	// }
}

func (d *TData) recalcXtitles() {
	for r, row := range d.drows {
		for c := range row.titles {
			grow := d.widths[c] - utf8.RuneCountInString(row.titles[c])
			d.drows[r].xtitles[c] = generateX(&row.titles[c], d.aligns[c], grow)
		}
	}
}

func (d *TData) recalcXtitle(r, c int, title *string) {
	d.drows[r].recalcXtitle(title, c, d.widths[c], d.aligns[c])
}

var gData TData

// TBox represents a cell
type TBox struct {
	*gtk.EventBox
	label *gtk.Label
}

func newTBox(xtitle *string) *TBox {
	label, _ := gtk.LabelNew(*xtitle)
	gTheme.normalMarkup(label, xtitle)
	ebox, _ := gtk.EventBoxNew()
	ebox.SetMarginEnd(int(ColumnSep))
	// ebox.SetVisibleWindow(true)
	// ebox.SetAboveChild(true)
	ebox.Add(label)

	ebox.Connect("enter-notify-event", func() {
		// fmt.Println("enter notify")
	})
	ebox.Connect("leave-notify-event", func() {
		// fmt.Println("leave notify")
	})

	return &TBox{ebox, label}
}

func (b *TBox) refreshMarkup(r, c int) {
	if Position < len(gData.drows) {
		if r == Position {
			gTheme.cursorMarkup(b.label, &gData.drows[r].xtitles[c])
		} else {
			gTheme.normalMarkup(b.label, &gData.drows[r].xtitles[c])
		}
	}
}

func (b *TBox) remove() {
	b.label.Destroy()
	b.Destroy()
}

// TRow holds an horizontal row of elements
type TRow struct {
	*gtk.Box
	tboxes []*TBox
}

func newTRow(titles []string, n int) *TRow {
	var tboxes []*TBox
	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	hbox.SetName(strconv.Itoa(n))
	hbox.AddEvents(int(gdk.SCROLL_MASK))

	for i := range gData.widths {
		var new *TBox
		if titles == nil {
			new = newTBox(&gData.empties[i])
		} else {
			new = newTBox(&gData.drows[i].xtitles[i])

		}
		hbox.Add(new)
		tboxes = append(tboxes, new)
	}

	hbox.Connect("button_press_event", func(_ *gtk.Box, event *gdk.Event) {
		name, _ := hbox.GetName()
		selectedRow, _ = strconv.Atoi(name)
	})

	hbox.Connect("scroll_event", func(_ *gtk.Box, event *gdk.Event) {
		fmt.Println("hbox scroll event")
	})

	return &TRow{hbox, tboxes}
}

func (r *TRow) fillWithEmpties() {
	for c, tbox := range r.tboxes {
		// tbox.label.updateTitle(gData.empties[i])
		gTheme.normalMarkup(tbox.label, &gData.empties[c])
	}
}

func (r *TRow) remove() {
	for _, tbox := range r.tboxes {
		tbox.remove()
	}
	r.Destroy()
}

func (r *TRow) refreshCursorMarkup(row int) {
	vo := VerticalOffset
	if row >= 0 && row < len(gData.drows)-vo {
		for col, tbox := range r.tboxes {
			gTheme.cursorMarkup(tbox.label, &gData.drows[row+vo].xtitles[col])
		}
	} else {
		for col, tbox := range r.tboxes {
			gTheme.cursorMarkup(tbox.label, &gData.empties[col])
		}
	}
}

func (r *TRow) refreshNormalMarkup(row int) {
	vo := VerticalOffset
	if row >= 0 && row < len(gData.drows)-vo {
		for _, sel := range Selection {
			if sel == row {
				for col, tbox := range r.tboxes {
					gTheme.selectMarkup(tbox.label, &gData.drows[row+vo].xtitles[col])
				}
				return
			}
		}
		for col, tbox := range r.tboxes {
			gTheme.normalMarkup(tbox.label, &gData.drows[row+vo].xtitles[col])
		}
	} else {
		for col, tbox := range r.tboxes {
			gTheme.normalMarkup(tbox.label, &gData.empties[col])
		}
	}
}
