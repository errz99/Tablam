package Tablam

import (
	"sort"
	"strconv"
	"strings"
)

func generateX(elem *string, align MyAlign, grow int) string {
	if grow < 0 {
		grow = 0
	}

	if elem == nil {
		return ""
	}

	sep := strings.Repeat(" ", LeftRightMargin)

	switch align {
	case AlignLeft:
		return sep + *elem + strings.Repeat(" ", grow) + sep
	case AlignRight:
		return sep + strings.Repeat(" ", grow) + *elem + sep
	case AlignCenter:
		a := grow / 2
		b := grow / 2
		if grow%2 != 0 {
			b++
		}
		return sep + strings.Repeat(" ", a) + *elem + strings.Repeat(" ", b) + sep
	default:
		return sep + *elem + strings.Repeat(" ", grow) + sep
	}
}

func sortData() {
	n := SortingColumn
	if ReverseSorting {
		sort.SliceStable(gData.drows, func(i, j int) bool {
			return gData.drows[j].titles[n] < gData.drows[i].titles[n]
		})
		ReverseSorting = false
	} else {
		sort.SliceStable(gData.drows, func(i, j int) bool {
			return gData.drows[i].titles[n] < gData.drows[j].titles[n]
		})
		ReverseSorting = true
	}

	if IndexColumn >= 0 {
		indexData()
	}
}

func indexData() {
	for i := range gData.drows {
		gData.drows[i].titles[IndexColumn] = strconv.Itoa(i + 1)
	}
}
