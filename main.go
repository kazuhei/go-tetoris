package main

import termbox "github.com/nsf/termbox-go"

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	var stage Stage
	stage = newStage(10, 20)
	stage.Start()
}
