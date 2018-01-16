package main

import "github.com/Hurricanezwf/gopass/ui"

func main() {
	ui := ui.New()
	defer ui.Close()

	eb := ui.EditBox()
	eb.Show()
}
