package main

import (
	"fmt"

	"github.com/ttasc/ttbox"
)

const (
	MinTermW = 50
	MinTermH = 15

	ColorPlayer  = 51
	ColorSpike   = 196
	ColorBlock   = 245
	ColorGround  = 239
	ColorText    = 250
	ColorTextDim = 240
	ColorAccent  = 220
	ColorLose    = 196

	BaseLineOffset = 3
	BaseLineChar   = '━'
)

func render(gs *GameState) {
	ttbox.Clear()

	if gs.TermW < MinTermW || gs.TermH < MinTermH {
		ttbox.DrawTextCenter(gs.TermH/2, " TERMINAL TOO SMALL! (Min 50x15) ", ColorText, ttbox.ColorDefault)
		ttbox.Present()
		return
	}

	drawGround(gs)
	drawObstacles(gs)

	if gs.Status != StatusLost {
		drawPlayer(gs)
	}

	drawParticles(gs)
	drawHeader(gs)

	if gs.Status != StatusPlaying {
		drawGameBanner(gs)
	} else {
		drawControlsGuide(gs.TermH)
	}

	ttbox.Present()
}

func drawParticles(gs *GameState) {
	for _, p := range gs.Particles {
		x := int(p.X)
		y := int(p.Y)

		if x >= 0 && x < gs.TermW && y >= 0 && y < gs.TermH {
			color := p.Color

			if p.Life/p.MaxLife < 0.4 {
				color = ColorGround
			}
			ttbox.SetCell(x, y, p.Char, color, ttbox.ColorDefault)
		}
	}
}

func drawHeader(gs *GameState) {
	speedMultiplier := gs.Speed / 25.0
	header := fmt.Sprintf(" SCORE: %d | HI-SCORE: %d | SPEED: %.1fx ", gs.Score, gs.HighScore, speedMultiplier)
	ttbox.DrawTextCenter(0, header, ColorText, ttbox.ColorDefault)
}

func drawGround(gs *GameState) {
	groundY := gs.TermH - BaseLineOffset

	for x := 0; x < gs.TermW; x++ {
		ttbox.SetCell(x, groundY, BaseLineChar, ColorGround, ttbox.ColorDefault)
	}

	offset := int(gs.Distance * 0.5)

	for y := groundY + 1; y < gs.TermH; y++ {
		for x := 0; x < gs.TermW; x++ {
			if (x+y+offset)%3 == 0 {
				ttbox.SetCell(x, y, '.', ColorTextDim, ttbox.ColorDefault)
			}
		}
	}
}

func drawObstacles(gs *GameState) {
	for _, obs := range gs.Obstacles {
		x := int(obs.X)
		y := int(obs.Y)

		if x < 0 || x >= gs.TermW {
			continue
		}

		switch obs.Type {
		case TypeSpike:
			ttbox.SetCell(x, y, '▲', ColorSpike, ttbox.ColorDefault)
		case TypeBlock:
			for dy := 0; dy < int(obs.Height); dy++ {
				for dx := 0; dx < int(obs.Width); dx++ {
					drawX := x + dx
					drawY := y + dy
					if drawX >= 0 && drawX < gs.TermW {
						ttbox.SetCell(drawX, drawY, '█', ColorBlock, ttbox.ColorDefault)
					}
				}
			}
		case CeilingTrap:
			drawW := int(obs.Width)

			for dy := 0; dy < 3; dy++ {
				for dx := 0; dx < drawW; dx++ {
					drawX := x + dx
					drawY := y + dy + 1
					if drawX >= 0 && drawX < gs.TermW && drawY >= 0 && drawY < gs.TermH {
						if dy == 2 {
							ttbox.SetCell(drawX, drawY, '▼', ColorSpike, ttbox.ColorDefault)
						} else {
							ttbox.SetCell(drawX, drawY, '█', ColorBlock, ttbox.ColorDefault)
						}
					}
				}
			}
		case TypeItem:
			ttbox.SetCell(x, y, '◆', ColorAccent, ttbox.ColorDefault)
		}
	}
}

func drawPlayer(gs *GameState) {
	px := gs.TermW / 4
	py := int(gs.Player.Y)

	char := '■'
	if gs.Player.CrouchTimer > 0 && !gs.Player.IsJumping {
		char = '_'
	}

	ttbox.SetCell(px, py, char, ColorPlayer, ttbox.ColorDefault)
}

func drawControlsGuide(termH int) {
	ttbox.DrawTextCenter(termH-2, " GEOMETRY DASH TERMINAL ", ColorAccent, ttbox.ColorDefault)
	ttbox.DrawTextCenter(termH-1, "[SPACE/UP]: Jump | [DOWN/S]: Fast Fall |[ESC]: Pause ", ColorTextDim, ttbox.ColorDefault)
}

func drawGameBanner(gs *GameState) {
	msgFg := ColorAccent
	msg := " * PAUSED * "
	subMsg := " [R] Restart   [ESC] Resume[Q] Quit "

	if gs.Status == StatusLost {
		msgFg = ColorLose
		msg = " * CRASHED! * "
		subMsg = fmt.Sprintf(" Final Score: %d | [R] Play Again   [ESC/Q] Exit ", gs.Score)
	}

	ttbox.ClearRect(0, gs.TermH-2, gs.TermW, 2)

	ttbox.SetAttr(true, false, false, false)
	ttbox.DrawTextCenter(gs.TermH-2, msg, msgFg, ttbox.ColorDefault)
	ttbox.ResetAttr()

	ttbox.DrawTextCenter(gs.TermH-1, subMsg, ColorTextDim, ttbox.ColorDefault)
}
