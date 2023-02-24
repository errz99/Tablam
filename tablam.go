package tablam

import (
	"fmt"
	"strings"
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

type Ebox struct {
	title    string
	titlex   string
	eventBox *gtk.EventBox
	label    *gtk.Label
	inUse    bool
}

func newEbox(title string, align MyAlign, width int, inUse bool) Ebox {
	grow := width - utf8.RuneCountInString(title)
	titlex := generateX(title, align, grow)

	label, _ := gtk.LabelNew(titlex)
	label.SetMarkup(dataMarkup[0] + titlex + dataMarkup[1])
	ebox, _ := gtk.EventBoxNew()
	ebox.Add(label)

	return Ebox{title, titlex, ebox, label, inUse}
}

func (e *Ebox) Update(title string, width int, inUse bool) {
	e.title = title
	grow := width - utf8.RuneCountInString(e.title)
	e.titlex = generateX(e.title, e.align, grow)
	e.label.SetMarkup(dataMarkup[0] + e.titlex + dataMarkup[1])
	e.inUse = inUse
}

func generateX(elem string, align MyAlign, grow int) string {
	if grow < 0 {
		grow = 0
	}

	sep := strings.Repeat(" ", leftRightMargin)

	if align == AlignLeft {
		return sep + elem + strings.Repeat(" ", grow) + sep

	} else if align == AlignRight {
		return sep + strings.Repeat(" ", grow) + elem + sep

	} else if align == AlignCenter {
		a := grow / 2
		b := grow / 2
		if grow%2 != 0 {
			b++
		}
		return sep + strings.Repeat(" ", a) + elem + strings.Repeat(" ", b) + sep
	} else {
		return sep + elem + strings.Repeat(" ", grow) + sep
	}
}

type EColumn struct {
	eboxes []Ebox
	align  MyAlign
	width  int
}

func NewEColumn(titles []string, align MyAlign) EColumn {
	var width int
	for _, title := range titles {
		if utf8.RuneCountInString(title) > width {
			title = utf8.RuneCountInString(title)
		}
	}

	var eboxes []Ebox
	for _, title := range titles {
		eboxes = append(eboxes, newEbox(title, align, width, true))
	}

	return EColumn{eboxes, align, width}
}

func Hello(name string) {
	fmt.Println("Hello", name, "from Tablam")
}
