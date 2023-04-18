package main

// gofmt -l -w .

import (
	"fmt"
	"log"
	"os"

	tb "github.com/errz99/Tablam"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	ModifierShift   bool
	ModifierAlt     bool
	ModifierControl bool

	head = []string{
		"Nr", "First N", "Second N", "Age", "City",
	}

	haligns = []tb.MyAlign{
		tb.AlignRight,
		tb.AlignLeft,
		tb.AlignLeft,
		tb.AlignLeft,
		tb.AlignCenter,
	}

	data = [][]string{
		[]string{"1", "Jorma", "Kaukonen", "70", "San Francisco"},
		[]string{"2", "John", "Smith", "8", "Sausalito"},
		[]string{"3", "Kim", "Jon-Un", "45", "Pyonyang"},
		[]string{"4", "Ian", "Anderson", "72", "Glasgow"},
		[]string{"5", "Patty", "Smith", "77", "New York"},
		[]string{"6", "Bob", "Dylan", "78", "Duluth"},
	}

	aligns = []tb.MyAlign{
		tb.AlignRight,
		tb.AlignRight,
		tb.AlignLeft,
		tb.AlignRight,
		tb.AlignLeft,
	}
)

func main() {
	initGui()
}

func initGui() {
	application, err := gtk.ApplicationNew("com.errz99.tablam_t2", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Gtk3+ Initialization failed.", err)
	}

	application.Connect("activate", func() {
		mainWindow(application)
	})

	os.Exit(application.Run(nil))
}

