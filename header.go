package Tablam

import (
	"errors"
	"strconv"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// // Index is an optional vertical column
// // on the left containing index numbers
// type Index struct {
// 	titles     []string
// 	titlesx    []string
// 	eventBoxes []*gtk.EventBox
// 	labels     []*gtk.Label
// 	align      MyAlign
// 	width      int
// 	hidden     bool
// 	startPoint int
// }

// func newIndex(amnt, sp, w int, align MyAlign) Index {
// 	var titles []string
// 	for i := 0; i < amnt; i++ {
// 		titles = append(titles, strconv.Itoa(i+sp))
// 	}
// 	return Index{
// 		titles:     titles,
// 		startPoint: sp,
// 		width:      w,
// 		align:      align,
// 	}
// }

// THeader is an optional first row for the Tablam grid
type THeader struct {
	*gtk.Box
	titles     []string
	titlexs    []string
	eventBoxes []*gtk.EventBox
	labels     []*gtk.Label
	aligns     []MyAlign
	widths     []int
	hidden     []bool
}

func NewHeader(titles []string, aligns []MyAlign) (*THeader, error) {
	if len(titles) != len(aligns) {
		return nil, errors.New("New Header: Different number of titles and aligns.")
	}
	gData = newTData(len(titles), 0)
	gTheme = newTheme()
	gTheme.setFontBold('c', true)

	hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	var titlexs []string
	var eboxes []*gtk.EventBox
	var labels []*gtk.Label
	widths := make([]int, len(titles))
	hidden := make([]bool, len(titles))

	for i, title := range titles {
		widths[i] = utf8.RuneCountInString(title)
		titlex := generateX(&title, aligns[i], 0)
		label, _ := gtk.LabelNew(titlex)
		gTheme.headMarkup(label, &titlex)
		ebox, _ := gtk.EventBoxNew()
		ebox.SetMarginEnd(int(ColumnSep))
		ebox.SetName(strconv.Itoa(i))
		ebox.Add(label)

		ebox.Connect("button_press_event", func(_ *gtk.EventBox, event *gdk.Event) {
			name, _ := ebox.GetName()
			active, _ := strconv.Atoi(name)
			ActiveHeaderCell = active
		})

		titlexs = append(titlexs, titlex)
		eboxes = append(eboxes, ebox)
		labels = append(labels, label)

		hbox.Add(ebox)
	}
	copy(gData.widths, widths)

	return &THeader{
		hbox,
		titles,
		titlexs,
		eboxes,
		labels,
		aligns,
		widths,
		hidden,
	}, nil
}

func (h *THeader) refreshMarkups() {
	for i, label := range h.labels {
		gTheme.headMarkup(label, &h.titlexs[i])
	}
}

func (h *THeader) widthsRecalc() {
	for i, label := range h.labels {
		if h.widths[i] != gData.widths[i] {
			grow := gData.widths[i] - utf8.RuneCountInString(h.titles[i])
			h.titlexs[i] = generateX(&h.titles[i], h.aligns[i], grow)
			gTheme.headMarkup(label, &h.titlexs[i])
		}
	}
	h.ShowAll()
}

func (h *THeader) SetColors(fg, bg string) {
	gTheme.setColors('h', fg, bg)
	h.refreshMarkups()
}
