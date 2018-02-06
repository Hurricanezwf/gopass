package ui

import (
	"bytes"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type Notify struct {
	msg string

	duration time.Duration
}

func NewNotify(msg string, duration time.Duration) *Notify {
	return &Notify{
		msg:      msg,
		duration: duration,
	}
}

func (n *Notify) Info() {
	n.show(termbox.AttrBold, termbox.ColorGreen)
}

func (n *Notify) Warn() {
	n.show(termbox.AttrBold, termbox.ColorYellow)
}

func (n *Notify) show(fg, bg termbox.Attribute) {
	w, _ := termbox.Size()

	placeholders := (w - len(n.msg)) / 2

	buf := bytes.NewBuffer(nil)
	for placeholders > 0 {
		placeholders--
		buf.WriteByte(' ')
	}
	buf.WriteString(n.msg)
	placeholders = w - buf.Len()
	for placeholders > 0 {
		placeholders--
		buf.WriteByte(' ')
	}

	TBPrint(0, 0, fg, bg, buf.String())

	time.Sleep(n.duration)
	Fill(0, 0, w, 1, termbox.Cell{Ch: ' '})
}