func mainWindow(app *gtk.Application) {
	window, _ := gtk.ApplicationWindowNew(app)
	window.SetTitle("Tablam_t2")
	window.SetDefaultSize(600, 400)
	window.Move(200, 200)
	window.AddEvents(int(gdk.SCROLL_MASK))

	window.Connect("delete_event", func() {
		window.Close()
	})

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	vbox.SetHAlign(gtk.ALIGN_CENTER)
	vbox.SetBorderWidth(8)
	window.Add(vbox)

	label, _ := gtk.LabelNew("Some Example")
	label.SetMarginBottom(8)
	vbox.Add(label)

	// Tablam
	tb.RowSep = 2
	tb.ColumnSep = 2
	tb.LeftRightMargin = 1
	// tb.SortingColumn = 1
	tb.IndexColumn = 0
	// tb.PageMode = true

	header, err := tb.NewHeader(head, haligns)
	if err == nil {
		header.SetMarginBottom(3)
		// header.SetColors("", "#99ee99")
		vbox.Add(header)
	} else {
		log.Println(err)
	}

	// scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	// scroll.SetVExpand(true)
	// scroll.SetHExpand(true)
	// vbox.Add(scroll)

	mytab, err := tb.NewTablam(header, len(data), len(aligns), 2)
	if err == nil {
		vbox.Add(mytab)
		mytab.FillWithData(data, aligns)

		if header != nil {
			header.Connect("button_press_event", func(_ *gtk.Box, event *gdk.Event) {
				active := tb.ActiveHeaderCell
				if active > 0 && active != 3 {
					mytab.SortData(active)
				}
			})
		}

		mytab.Connect("scroll_event", func(_ *gtk.Grid, event *gdk.Event) {
			fmt.Println("mytab scroll event")
		})

	} else {
		log.Println(err)
		os.Exit(1)
	}

	window.Connect("button_press_event", func(_ *gtk.ApplicationWindow, event *gdk.Event) {
		if ModifierControl {
			mytab.SelectARow()
		} else {
			mytab.UpdateCursor(true)
		}
	})

	window.Connect("scroll_event", func(_ *gtk.ApplicationWindow, event *gdk.Event) {
		eventScroll := gdk.EventScrollNewFromEvent(event)
		switch eventScroll.Direction() {
		case gdk.SCROLL_UP:
			// if tb.VerticalOffset > 0 && tb.VerticalOffset < mytab.DataRows() - mytab.TotalRows() {
			if tb.VerticalOffset > 0 {
				tb.VerticalOffset--
				// fmt.Println("offset:", tb.VerticalOffset)
				mytab.UpdateWithOffset(-1)
			}
		case gdk.SCROLL_DOWN:
			if tb.VerticalOffset < mytab.DataRows() && mytab.DataRows() > mytab.TotalRows() {
				tb.VerticalOffset++
				// fmt.Println("offset:", tb.VerticalOffset)
				mytab.UpdateWithOffset(1)
			}
		default:
		}
	})

	window.Connect("key_press_event", func(_ *gtk.ApplicationWindow, event *gdk.Event) {
		evKey := gdk.EventKeyNewFromEvent(event)
		switch evKey.KeyVal() {
		case gdk.KEY_Alt_L, gdk.KEY_Alt_R:
			ModifierAlt = true
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			ModifierControl = true
		case gdk.KEY_Shift_L, gdk.KEY_Shift_R:
			ModifierShift = true

		case gdk.KEY_Escape:
			if tb.Position >= 0 {
				mytab.RemoveCursor()
			} else {
				mytab.ClearSelected()
			}

		case gdk.KEY_Return:
			if ModifierControl {
				mytab.SelectCursorRow()
			} else {
				rd := mytab.ActiveRowData()
				if rd != nil {
					fmt.Println(rd)
				}
			}

		case gdk.KEY_Page_Up:
			mytab.DecFontSize()

		case gdk.KEY_Page_Down:
			mytab.IncFontSize()

		case gdk.KEY_Up:
			mytab.CursorUp()
		case gdk.KEY_Down:
			mytab.CursorDown()

		case gdk.KEY_F12:
			mytab.UpdateCell(0, 4, "Santos Franciscos")

		case gdk.KEY_F11:
			mytab.UpdateCell(0, 4, "Frisco")

		case gdk.KEY_F10:
			mytab.UpdateCell(4, 4, "News Yorks")

		case gdk.KEY_F9:
			mytab.UpdateCell(4, 4, "York")

		case gdk.KEY_F8:
			mytab.UpdateRow(0, []string{"1", "Miguel", "Arias", "65", "Málaga"})

		case gdk.KEY_F7:
			mytab.UpdateRow(0, data[1])

		case gdk.KEY_F6:
			if mytab.HeaderFontBold() {
				mytab.SetHeaderFontBold(false)
			} else {
				mytab.SetHeaderFontBold(true)
			}

		case gdk.KEY_F5:
			if mytab.NormalFontBold() {
				mytab.SetNormalFontBold(false)
			} else {
				mytab.SetNormalFontBold(true)
			}

		case gdk.KEY_F4:
			mytab.SetCursorFontBold(true)

		case gdk.KEY_Delete:
			mytab.RemoveARow(tb.Position)

		case gdk.KEY_Insert:
			mytab.InsertARow(tb.Position,
				[]string{"7", "Jerry", "Garcia", "80", "Santos Franciscos"})

		case gdk.KEY_A:
			mytab.AddDataRow([]string{"9", "Jerry", "Garcia", "80", "Santos Franciscos"})

		case gdk.KEY_C:
			mytab.SetCursorColors("", "#99ee99")

		case gdk.KEY_N:
			mytab.SetNormalColors("#aabbcc", "#554433")

		case gdk.KEY_F3:
			mytab.SortData(4)

		case gdk.KEY_F2:
			// mytab.SortData(2)
			mytab.ShowOrHideColumn(4)

		case gdk.KEY_F1:
			// mytab.SortData(1)
			mytab.ShowOrHideColumn(1)

		case gdk.KEY_T:
			mytab.SetFont("")

		case gdk.KEY_I:
			mytab.SetFont("Iosevka Term")

		default:
		}
	})

	window.Connect("key_release_event", func(_ *gtk.ApplicationWindow, event *gdk.Event) {
		evKey := gdk.EventKeyNewFromEvent(event)
		switch evKey.KeyVal() {
		case gdk.KEY_Alt_L, gdk.KEY_Alt_R:
			ModifierAlt = false
		case gdk.KEY_Control_L, gdk.KEY_Control_R:
			ModifierControl = false
		case gdk.KEY_Shift_L, gdk.KEY_Shift_R:
			ModifierShift = false
		default:
		}
	})

	window.ShowAll()
}
