package ui

import (
	"errors"
	"fmt"

	"github.com/Hurricanezwf/gopass/log"
	termbox "github.com/nsf/termbox-go"
)

func Open() error {
	if err := termbox.Init(); err != nil {
		return fmt.Errorf("Init UI failed, %v", err)
	}
	NewUI().Open()
	return nil
}

func Close() error {
	termbox.Close()
	return nil
}

type KeyHandler func(ui *UI, ch rune) error

type UI struct {
	EditBox *EditBox
	ListBox *ListBox

	// 按键捕获
	handlers map[termbox.Key]KeyHandler
}

func NewUI() *UI {
	return &UI{
		EditBox: NewEditBox(),
		ListBox: NewListBox(),
		handlers: map[termbox.Key]KeyHandler{
			termbox.KeyEsc:        KeyEscHandler,
			termbox.KeyBackspace2: KeyBackspace2Handler,
			termbox.KeyBackspace:  KeyBackspaceHandler,
			termbox.KeyArrowUp:    KeyArrowUpHandler,
			termbox.KeyArrowDown:  KeyArrowDownHandler,
			termbox.KeyCtrlQ:      KeyCtrlQHandler,
			termbox.KeySpace:      KeySpaceHandler,
			termbox.KeyEnter:      KeyEnterHandler,
		},
	}
}

func (ui *UI) Open() {
	ui.EditBox.Open(ui)
	ui.ListBox.Open(ui)
	termbox.Flush()

	defer ui.ListBox.Close()

	// watch event
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			log.Debug("Receive key %d\n", ev.Key)
			h, ok := ui.handlers[ev.Key]
			if !ok {
				KeyDefaultHandler(ui, ev.Ch)
				continue
			}
			if h != nil {
				if err := h(ui, ev.Ch); err != nil {
					return
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func KeyEscHandler(ui *UI, ch rune) error {
	return errors.New("exit actively")
}

func KeyBackspace2Handler(ui *UI, ch rune) error {
	ui.EditBox.DeleteRune()
	ui.ListBox.NotifyMatch()
	return nil
}

func KeyBackspaceHandler(ui *UI, ch rune) error {
	return KeyBackspace2Handler(ui, ch)
}

func KeyArrowUpHandler(ui *UI, ch rune) error {
	return ui.ListBox.Prev()
}

func KeyArrowDownHandler(ui *UI, ch rune) error {
	return ui.ListBox.Next()
}

func KeyCtrlQHandler(ui *UI, ch rune) error {
	ui.EditBox.Drop()
	ui.ListBox.NotifyMatch()
	return nil
}

func KeySpaceHandler(ui *UI, ch rune) error {
	ui.EditBox.InsertRune(' ')
	return nil
}

func KeyEnterHandler(ui *UI, ch rune) error {
	return ui.ListBox.CopySel()
}

func KeyDefaultHandler(ui *UI, ch rune) error {
	if ch == 0 {
		return fmt.Errorf("Empty rune")
	}
	ui.EditBox.InsertRune(ch)
	ui.ListBox.NotifyMatch()
	return nil
}
