package ui

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
)

const (
	boxSpan      = 30   // 编辑框的长度,单位cell
	boxTopMargin = 0.25 // 编辑框距离顶端的百分比
)

var isInit = false

func Init() error {
	var err error

	if err = termbox.Init(); err != nil {
		return err
	}

	isInit = true
	termbox.SetInputMode(termbox.InputEsc)
	return nil
}

func Close() {
	isInit = false
	termbox.Close()
}

func DrawAll(l []string) error {
	if isInit == false {
		return errors.New("Not init")
	}

	color := termbox.ColorDefault
	termbox.Clear(color, color)

	w, h := termbox.Size()

	// show edit box
	boxX := (w - boxSpan) / 2
	boxY := int(boxTopMargin * float64(h))
	termbox.SetCell(boxX-1, boxY-1, '┌', color, color)
	termbox.SetCell(boxX-1, boxY, '│', color, color)
	termbox.SetCell(boxX-1, boxY+1, '└', color, color)
	termbox.SetCell(boxX+boxSpan, boxY-1, '┐', color, color)
	termbox.SetCell(boxX+boxSpan, boxY, '│', color, color)
	termbox.SetCell(boxX+boxSpan, boxY+1, '┘', color, color)

	termbox.Flush()

	drawAllWatch()
	return nil
}

func drawAllWatch() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			default:
				if ev.Ch != 0 {
					//TODO: insert edit box
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
