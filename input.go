package termloop

import "github.com/gdamore/tcell/v3"

type input struct {
	endKey tcell.Key
	eventQ chan tcell.Event
	ctrl   chan bool
}

func newInput() *input {
	i := input{eventQ: make(chan tcell.Event),
		ctrl:   make(chan bool, 2),
		endKey: tcell.Key(tcell.KeyCtrlC)}
	return &i
}

func (i *input) start(s *Screen) {
	go poll(i, s)
}

func (i *input) stop() {
	i.ctrl <- true
}

func poll(i *input, s *Screen) {
	eventChannel := s.EventQ()
loop:
	for {
		select {
		case <-i.ctrl:
			break loop
		case ev := <-eventChannel:
			i.eventQ <- ev
		}
	}
}
