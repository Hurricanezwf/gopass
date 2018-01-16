package ui

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

type ListBoxConfig struct {
	// 列表距离页面顶端的距离
	boxTopMargin int

	// 列表距离页面左侧的距离
	boxLeftMargin int

	// 列表宽度, 默认30px
	boxSpanX int
	// 列表高度, 默认60px
	boxSpanY int

	// 列表颜色
	boxColor termbox.Attribute
}

func DefaultListBoxConfig() *ListBoxConfig {
	lb := &ListBoxConfig{
		boxTopMargin:  0,
		boxLeftMargin: 0,
		boxSpanX:      50,
		boxSpanY:      20,
		boxColor:      termbox.ColorCyan,
	}

	w, h := termbox.Size()
	lb.boxTopMargin = int(0.3*float64(h)) + 1
	lb.boxLeftMargin = (w - lb.boxSpanX) / 2
	return lb
}

type ListBox struct {
	// 如果conf为nil, 将使用默认配置
	Conf *ListBoxConfig

	boxX int
	boxY int
}

func NewListBox() *ListBox {
	return &ListBox{
		Conf: DefaultListBoxConfig(),
	}
}

func (lb *ListBox) Show() {
	var (
		color    = lb.Conf.boxColor
		boxX     = lb.Conf.boxLeftMargin
		boxY     = lb.Conf.boxTopMargin
		boxSpanX = lb.Conf.boxSpanX
		boxSpanY = lb.Conf.boxSpanY
	)

	// show list box
	//termbox.SetCell(boxX-1, boxY-1, '┌', color, color)
	//termbox.SetCell(boxX-1, boxY+1, '│', color, color)
	//termbox.SetCell(boxX+boxSpanX, boxY+1, '│', color, color)
	//termbox.SetCell(boxX-1, boxY, '│', color, color)
	termbox.SetCell(boxX-1, boxY+boxSpanY, '└', color, termbox.ColorDefault)
	termbox.SetCell(boxX+boxSpanX, boxY+boxSpanY, '┘', color, termbox.ColorDefault)
	//termbox.SetCell(boxX+boxSpanX, boxY-1, '┐', color, color)
	//termbox.SetCell(boxX+boxSpan, boxY, '│', color, color)
	Fill(boxX-1, boxY+1, 1, boxSpanY, termbox.Cell{Ch: '│', Fg: color})
	Fill(boxX+boxSpanX, boxY+1, 1, boxSpanY, termbox.Cell{Ch: '│', Fg: color})
	Fill(boxX, boxY+boxSpanY, boxSpanX, 1, termbox.Cell{Ch: '─', Fg: color})
	//Fill(boxX, boxY-1, boxSpan, 1, termbox.Cell{Ch: '─'})
	//Fill(boxX, boxY+1, boxSpan, 1, termbox.Cell{Ch: '─'})

	TBPrint(boxX, boxY+1, termbox.ColorDefault, termbox.ColorDefault, fmt.Sprintf("%-48s", "(*) VPN密码"))
	TBPrint(boxX, boxY+2, termbox.ColorDefault, termbox.ColorDefault, fmt.Sprintf("%-48s", "(*) Bak00密码"))

	termbox.Flush()

	// cache
	lb.boxX = boxX
	lb.boxY = boxY
}

//func (lb *ListBox) watch() {
//	for {
//		switch ev := termbox.PollEvent(); ev.Type {
//		case termbox.EventKey:
//			log.Debug("ListBox receive key %d\n", ev.Key)
//			switch ev.Key {
//			case termbox.KeyEsc:
//				return
//			}
//		case termbox.EventError:
//			panic(ev.Err)
//		}
//	}
//}
