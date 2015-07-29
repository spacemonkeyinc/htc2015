// Copyright (C) 2015 Space Monkey, Inc.

package game

import (
	"flag"
	"time"
)

var (
	gridWidth          = flag.Int("logic.grid-width", 24, "width of grid")
	gridHeight         = flag.Int("logic.grid-height", 16, "height of grid")
	wallCount          = flag.Int("logic.grid-walls", 8, "# of walls")
	gridEnclosed       = flag.Bool("logic.grid-enclosed", false, "true if the grid should be enclosed")
	turnTimeout        = flag.Duration("logic.turn-timeout", time.Second/2, "timeout before player action is ignored for the turn")
	connectBackTimeout = flag.Duration("logic.connect-back-timeout", 10*time.Second, "timeout before we assume player has left the game")
	numPlayers         = flag.Int("logic.players", 2, "number of players in a game")
	turnTicks          = flag.Int("logic.turn-ticks", 2, "how many game ticks per player action")
	playerHealth       = flag.Int("logic.player-health", 300, "starting player health")
	maxPlayerHealth    = flag.Int("logic.player-max-health", 300, "max player health")
	playerEnergy       = flag.Int("logic.player-energy", 5, "starting player energy")
	maxPlayerEnergy    = flag.Int("logic.player-max-energy", 10, "max player energy")
	healthLoss         = flag.Int("logic.health-loss", 1, "health lost per turn")
	laserDamage        = flag.Int("logic.laser-damage", 50, "amount of damage when hit by a laser")
	laserLifetime      = flag.Int("logic.laser-distance", 32, "how many cells a laser travels")
	laserEnergy        = flag.Int("logic.laser-energy", 1, "how much energy need to fire a laser")
	batteryPower       = flag.Int("logic.battery-power", 5, "how much power a battery gives you")
	batteryHealth      = flag.Int("logic.battery-health", 20, "how much health a battery gives you")
	batteryTicks       = flag.Int("logic.battery-ticks", 15, "ticks between battery pack spawn")
	maxBatteries       = flag.Int("logic.max-batteries", 5, "the maximum number of batteries on the grid")
	gridFile           = flag.String("gridfile", "", "file containing grid to use")
)

type Config struct {
	Width              int
	Height             int
	Walls              int
	Enclosed           bool
	TurnTimeout        time.Duration
	ConnectBackTimeout time.Duration
	TurnTicks          int
	NumPlayers         int
	PlayerHealth       int
	MaxPlayerHealth    int
	PlayerEnergy       int
	MaxPlayerEnergy    int
	HealthLoss         int
	LaserDamage        int
	LaserLifetime      int
	LaserEnergy        int
	BatteryPower       int
	BatteryHealth      int
	BatteryTicks       int
	MaxBatteries       int
	GridFile           string
}

func DefaultConfig() *Config {
	return &Config{
		Width:              *gridWidth,
		Height:             *gridHeight,
		Walls:              *wallCount,
		Enclosed:           *gridEnclosed,
		TurnTimeout:        *turnTimeout,
		ConnectBackTimeout: *connectBackTimeout,
		TurnTicks:          *turnTicks,
		NumPlayers:         *numPlayers,
		PlayerHealth:       *playerHealth,
		MaxPlayerHealth:    *maxPlayerHealth,
		PlayerEnergy:       *playerEnergy,
		MaxPlayerEnergy:    *maxPlayerEnergy,
		HealthLoss:         *healthLoss,
		LaserDamage:        *laserDamage,
		LaserLifetime:      *laserLifetime,
		LaserEnergy:        *laserEnergy,
		BatteryPower:       *batteryPower,
		BatteryHealth:      *batteryHealth,
		BatteryTicks:       *batteryTicks,
		MaxBatteries:       *maxBatteries,
		GridFile:           *gridFile,
	}
}
