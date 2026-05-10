package main

import (
	"math/rand"
)

const (
	Gravity         = 290.0
	MaxSpeed        = 75.0
	GroundOffset    = 3
	ScoreMultiplier = 5.0
	SpeedIncrement  = 1.5
	SkyLimit        = 2.0
)

func updateLogic(gs *GameState, dt float64) {
	for i := len(gs.Particles) - 1; i >= 0; i-- {
		p := gs.Particles[i]
		p.X += p.VelX * dt
		p.Y += p.VelY * dt
		p.Life -= dt
		if p.Life <= 0 || p.X < 0 || p.X > float64(gs.TermW) || p.Y > float64(gs.TermH) {
			gs.Particles = append(gs.Particles[:i], gs.Particles[i+1:]...)
		}
	}

	if gs.Status != StatusPlaying {
		return
	}

	gs.Distance += gs.Speed * dt
	gs.Score = int(gs.Distance / ScoreMultiplier)

	if gs.Speed < MaxSpeed {
		gs.Speed += SpeedIncrement * dt
	}

	groundY := float64(gs.TermH - GroundOffset - 1)
	playerX := float64(gs.TermW / 4)

	if gs.Player.Y == 0 {
		gs.Player.Y = groundY
	}

	gs.Player.VelocityY += Gravity * dt
	gs.Player.Y += gs.Player.VelocityY * dt

	if gs.Player.Y >= groundY {
		gs.Player.Y = groundY
		gs.Player.VelocityY = 0
		gs.Player.IsJumping = false
	}

	if gs.Player.Y < SkyLimit {
		gs.Player.Y = SkyLimit
		gs.Player.VelocityY = 0
	}

	if gs.Player.CrouchTimer > 0 {
		gs.Player.CrouchTimer -= dt
	}

	var hitboxY float64
	if gs.Player.CrouchTimer > 0 && !gs.Player.IsJumping {
		gs.Player.Height = 0.5
		hitboxY = gs.Player.Y + 0.5
	} else {
		gs.Player.Height = 1.0
		hitboxY = gs.Player.Y
	}

	if !gs.Player.IsJumping {
		gs.ParticleTimer -= dt
		if gs.ParticleTimer <= 0 {
			spawnTrailParticle(gs, playerX, gs.Player.Y)
			gs.ParticleTimer = 0.08
		}
	}

	gs.SpawnTimer -= dt
	if gs.SpawnTimer <= 0 {
		spawnObstacle(gs, groundY)
	}

	for i := len(gs.Obstacles) - 1; i >= 0; i-- {
		obs := gs.Obstacles[i]
		obs.X -= gs.Speed * dt

		if obs.X+obs.Width < 0 {
			gs.Obstacles = append(gs.Obstacles[:i], gs.Obstacles[i+1:]...)
			continue
		}

		var isHit bool
		if obs.Type == TypeItem {
			isHit = checkCollision(playerX, hitboxY, gs.Player.Width, gs.Player.Height, obs.X - 1.5, obs.Y - 1.5, obs.Width + 3.0, obs.Height + 3.0)
		} else {
			isHit = checkCollision(playerX, hitboxY, gs.Player.Width, gs.Player.Height, obs.X, obs.Y, obs.Width, obs.Height)
		}

		if isHit {
			if obs.Type == TypeItem {
				gs.Score += 25
				spawnItemExplosion(gs, obs.X, obs.Y)
				gs.Obstacles = append(gs.Obstacles[:i], gs.Obstacles[i+1:]...)
			} else {
				gs.Status = StatusLost
				if gs.Score > gs.HighScore {
					gs.HighScore = gs.Score
				}
				spawnDeathExplosion(gs, playerX, gs.Player.Y)
			}
		}
	}
}

func spawnObstacle(gs *GameState, groundY float64) {
	baseInterval := 1.5 - (gs.Speed/MaxSpeed)*0.7

	r := rand.Float64()
	obsType := TypeSpike
	if r > 0.4 && r <= 0.6 {
		obsType = TypeBlock
	} else if r > 0.6 && r <= 0.8 {
		obsType = CeilingTrap
	} else if r > 0.8 {
		obsType = TypeItem
	}

	var w, h, y float64
	var extraDelay float64

	switch obsType {
	case TypeSpike:
		w, h, y = 1.0, 1.0, groundY
		extraDelay = 0.0
	case TypeBlock:
		w, h, y = 2.0, 2.0, groundY-1.0
		extraDelay = 0.2
	case CeilingTrap:
		w, h, y = 3.0, 3.0, groundY-2.5
		extraDelay = 0.2
	case TypeItem:
		w, h, y = 1.0, 1.0, groundY-2.5
		extraDelay = -0.2
	}

	gs.SpawnTimer = baseInterval + rand.Float64()*0.4 + extraDelay

	gs.Obstacles = append(gs.Obstacles, &Obstacle{
		Type:   obsType,
		X:      float64(gs.TermW),
		Y:      y,
		Width:  w,
		Height: h,
	})
}

func spawnTrailParticle(gs *GameState, x, y float64) {
	gs.Particles = append(gs.Particles, &Particle{
		X:       x - 1,
		Y:       y, // + 0.5,
		VelX:    -gs.Speed * 0.8,
		VelY:    0.0, // -5.0 + rand.Float64()*10.0,
		Life:    0.3,
		MaxLife: 0.3,
		Char:    '.',
		Color:   245,
	})
}

func spawnDeathExplosion(gs *GameState, x, y float64) {
	chars :=[]rune{'x', '*', '+', ':', '■', '.', '\''}
	for i := 0; i < 25; i++ {
		gs.Particles = append(gs.Particles, &Particle{
			X:       x,
			Y:       y,
			VelX:    -40.0 + rand.Float64()*80.0,
			VelY:    -50.0 + rand.Float64()*100.0,
			Life:    1.0 + rand.Float64(),
			MaxLife: 2.0,
			Char:    chars[rand.Intn(len(chars))],
			Color:   196,
		})
	}
}

func spawnItemExplosion(gs *GameState, x, y float64) {
	for i := 0; i < 8; i++ {
		gs.Particles = append(gs.Particles, &Particle{
			X:       x,
			Y:       y,
			VelX:    -15.0 + rand.Float64()*30.0,
			VelY:    -20.0 + rand.Float64()*40.0,
			Life:    0.5,
			MaxLife: 0.5,
			Char:    '+',
			Color:   220,
		})
	}
}

func checkCollision(px, py, pw, ph, ox, oy, ow, oh float64) bool {
	leniency := 0.2
	return px+leniency < ox+ow && px+pw-leniency > ox && py+leniency < oy+oh && py+ph-leniency > oy
}
