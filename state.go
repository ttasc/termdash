package main

type GameStatus int

const (
	StatusPlaying GameStatus = iota
	StatusPaused
	StatusLost
)

type ObstacleType int

const (
	TypeSpike ObstacleType = iota
	TypeBlock
	TypeBird
	TypeItem
)

type Player struct {
	Y           float64
	VelocityY   float64
	IsJumping   bool
	CrouchTimer float64
}

type Obstacle struct {
	Type   ObstacleType
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type Particle struct {
	X       float64
	Y       float64
	VelX    float64
	VelY    float64
	Life    float64
	MaxLife float64
	Char    rune
	Color   int
}

type GameState struct {
	Status GameStatus

	Player    *Player
	Obstacles []*Obstacle
	Particles []*Particle

	Score     int
	HighScore int

	Speed      float64
	SpawnTimer float64
	Distance   float64

	ParticleTimer float64

	TermW int
	TermH int
}

func NewGameState() *GameState {
	gs := &GameState{}
	gs.Reset()
	return gs
}

func (gs *GameState) Reset() {
	gs.Status = StatusPlaying
	gs.Obstacles = make([]*Obstacle, 0)
	gs.Particles = make([]*Particle, 0)
	gs.Score = 0
	gs.Distance = 0
	gs.Speed = 25.0
	gs.SpawnTimer = 1.5
	gs.ParticleTimer = 0.0
	gs.Player = &Player{VelocityY: 0, IsJumping: false, Y: 0}
}
