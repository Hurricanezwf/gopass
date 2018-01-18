package ui

import (
	"fmt"
	"time"

	"github.com/Hurricanezwf/gopass/log"
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

	// 列表选中的颜色
	boxColor termbox.Attribute

	// 列表一页中最多显示的数据条数, 默认5条
	pageDataNum int
}

func DefaultListBoxConfig() *ListBoxConfig {
	lb := &ListBoxConfig{
		boxTopMargin:  0,
		boxLeftMargin: 0,
		boxSpanX:      50,
		boxSpanY:      20,
		boxColor:      termbox.ColorYellow,
		pageDataNum:   5,
	}

	w, h := termbox.Size()
	lb.boxTopMargin = int(0.3*float64(h)) + 1
	lb.boxLeftMargin = (w - lb.boxSpanX) / 2
	return lb
}

type ListBox struct {
	// 如果conf为nil, 将使用默认配置
	Conf *ListBoxConfig

	// 数据列表
	data       []string
	curPageIdx int
	curDataIdx int
}

func NewListBox() *ListBox {
	return &ListBox{
		Conf:       DefaultListBoxConfig(),
		data:       make([]string, 12),
		curPageIdx: 0,
		curDataIdx: 0,
	}
}

func (lb *ListBox) Show() {
	lb.draw(0)
}

func (lb *ListBox) draw(pageNo int) {
	var (
		tmpCount      = 0
		colorSelected = lb.Conf.boxColor
		boxX          = lb.Conf.boxLeftMargin
		boxY          = lb.Conf.boxTopMargin
		fgColor       = termbox.ColorDefault
		bgColor       = termbox.ColorDefault
		boxSpanX      = lb.Conf.boxSpanX
		//boxSpanY = lb.Conf.boxSpanY
	)

	lb.data[0] = "VPN密码"
	lb.data[1] = "Bak00密码"
	lb.data[2] = "sdjfkldsjkfjsdkfjsdkjfklsdjfksdjfksjdklfjskldjfksdjfksjdklfjskldjfklsdjfjsdklfj"
	lb.data[3] = "测试密码3"
	lb.data[4] = "测试密码4"
	lb.data[5] = "测试密码5"
	lb.data[6] = "测试密码6"
	lb.data[7] = "测试密码7"
	lb.data[8] = "测试密码8"
	lb.data[9] = "测试密码9"
	lb.data[10] = "测试密码10"
	lb.data[11] = "测试密码11"

	pageStart, pageEnd, pageCount, err := lb.calcPageIdx(pageNo)
	if err != nil {
		log.Warn("ListBox, calcPageIdx failed, %v", err)
		return
	}
	log.Debug("pageCount=%d, pageStart=%d, pageEnd=%d", pageCount, pageStart, pageEnd)

	lb.Clear()

	for i := pageStart; i <= pageEnd; i++ {
		tmpCount++

		data := lb.data[i]
		if len(data) > boxSpanX-3 {
			data = data[:boxSpanX-3]
		}

		if i == lb.curDataIdx {
			fgColor = colorSelected
		} else {
			fgColor = termbox.ColorDefault
		}

		TBPrint(boxX-1, boxY+tmpCount*2, fgColor, bgColor, fmt.Sprintf("(*) %s", data))
	}

	termbox.Flush()
}

func (lb *ListBox) Clear() {
	var (
		boxX        = lb.Conf.boxLeftMargin
		boxY        = lb.Conf.boxTopMargin
		boxSpanX    = lb.Conf.boxSpanX
		pageDataNum = lb.Conf.pageDataNum
	)

	Fill(boxX-1, boxY+2, boxSpanX, 2*pageDataNum, termbox.Cell{Ch: ' '})
}

func (lb *ListBox) KeyArrowUpHandler() error {
	var (
		err       error
		dataIdx   = lb.curDataIdx - 1
		pageNo    = lb.curPageIdx
		pageCount = 0
		pageStart = 0
		pageEnd   = 0
	)

	if dataIdx < 0 {
		dataIdx = len(lb.data) - 1
	}

	// 最多跨一页搜索
	for i := 0; i < 2; i++ {
		log.Debug("Find page(%d), pageStart:%d, pageEnd:%d, pageCount:%d, dataIdx:%d", pageNo, pageStart, pageEnd, pageCount, dataIdx)

		pageStart, pageEnd, pageCount, err = lb.calcPageIdx(pageNo)
		if err != nil {
			log.Warn("ListBox, calcPageIdx failed, %v", err)
			return err
		}
		if dataIdx >= pageStart && dataIdx <= pageEnd {
			break
		}
		pageNo = lb.prePageNo(pageNo, pageCount)
	}

	lb.curDataIdx = dataIdx
	lb.curPageIdx = pageNo
	lb.draw(pageNo)

	return nil
}

func (lb *ListBox) KeyArrowDownHandler() error {
	var (
		err       error
		dataIdx   = lb.curDataIdx + 1
		pageNo    = lb.curPageIdx
		pageCount = 0
		pageStart = 0
		pageEnd   = 0
	)

	if dataIdx >= len(lb.data) {
		dataIdx = 0
	}

	// 最多跨一页搜索
	for i := 0; i < 2; i++ {
		log.Debug("Find page(%d), pageStart:%d, pageEnd:%d, pageCount:%d, dataIdx:%d", pageNo, pageStart, pageEnd, pageCount, dataIdx)

		pageStart, pageEnd, pageCount, err = lb.calcPageIdx(pageNo)
		if err != nil {
			log.Warn("ListBox, calcPageIdx failed, %v", err)
			return err
		}
		if dataIdx >= pageStart && dataIdx <= pageEnd {
			break
		}
		pageNo = lb.nextPageNo(pageNo, pageCount)
	}

	lb.curDataIdx = dataIdx
	lb.curPageIdx = pageNo
	lb.draw(pageNo)

	return nil
}

func (lb *ListBox) KeyEnterHandler() error {
	notify := NewNotify("Copy OK", time.Second)
	notify.Info()
	return nil
}

func (lb *ListBox) calcPageCount() int {
	var (
		pageCount   = 0
		pageDataNum = lb.Conf.pageDataNum
		dataCount   = len(lb.data)
	)

	if dataCount%pageDataNum > 0 {
		pageCount++
	}
	return pageCount + dataCount/pageDataNum
}

func (lb *ListBox) calcPageIdx(pageNo int) (pageStart, pageEnd, pageCount int, err error) {
	var (
		dataCount = 0
		pageSize  = lb.Conf.pageDataNum
	)

	pageCount = lb.calcPageCount()
	if pageNo > pageCount {
		return 0, 0, 0, fmt.Errorf("PageNo(%d) exceeds PageCount(%d)", pageNo, pageCount)
	}

	pageStart = pageNo * pageCount
	pageEnd = pageStart + pageSize - 1
	dataCount = len(lb.data)
	if pageStart >= dataCount {
		return 0, 0, 0, fmt.Errorf("PageStart(%d) >= DataCount(%d)", pageStart, dataCount)
	}
	if pageEnd >= dataCount {
		pageEnd = dataCount - 1
	}

	return pageStart, pageEnd, pageCount, nil
}

func (lb *ListBox) nextPageNo(curPageNo, pageCount int) int {
	curPageNo++
	if curPageNo > pageCount {
		curPageNo = 0
	}
	return curPageNo
}

func (lb *ListBox) prePageNo(curPageNo, pageCount int) int {
	curPageNo--
	if curPageNo < 0 {
		curPageNo = pageCount
	}
	return curPageNo
}
