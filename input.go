package main

import (
	"unicode"

	"github.com/ttasc/ttbox"
)

const (
	JumpForce     = -45.0
	FastFallForce = 100.0
)

func handleInput(gs *GameState, evt ttbox.Event) bool {
	if evt.Type != ttbox.EventKey {
		return true
	}

	switch gs.Status {
	case StatusPlaying:
		return handlePlayingInput(gs, evt)
	case StatusPaused:
		return handlePausedInput(gs, evt)
	case StatusLost:
		return handleGameOverInput(gs, evt)
	}

	return true
}

func handlePlayingInput(gs *GameState, evt ttbox.Event) bool {
	if evt.Key == ttbox.KeyEscape {
		gs.Status = StatusPaused
		return true
	}

	isJumpKey := evt.Key == ttbox.KeyArrowUp || evt.Ch == ' ' || unicode.ToLower(evt.Ch) == 'k'
	isDownKey := evt.Key == ttbox.KeyArrowDown || unicode.ToLower(evt.Ch) == 'j'

	if isDownKey {
		if gs.Player.IsJumping {
			if gs.Player.VelocityY < 0 {
				gs.Player.VelocityY = 0
			}
			gs.Player.VelocityY += FastFallForce
		} else {
			gs.Player.CrouchTimer = 0.8
		}
		return true
	}

	if isJumpKey {
		if !gs.Player.IsJumping {
			gs.Player.VelocityY = JumpForce
			gs.Player.IsJumping = true
			gs.Player.CrouchTimer = 0.0
		}
		return true
	}

	return true
}

func handlePausedInput(gs *GameState, evt ttbox.Event) bool {
	if evt.Key == ttbox.KeyEscape {
		gs.Status = StatusPlaying
		return true
	}
	if evt.Ch != 0 {
		ch := unicode.ToLower(evt.Ch)
		if ch == 'r' {
			gs.Reset()
			return true
		}
		if ch == 'q' {
			return false
		}
	}
	return true
}

func handleGameOverInput(gs *GameState, evt ttbox.Event) bool {
	if evt.Key == ttbox.KeyEscape {
		return false
	}
	if evt.Ch != 0 {
		ch := unicode.ToLower(evt.Ch)
		if ch == 'r' {
			gs.Reset()
			return true
		}
		if ch == 'q' {
			return false
		}
	}
	return true
}
