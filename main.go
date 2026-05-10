package main

import (
	"fmt"
	"time"

	"github.com/ttasc/ttbox"
)

func main() {

	state := NewGameState()

	if err := ttbox.Init(); err != nil {
		fmt.Printf("Error initializing TUI: %v\n", err)
		return
	}

	defer ttbox.Close()

	ttbox.HideCursorFunc()

	lastTime := time.Now()
	isRunning := true

	for isRunning {

		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		if dt > 0.1 {
			dt = 0.1
		}

		termW, termH := ttbox.Size()
		state.TermW = termW
		state.TermH = termH

		evt, err := ttbox.PollEventTimeout(16 * time.Millisecond)

		if err == nil {

			if evt.Type == ttbox.EventResize {
				if state.Status == StatusPlaying {
					state.Status = StatusPaused
				}
			} else if evt.Type == ttbox.EventKey && evt.Key == ttbox.KeyCtrlC {

				isRunning = false
			} else {

				isRunning = handleInput(state, evt)
			}
		}

		if state.Status != StatusPaused {
			updateLogic(state, dt)
		}

		render(state)
	}
}
