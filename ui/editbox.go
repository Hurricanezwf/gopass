package ui

import (
	"bytes"
	"container/list"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

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

	ui *UI

	// 编辑框内容
	text     []byte
	runeSize *list.List

	boxX int
	boxY int
}

func NewEditBox() *EditBox {
	return &EditBox{
		Conf:     DefaultEditBoxConfig(),
		text:     make([]byte, 0),
		runeSize: list.New(),
	}
}

func (eb *EditBox) Open(ui *UI) {
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
	eb.ui = ui
}

func (eb *EditBox) Value() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(eb.text)
	return buf.Bytes()
}

func (eb *EditBox) InsertRune(r rune) {
	b := make([]byte, utf8.RuneLen(r))
	n := utf8.EncodeRune(b, r)
	eb.text = ByteSliceInsert(eb.text, len(eb.text), b[:n])
	eb.runeSize.PushBack(n)
	eb.flush()
}

func (eb *EditBox) DeleteRune() {
	if len(eb.text) > 0 {
		element := eb.runeSize.Back()
		rmSize := element.Value.(int)
		eb.text = eb.text[:len(eb.text)-rmSize]
		eb.runeSize.Remove(element)
	}
	eb.flush()
}

func (eb *EditBox) flush() {
	var (
		txt              []byte
		boxX             = eb.boxX
		boxY             = eb.boxY
		color            = eb.Conf.boxColor
		space            = 2
		span             = eb.Conf.boxSpan - space*2
		cursorPos        = boxX + space
		cursorRightLimit = boxX + span
	)

	displayLen := 0
	for back := eb.runeSize.Back(); back != nil; back = back.Prev() {
		l := back.Value.(int)
		if cursorPos > cursorRightLimit {
			break
		}
		if l > 1 {
			// 只有非英文字符才大于1字节，中文显示用两个字节
			cursorPos += 2
		} else {
			cursorPos += 1
		}
		displayLen += l
	}
	txt = eb.text[len(eb.text)-displayLen:]
	eb.clear()
	TBPrint(boxX+space, boxY, color, color, string(txt))
	termbox.SetCursor(cursorPos, boxY)
	termbox.Flush()
}

// clear会清空编辑框的显示内容
func (eb *EditBox) clear() {
	var (
		boxX = eb.boxX
		boxY = eb.boxY
		span = eb.Conf.boxSpan - 1
	)
	Fill(boxX+1, boxY, span, 1, termbox.Cell{Ch: ' '})
}

func (eb *EditBox) Drop() error {
	eb.text = eb.text[:0]
	eb.runeSize = list.New()
	eb.flush()
	return nil
}
