package Tablam

import "strings"

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
