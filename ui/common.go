package ui

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

//
//func byte_slice_grow(s []byte, desired_cap int) []byte {
//	if cap(s) < desired_cap {
//		ns := make([]byte, len(s), desired_cap)
//		copy(ns, s)
//		return ns
//	}
//	return s
//}

func ByteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = ByteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}

func ByteSliceRemove(text []byte, from, to int) []byte {
	sz := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-sz]
	return text
}

func ByteSliceGrow(b []byte, newCap int) []byte {
	if cap(b) >= newCap {
		return b
	}
	nb := make([]byte, len(b), newCap)
	copy(nb, b)
	return nb
}

func Fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
	termbox.Flush()
}

func TBPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
	termbox.Flush()
}
