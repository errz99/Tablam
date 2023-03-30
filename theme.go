package tablam

import (
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

const fontSizeUnit = 1024
const minFontSize = 4
const initialFontSize = 11

const headerColor = ""
const headerBackground = ""
const normalColor = "#282828"
const normalBackground = "#ffffff"
const cursorColor = "#ffffff"
const cursorBackground = "#6666dd"
const selectColor = "#282828"
const selectBackground = "#d8d8d8"

type ThStyle struct {
	boldA      string
	boldB      string
	color      string
	background string
	chainA     string
	chainB     string
}

func newThStyle(fg, bg, size string, bold bool) ThStyle {
	boldA := ""
	boldB := ""
	color := ""
	background := ""

	if bold {
		boldA = "<b>"
		boldB = "</b>"
	}

	if fg != "" {
		color = " color='" + fg + "'"
	}
	if bg != "" {
		background = " background='" + bg + "'"
	}

	chainA := "<span" + size + background + color + "><tt>" + boldA
	chainB := boldB + "</tt></span>"

	return ThStyle{
		boldA,
		boldB,
		color,
		background,
		chainA,
		chainB,
	}
}

func (ts *ThStyle) setFontBold(bold bool) {
	if bold {
		ts.boldA = "<b>"
		ts.boldB = "</b>"
	} else {
		ts.boldA = ""
		ts.boldB = ""
	}
}

func (ts ThStyle) colors() (string, string) {
	return ts.color, ts.background
}

func (ts *ThStyle) setColors(fg, bg string) {
	if fg != "" {
		ts.color = " color='" + fg + "'"
	} else {
		ts.color = ""
	}
	if bg != "" {
		ts.background = " background='" + bg + "'"
	} else {
		ts.background = ""
	}
}

func (ts *ThStyle) setChains(size string) {
	ts.chainA = "<span" + size + ts.background + ts.color + "><tt>" + ts.boldA
	ts.chainB = ts.boldB + "</tt></span>"
}

type Theme struct {
	fsize    int
	fontSize string

	hStyle ThStyle
	nStyle ThStyle
	cStyle ThStyle
	sStyle ThStyle
}

func newTheme() Theme {
	fontSize := " size='" + strconv.Itoa(initialFontSize*fontSizeUnit) + "'"
	return Theme{
		fsize:    initialFontSize,
		fontSize: fontSize,

		hStyle: newThStyle(headerColor, headerBackground, fontSize, true),
		nStyle: newThStyle(normalColor, normalBackground, fontSize, false),
		cStyle: newThStyle(cursorColor, cursorBackground, fontSize, false),
		sStyle: newThStyle(selectColor, selectBackground, fontSize, false),
	}
}

func (t Theme) themeFontSize() int {
	return t.fsize
}

func (t *Theme) setFontSize(size int) {
	t.fsize = size
	if size > 0 {
		t.fontSize = " size='" + strconv.Itoa(size*fontSizeUnit) + "'"
	} else {
		t.fontSize = ""
	}
	t.hStyle.setChains(t.fontSize)
	t.nStyle.setChains(t.fontSize)
	t.cStyle.setChains(t.fontSize)
	t.sStyle.setChains(t.fontSize)
}

func (t *Theme) incFontSize() {
	t.fsize++
	t.fontSize = " size='" + strconv.Itoa(t.fsize*fontSizeUnit) + "'"
	t.hStyle.setChains(t.fontSize)
	t.nStyle.setChains(t.fontSize)
	t.cStyle.setChains(t.fontSize)
	t.sStyle.setChains(t.fontSize)
}

func (t *Theme) decFontSize() bool {
	if t.fsize > minFontSize {
		t.fsize--
		t.fontSize = " size='" + strconv.Itoa(t.fsize*fontSizeUnit) + "'"
		t.hStyle.setChains(t.fontSize)
		t.nStyle.setChains(t.fontSize)
		t.cStyle.setChains(t.fontSize)
		t.sStyle.setChains(t.fontSize)
		return true
	}
	return false
}

func (t Theme) fontBold(c byte) bool {
	var isBold bool
	switch c {
	case 'h':
		if t.hStyle.boldA != "" {
			isBold = true
		}
	case 'n':
		if t.nStyle.boldA != "" {
			isBold = true
		}
	case 'c':
		if t.cStyle.boldA != "" {
			isBold = true
		}
	case 's':
		if t.sStyle.boldA != "" {
			isBold = true
		}
	default:
	}
	return isBold
}

func (t *Theme) setFontBold(c byte, bold bool) {
	switch c {
	case 'h':
		t.hStyle.setFontBold(bold)
		t.hStyle.setChains(t.fontSize)
	case 'n':
		t.nStyle.setFontBold(bold)
		t.nStyle.setChains(t.fontSize)
	case 'c':
		t.cStyle.setFontBold(bold)
		t.cStyle.setChains(t.fontSize)
	case 's':
		t.sStyle.setFontBold(bold)
		t.sStyle.setChains(t.fontSize)
	default:
	}
}

func (t Theme) colors(c byte) (string, string) {
	switch c {
	case 'h':
		return t.hStyle.colors()
	case 'n':
		return t.nStyle.colors()
	case 'c':
		return t.cStyle.colors()
	case 's':
		return t.sStyle.colors()
	default:
	}
	return "", ""
}

func (t *Theme) setColors(c byte, fg, bg string) {
	switch c {
	case 'h':
		t.hStyle.setColors(fg, bg)
		t.hStyle.setChains(t.fontSize)
	case 'n':
		t.nStyle.setColors(fg, bg)
		t.nStyle.setChains(t.fontSize)
	case 'c':
		t.cStyle.setColors(fg, bg)
		t.cStyle.setChains(t.fontSize)
	case 's':
		t.sStyle.setColors(fg, bg)
		t.sStyle.setChains(t.fontSize)
	default:
	}
}

func (t Theme) headMarkup(label *gtk.Label, text *string) {
	label.SetMarkup(t.hStyle.chainA + *text + t.hStyle.chainB)
}

func (t Theme) normalMarkup(label *gtk.Label, text *string) {
	label.SetMarkup(t.nStyle.chainA + *text + t.nStyle.chainB)
}

func (t Theme) cursorMarkup(label *gtk.Label, text *string) {
	label.SetMarkup(t.cStyle.chainA + *text + t.cStyle.chainB)
}

func (t Theme) selectMarkup(label *gtk.Label, text *string) {
	label.SetMarkup(t.sStyle.chainA + *text + t.sStyle.chainB)
}
