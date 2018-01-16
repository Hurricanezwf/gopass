package ui

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

//func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
//	for _, c := range msg {
//		termbox.SetCell(x, y, c, fg, bg)
//		x += runewidth.RuneWidth(c)
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
//func rune_advance_len(r rune, pos int) int {
//	if r == '\t' {
//		return tabstop_length - pos%tabstop_length
//	}
//	return runewidth.RuneWidth(r)
//}
//
//func voffset_coffset(text []byte, boffset int) (voffset, coffset int) {
//	text = text[:boffset]
//	for len(text) > 0 {
//		r, size := utf8.DecodeRune(text)
//		text = text[size:]
//		coffset += 1
//		voffset += rune_advance_len(r, voffset)
//	}
//	return
//}
//
//func byte_slice_remove(text []byte, from, to int) []byte {
//	size := to - from
//	copy(text[from:], text[to:])
//	text = text[:len(text)-size]
//	return text
//}
//
//func byte_slice_insert(text []byte, offset int, what []byte) []byte {
//	n := len(text) + len(what)
//	text = byte_slice_grow(text, n)
//	text = text[:n]
//	copy(text[offset+len(what):], text[offset:])
//	copy(text[offset:], what)
//	return text
//}

type EditBoxConfig struct {
	// 编辑框距离页面顶端的距离
	boxTopMargin int

	// 编辑框距离页面左边的距离
	boxLeftMargin int

	// 编辑框长度, 默认30px
	boxSpan int

	// 编辑框颜色
	boxColor termbox.Attribute
}

func DefaultEditBoxConfig() *EditBoxConfig {
	eb := &EditBoxConfig{
		boxTopMargin:  0,
		boxLeftMargin: 0,
		boxSpan:       50,
		boxColor:      termbox.ColorDefault,
	}

	w, h := termbox.Size()
	eb.boxTopMargin = int(0.3 * float64(h))
	eb.boxLeftMargin = (w - eb.boxSpan) / 2
	return eb
}

type EditBox struct {
	// 如果conf为nil, 将使用默认配置
	Conf *EditBoxConfig

	// 编辑框内容
	text []byte

	boxX int
	boxY int
}

func NewEditBox() *EditBox {
	return &EditBox{
		Conf: DefaultEditBoxConfig(),
		text: make([]byte, 0),
	}
}

func (eb *EditBox) Show() {
	var (
		color   = eb.Conf.boxColor
		boxX    = eb.Conf.boxLeftMargin
		boxY    = eb.Conf.boxTopMargin
		boxSpan = eb.Conf.boxSpan
	)

	// show edit box
	termbox.SetCell(boxX-1, boxY-1, '┌', color, color)
	termbox.SetCell(boxX-1, boxY, '│', color, color)
	termbox.SetCell(boxX-1, boxY+1, '└', color, color)
	termbox.SetCell(boxX+boxSpan, boxY-1, '┐', color, color)
	termbox.SetCell(boxX+boxSpan, boxY, '│', color, color)
	termbox.SetCell(boxX+boxSpan, boxY+1, '┘', color, color)
	Fill(boxX, boxY-1, boxSpan, 1, termbox.Cell{Ch: '─'})
	Fill(boxX, boxY+1, boxSpan, 1, termbox.Cell{Ch: '─'})

	TBPrint(boxX, boxY, color, color, "→")
	TBPrint(boxX+16, boxY-3, color, color, "Press ESC to quit")
	termbox.SetCursor(boxX+2, boxY)

	termbox.Flush()

	// cache
	eb.boxX = boxX
	eb.boxY = boxY
}

//func (eb *EditBox) watch() {
//	for {
//		switch ev := termbox.PollEvent(); ev.Type {
//		case termbox.EventKey:
//			log.Debug("EditBox receive key %d\n", ev.Key)
//			switch ev.Key {
//			case termbox.KeyEsc:
//				return
//			case termbox.KeyBackspace2:
//				eb.KeyBackspace2Handler()
//			case termbox.KeyBackspace:
//				eb.KeyBackspaceHandler()
//			case termbox.KeyCtrlQ:
//				eb.KeyCtrlQHandler()
//			case termbox.KeySpace:
//				eb.insertRune(' ')
//			default:
//				if ev.Ch != 0 {
//					eb.insertRune(ev.Ch)
//				}
//			}
//		case termbox.EventError:
//			panic(ev.Err)
//		}
//	}
//}

