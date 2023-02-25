package tablam

import (
	// "fmt"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gtk"
)

const dma string = "<span background=\"white\"><tt>"
const dmb string = "</tt></span>"
const cma string = "<span foreground=\"white\" background=\"#6666dd\"><tt>"
const cmb string = "</tt></span>"

var dataMarkup = [2]string{dma, dmb}
var cursorMarkup = [2]string{cma, cmb}

var rowSep uint = 4
var columnSep uint = 4
var leftRightMargin = 1

type MyAlign uint8

const (
	AlignLeft MyAlign = iota + 1
	AlignCenter
	AlignRight
)

var data [][]string
var aligns []MyAlign

// TBox represents an element
type TBox struct {
	title    string
	titlex   string
	eventBox *gtk.EventBox
	label    *gtk.Label
	inUse    bool
}

func newTBox(width int) TBox {
	titlex := generateX("", AlignLeft, width)

	label, _ := gtk.LabelNew(titlex)
	label.SetMarkup(dataMarkup[0] + titlex + dataMarkup[1])
	ebox, _ := gtk.EventBoxNew()
	ebox.Add(label)

	return TBox{"", titlex, ebox, label, true}
}

func (b *TBox) Fill(title string, align MyAlign, width int, inUse bool) {
	grow := width - utf8.RuneCountInString(title)
	b.title = title
	b.titlex = generateX(title, align, grow)
	b.label.SetMarkup(dataMarkup[0] + b.titlex + dataMarkup[1])
}

func (e *TBox) update(title string, width int, align MyAlign, inUse bool) {
	e.title = title
	grow := width - utf8.RuneCountInString(e.title)
	e.titlex = generateX(e.title, align, grow)
	e.label.SetMarkup(dataMarkup[0] + e.titlex + dataMarkup[1])
	e.inUse = inUse
}

// TColumn holds a vertical column of elements
type TColumn struct {
	tboxes []TBox
	align  MyAlign
	width  int
}

func NewTColumn(rows, width int) TColumn {
	var tboxes []TBox
	for i := 0; i < rows; i++ {
		tboxes = append(tboxes, newTBox(width))
	}
	return TColumn{tboxes, AlignLeft, width}
}

func (c *TColumn) FillWithText(n, rows int, align MyAlign) {
	var width int

	for i := 0; i < len(data); i++ {
		if utf8.RuneCountInString(data[i][n]) > width {
			width = utf8.RuneCountInString(data[i][n])
		}
	}
	for i := 0; i < rows; i++ {
		c.tboxes[i].Fill(data[i][n], align, width, true)
	}

	c.align = align
	c.width = width
}

func (c *TColumn) CompleteWithWhite(n int) {
	for i := n; i < len(c.tboxes); i++ {
		c.tboxes[i].Fill("", c.align, c.width, false)
	}
}

// Tablam is a gtk grid with an array of TColumn
type Tablam struct {
	*gtk.Grid
	ecols []TColumn
}

func NewTablam(rows, cols, width int) Tablam {
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(2)
	grid.SetColumnSpacing(2)

	var ecols []TColumn

	for i := 0; i < cols; i++ {
		ecol := NewTColumn(rows, width)
		ecols = append(ecols, ecol)
	}

	for i, ecol := range ecols {
		for j, tbox := range ecol.tboxes {
			grid.Attach(tbox.eventBox, i, j, 1, 1)
		}
	}

	return Tablam{grid, ecols}
}

func (t *Tablam) FillWithData(d [][]string, aligns []MyAlign) {
	if len(d) == 0 || len(t.ecols) == 0 {
		return
	}
	data = d

	rows := len(data)
	cols := len(data[0])

	if cols > len(t.ecols) {
		cols = len(t.ecols)
	}
	if rows > len(t.ecols[0].tboxes) {
		rows = len(t.ecols[0].tboxes)
	}

	for i := 0; i < cols; i++ {
		t.ecols[i].FillWithText(i, rows, aligns[i])
	}
	for i := 0; i < len(t.ecols); i++ {
		t.ecols[i].CompleteWithWhite(rows)
	}
}
