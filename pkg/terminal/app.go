package terminal

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"

	"temi/pkg"
)

func Run(loader pkg.MemStatsLoader) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	controller := newController()

	ev := ui.PollEvents()
	tick := time.Tick(time.Second)

	for {
		select {
		case e := <-ev:
			switch e.Type {
			case ui.KeyboardEvent:
				// quit on any keyboard event
				return
			case ui.ResizeEvent:
				controller.Resize()
			}
		case <-tick:
			stat, err := loader.Load()
			if err != nil {
				log.Println(err)
				break
			}
			// update dashboard every second
			controller.Render(stat)
		}
	}
}