func (eb *EditBox) InsertRune(r rune) {
	b := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(b, r)
	eb.text = ByteSliceInsert(eb.text, len(eb.text), b)
	eb.flush()
}

func (eb *EditBox) DeleteRune() {
	if len(eb.text) > 0 {
		eb.text = eb.text[:len(eb.text)-1]
	}
	eb.flush()
}

func (eb *EditBox) flush() {
	var (
		txt   []byte
		boxX  = eb.boxX
		boxY  = eb.boxY
		color = eb.Conf.boxColor
		space = 2
		span  = eb.Conf.boxSpan - space*2
	)

	if len(eb.text) <= span {
		txt = eb.text
	} else {
		txt = eb.text[len(eb.text)-span:]
	}
	eb.clear()
	TBPrint(boxX+space, boxY, color, color, string(txt))
	termbox.SetCursor(boxX+space+len(txt), boxY)
	termbox.Flush()
}

func (eb *EditBox) clear() {
	var (
		boxX = eb.boxX
		boxY = eb.boxY
		span = eb.Conf.boxSpan - 1
	)
	Fill(boxX+1, boxY, span, 1, termbox.Cell{Ch: ' '})
}

func (eb *EditBox) KeyBackspaceHandler() error {
	return eb.KeyBackspace2Handler()
}

func (eb *EditBox) KeyBackspace2Handler() error {
	eb.DeleteRune()
	return nil
}

func (eb *EditBox) KeyCtrlQHandler() error {
	eb.text = eb.text[:0]
	eb.flush()
	return nil
}

