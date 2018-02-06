package ui

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/password"
	"github.com/atotto/clipboard"
	termbox "github.com/nsf/termbox-go"
)

type ListBoxConfig struct {
	// 列表距离页面顶端的距离
	boxTopMargin int

	// 列表距离页面左侧的距离
	boxLeftMargin int

	// 列表宽度, 默认50px
	boxSpanX int
	// 列表高度, 默认20px
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

	ui *UI

	// 数据列表
	dataAll    [][]byte
	dataDraw   [][]byte
	curPageIdx int
	curDataIdx int

	// 控制列表匹配输入框内容,输入框内容发生改变，即向channel写
	matchC chan struct{}
}

func NewListBox() *ListBox {
	return &ListBox{
		Conf:       DefaultListBoxConfig(),
		curPageIdx: 0,
		curDataIdx: 0,
		matchC:     make(chan struct{}, 100),
	}
}

func (lb *ListBox) Open(ui *UI) {
	keys, err := password.ListKeys()
	if err != nil {
		log.Warn("Password list keys failed, %v", err)
		return
	}

	lb.ui = ui
	lb.dataAll = keys
	lb.dataDraw = lb.dataAll[:]

	// listen keys search
	go lb.match()

	lb.draw(0)
}

func (lb *ListBox) Close() {
	close(lb.matchC)
	lb.matchC = nil
}

// NotifyMatch: 告知ListBox编辑框的内容发生改变
func (lb *ListBox) NotifyMatch() {
	if lb.matchC == nil {
		return
	}
	select {
	case lb.matchC <- struct{}{}:
	default:
		log.Warn("Match channel is full")
	}
}

func (lb *ListBox) match() {
	var (
		changed bool
		ticker  = time.NewTicker(100 * time.Millisecond)
	)

	for {
		select {
		case _, ok := <-lb.matchC:
			if !ok {
				return
			}
			changed = true
		case <-ticker.C:
			if !changed {
				continue
			}

			changed = false
			v := lb.ui.EditBox.Value()
			lb.filter(v)
			lb.draw(0)
			//log.Debug("%s", string(v))
		}
	}
}

func (lb *ListBox) filter(key []byte) {
	if len(key) <= 0 {
		lb.dataDraw = lb.dataAll[:]
		return
	}

	pattern := fmt.Sprintf("(?i)%s", string(key))
	dataDraw := make([][]byte, 0)
	for _, d := range lb.dataAll {
		ok, err := regexp.Match(pattern, d)
		if err != nil {
			log.Warn("%s match pattern(%s) failed", string(d), pattern)
			continue
		}
		if ok {
			dataDraw = append(dataDraw, d)
		}
	}
	lb.dataDraw = dataDraw
	log.Debug("After filter, dataDraw size:%d", len(dataDraw))
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

	/*
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
	*/
	lb.Clear()

	pageStart, pageEnd, pageCount, err := lb.calcPageIdx(pageNo)
	if err != nil {
		log.Warn("ListBox, calcPageIdx failed, %v", err)
		return
	}
	log.Debug("pageCount=%d, pageStart=%d, pageEnd=%d", pageCount, pageStart, pageEnd)

	for i := pageStart; i <= pageEnd; i++ {
		tmpCount++

		data := lb.dataDraw[i]
		if len(data) > boxSpanX-4 {
			data = data[:boxSpanX-4]
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
		dataIdx = len(lb.dataDraw) - 1
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

	if dataIdx >= len(lb.dataDraw) {
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

// xsel or xclip will be needed
func (lb *ListBox) KeyEnterHandler() error {
	key := lb.dataDraw[lb.curDataIdx]
	pass, err := password.Get(key)
	if err != nil {
		log.Warn("Get password failed, %v", err)
		NewNotify("Copy Failed", time.Second).Warn()
		return err
	}
	if err = clipboard.WriteAll(string(pass)); err != nil {
		log.Warn("Copy to clipboard failed, %v", err)
		NewNotify("Copy Failed", time.Second).Warn()
		return err
	}
	NewNotify("Copy OK", time.Second).Info()
	return nil
}

func (lb *ListBox) calcPageCount() int {
	var (
		pageCount   = 0
		pageDataNum = lb.Conf.pageDataNum
		dataCount   = len(lb.dataDraw)
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
	dataCount = len(lb.dataDraw)
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
