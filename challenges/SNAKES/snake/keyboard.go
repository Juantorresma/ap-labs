package snake

import "github.com/nsf/termbox-go"

type keyboardEventType int

// Keyboard events
const (
	MOVE keyboardEventType = 1 + iota
	RETRY
	END
)

type keyboardEvent struct {
	eventType keyboardEventType
	key       termbox.Key
}

func keyToDirection(k termbox.Key) direction {
	switch k {
	case termbox.KeyArrowLeft:
		return LEFT
	case termbox.KeyArrowDown:
		return DOWN
	case termbox.KeyArrowRight:
		return RIGHT
	case termbox.KeyArrowUp:
		return UP
	default:
		return 0
	}
}

func listenToKeyboard(evChan chan keyboardEvent) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch kp := termbox.PollEvent(); kp.Type {
		case termbox.EventKey:
			switch kp.Key {
			case termbox.KeyArrowLeft:
				evChan <- keyboardEvent{eventType: MOVE, key: kp.Key}
			case termbox.KeyArrowDown:
				evChan <- keyboardEvent{eventType: MOVE, key: kp.Key}
			case termbox.KeyArrowRight:
				evChan <- keyboardEvent{eventType: MOVE, key: kp.Key}
			case termbox.KeyArrowUp:
				evChan <- keyboardEvent{eventType: MOVE, key: kp.Key}
			case termbox.KeyEsc:
				evChan <- keyboardEvent{eventType: END, key: kp.Key}
			default:
				if kp.Ch == 'r' {
					evChan <- keyboardEvent{eventType: RETRY, key: kp.Key}
				}
			}
		case termbox.EventError:
			panic(kp.Err)
		}
	}
}