//const preferred_horizontal_threshold = 5
//const tabstop_length = 8
//
//type EditBox struct {
//	text           []byte
//	line_voffset   int
//	cursor_boffset int // cursor offset in bytes
//	cursor_voffset int // visual cursor offset in termbox cells
//	cursor_coffset int // cursor offset in unicode code points
//}
//
//// Draws the EditBox in the given location, 'h' is not used at the moment
//func (eb *EditBox) Draw(x, y, w, h int) {
//	eb.AdjustVOffset(w)
//
//	const coldef = termbox.ColorDefault
//	fill(x, y, w, h, termbox.Cell{Ch: ' '})
//
//	t := eb.text
//	lx := 0
//	tabstop := 0
//	for {
//		rx := lx - eb.line_voffset
//		if len(t) == 0 {
//			break
//		}
//
//		if lx == tabstop {
//			tabstop += tabstop_length
//		}
//
//		if rx >= w {
//			termbox.SetCell(x+w-1, y, '→',
//				coldef, coldef)
//			break
//		}
//
//		r, size := utf8.DecodeRune(t)
//		if r == '\t' {
//			for ; lx < tabstop; lx++ {
//				rx = lx - eb.line_voffset
//				if rx >= w {
//					goto next
//				}
//
//				if rx >= 0 {
//					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
//				}
//			}
//		} else {
//			if rx >= 0 {
//				termbox.SetCell(x+rx, y, r, coldef, coldef)
//			}
//			lx += runewidth.RuneWidth(r)
//		}
//	next:
//		t = t[size:]
//	}
//
//	if eb.line_voffset != 0 {
//		termbox.SetCell(x, y, '←', coldef, coldef)
//	}
//}
//
//// Adjusts line visual offset to a proper value depending on width
//func (eb *EditBox) AdjustVOffset(width int) {
//	ht := preferred_horizontal_threshold
//	max_h_threshold := (width - 1) / 2
//	if ht > max_h_threshold {
//		ht = max_h_threshold
//	}
//
//	threshold := width - 1
//	if eb.line_voffset != 0 {
//		threshold = width - ht
//	}
//	if eb.cursor_voffset-eb.line_voffset >= threshold {
//		eb.line_voffset = eb.cursor_voffset + (ht - width + 1)
//	}
//
//	if eb.line_voffset != 0 && eb.cursor_voffset-eb.line_voffset < ht {
//		eb.line_voffset = eb.cursor_voffset - ht
//		if eb.line_voffset < 0 {
//			eb.line_voffset = 0
//		}
//	}
//}
//
//func (eb *EditBox) MoveCursorTo(boffset int) {
//	eb.cursor_boffset = boffset
//	eb.cursor_voffset, eb.cursor_coffset = voffset_coffset(eb.text, boffset)
//}
//
//func (eb *EditBox) RuneUnderCursor() (rune, int) {
//	return utf8.DecodeRune(eb.text[eb.cursor_boffset:])
//}
//
//func (eb *EditBox) RuneBeforeCursor() (rune, int) {
//	return utf8.DecodeLastRune(eb.text[:eb.cursor_boffset])
//}
//
//func (eb *EditBox) MoveCursorOneRuneBackward() {
//	if eb.cursor_boffset == 0 {
//		return
//	}
//	_, size := eb.RuneBeforeCursor()
//	eb.MoveCursorTo(eb.cursor_boffset - size)
//}
//
//func (eb *EditBox) MoveCursorOneRuneForward() {
//	if eb.cursor_boffset == len(eb.text) {
//		return
//	}
//	_, size := eb.RuneUnderCursor()
//	eb.MoveCursorTo(eb.cursor_boffset + size)
//}
//
//func (eb *EditBox) MoveCursorToBeginningOfTheLine() {
//	eb.MoveCursorTo(0)
//}
//
//func (eb *EditBox) MoveCursorToEndOfTheLine() {
//	eb.MoveCursorTo(len(eb.text))
//}
//
//func (eb *EditBox) DeleteRuneBackward() {
//	if eb.cursor_boffset == 0 {
//		return
//	}
//
//	eb.MoveCursorOneRuneBackward()
//	_, size := eb.RuneUnderCursor()
//	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
//}
//
//func (eb *EditBox) DeleteRuneForward() {
//	if eb.cursor_boffset == len(eb.text) {
//		return
//	}
//	_, size := eb.RuneUnderCursor()
//	eb.text = byte_slice_remove(eb.text, eb.cursor_boffset, eb.cursor_boffset+size)
//}
//
//func (eb *EditBox) DeleteTheRestOfTheLine() {
//	eb.text = eb.text[:eb.cursor_boffset]
//}
//
//func (eb *EditBox) InsertRune(r rune) {
//	var buf [utf8.UTFMax]byte
//	n := utf8.EncodeRune(buf[:], r)
//	eb.text = byte_slice_insert(eb.text, eb.cursor_boffset, buf[:n])
//	eb.MoveCursorOneRuneForward()
//}
//
//// Please, keep in mind that cursor depends on the value of line_voffset, which
//// is being set on Draw() call, so.. call this method after Draw() one.
//func (eb *EditBox) CursorX() int {
//	return eb.cursor_voffset - eb.line_voffset
//}

/*
var edit_box EditBox

const edit_box_width = 30

func redraw_all() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	midy := h / 2
	midx := (w - edit_box_width) / 2

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy+1, '┘', coldef, coldef)
	fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})

	edit_box.Draw(midx, midy, edit_box_width, 1)
	termbox.SetCursor(midx+edit_box.CursorX(), midy)

	tbprint(midx+6, midy+3, coldef, coldef, "Press ESC to quit")
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	redraw_all()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				edit_box.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				edit_box.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				edit_box.DeleteRuneBackward()
			case termbox.KeyDelete, termbox.KeyCtrlD:
				edit_box.DeleteRuneForward()
			case termbox.KeyTab:
				edit_box.InsertRune('\t')
			case termbox.KeySpace:
				edit_box.InsertRune(' ')
			case termbox.KeyCtrlK:
				edit_box.DeleteTheRestOfTheLine()
			case termbox.KeyHome, termbox.KeyCtrlA:
				edit_box.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				edit_box.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					edit_box.InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all()
	}
}
*/
