package ui

import (
	"github.com/Hurricanezwf/gopass/log"
	termbox "github.com/nsf/termbox-go"
)

func Open() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	NewUI().Open()
}

type UI struct {
	EditBox *EditBox
	ListBox *ListBox
}

func NewUI() *UI {
	return &UI{
		EditBox: NewEditBox(),
		ListBox: NewListBox(),
	}
}

func (ui *UI) Open() {
	var (
		editBox = ui.EditBox
		listBox = ui.ListBox
	)

	editBox.Open(ui)
	listBox.Open(ui)
	defer listBox.Close()
	termbox.Flush()

	// watch event
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			log.Debug("Receive key %d\n", ev.Key)
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyBackspace2:
				editBox.KeyBackspace2Handler()
			case termbox.KeyBackspace:
				editBox.KeyBackspaceHandler()
			case termbox.KeyArrowUp:
				listBox.KeyArrowUpHandler()
			case termbox.KeyArrowDown:
				listBox.KeyArrowDownHandler()
			case termbox.KeyCtrlQ:
				editBox.KeyCtrlQHandler()
			case termbox.KeySpace:
				editBox.InsertRune(' ')
			case termbox.KeyEnter:
				listBox.KeyEnterHandler()
			default:
				editBox.DefaultHandler(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

//const (
//	boxSpan      = 30   // 编辑框的长度,单位cell
//	boxTopMargin = 0.25 // 编辑框距离顶端的百分比
//)

//var (
//	refCount int
//)
//
//type UI struct{}
//
//func New() *UI {
//	if refCount <= 0 {
//		if err := termbox.Init(); err != nil {
//			return nil
//		}
//		//termbox.SetInputMode(termbox.InputEsc)
//	}
//	refCount++
//	return &UI{}
//}
//
//func (u *UI) Close() {
//	refCount--
//	if refCount <= 0 {
//		termbox.Close()
//	}
//}
//
//func (u *UI) EditBox() *EditBox {
//	return NewEditBox()
//}
//
//func (u *UI) ListBox() *ListBox {
//	return NewListBox()
//}

//func DrawSearchBox(l []string) error {
//	if isInit == false {
//		return errors.New("Not init")
//	}
//
//	color := termbox.ColorDefault
//	termbox.Clear(color, color)
//
//	w, h := termbox.Size()
//
//	// show edit box
//	boxX := (w - boxSpan) / 2
//	boxY := int(boxTopMargin * float64(h))
//	termbox.SetCell(boxX-1, boxY-1, '┌', color, color)
//	termbox.SetCell(boxX-1, boxY, '│', color, color)
//	termbox.SetCell(boxX-1, boxY+1, '└', color, color)
//	termbox.SetCell(boxX+boxSpan, boxY-1, '┐', color, color)
//	termbox.SetCell(boxX+boxSpan, boxY, '│', color, color)
//	termbox.SetCell(boxX+boxSpan, boxY+1, '┘', color, color)
//	fill(boxX, boxY-1, boxSpan, 1, termbox.Cell{Ch: '─'})
//	fill(boxX, boxY+1, boxSpan, 1, termbox.Cell{Ch: '─'})
//
//	tbprint(boxX, boxY, color, color, "→")
//	tbprint(boxX+6, boxY-3, color, color, "Press ESC to quit")
//	termbox.SetCursor(boxX+2, boxY)
//
//	termbox.Flush()
//
//	drawSearchBoxWatcher()
//	return nil
//}
//
//func DrawSearchBoxWatcher() {
//	for {
//		switch ev := termbox.PollEvent(); ev.Type {
//		case termbox.EventKey:
//			switch ev.Key {
//			case termbox.KeyEsc:
//				return
//			default:
//				if ev.Ch != 0 {
//					//TODO: insert edit box
//				}
//			}
//		case termbox.EventError:
//			panic(ev.Err)
//		}
//	}
//}
//
//func fill(x, y, w, h int, cell termbox.Cell) {
//	for ly := 0; ly < h; ly++ {
//		for lx := 0; lx < w; lx++ {
//			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
//		}
//	}
//}
//
//func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
//	for _, c := range msg {
//		termbox.SetCell(x, y, c, fg, bg)
//		x += runewidth.RuneWidth(c)
//	}
//}
