package main

import "github.com/Hurricanezwf/gopass/ui"

func main() {
	ui.Init()
	defer ui.Close()

	ui.DrawAll(nil)
}
