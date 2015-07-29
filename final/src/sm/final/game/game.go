// Copyright (C) 2015 Space Monkey, Inc.

package game

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	std_errors "errors"
	"fmt"
	math_rand "math/rand"
	"strings"
	"sync"
	"time"

	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/spacelog"

	"sm/final/grid"
	"sm/final/renderer"
)

var (
	GameError = errors.NewClass("game error", errors.NoCaptureStack())
	JoinError = GameError.NewClass("join error")

	logger = spacelog.GetLogger()
)

type Command string

const (
	selfDestruct Command = "kaboom"
	Join         Command = "join"
	Noop         Command = "noop"
	MoveForward  Command = "move"
	RotateLeft   Command = "left"
	RotateRight  Command = "right"
	FireLaser    Command = "fire"
)

func (c Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

func (c *Command) UnmarshalJSON(p []byte) (err error) {
	var raw string
	err = json.Unmarshal(p, &raw)
	if err != nil {
		return err
	}
	*c, err = CommandFromString(raw)
	return nil
}

func CommandFromString(s string) (Command, error) {
	switch Command(strings.ToLower(s)) {
	case Join:
		return Join, nil
	case Noop:
		return Noop, nil
	case MoveForward:
		return MoveForward, nil
	case RotateLeft:
		return RotateLeft, nil
	case RotateRight:
		return RotateRight, nil
	case FireLaser:
		return FireLaser, nil
	}
	return "", std_errors.New(fmt.Sprintf("%s is not a valid command", s))
}

type GameStatus string

const (
	Running GameStatus = "running"
	Lost    GameStatus = "lost"
	Won     GameStatus = "won"
	Draw    GameStatus = "draw"
)

func (s GameStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *GameStatus) UnmarshalJSON(p []byte) (err error) {
	var raw string
	err = json.Unmarshal(p, &raw)
	if err != nil {
		return err
	}
	switch GameStatus(strings.ToLower(raw)) {
	case Running:
		*s = Running
	case Lost:
		*s = Lost
	case Won:
		*s = Won
	case Draw:
		*s = Draw
	default:
		return std_errors.New(fmt.Sprintf("%s is not a valid game status", raw))
	}
	return nil
}

type TurnState struct {
	Status      GameStatus       `json:"status"`
	Health      int              `json:"health"`
	Energy      int              `json:"energy"`
	Orientation grid.Orientation `json:"orientation"`
	Grid        string           `json:"grid"`
	Config      *GameConfig      `json:"config,omitempty"`
}

type Player struct {
	Id          string
	Moniker     string
	Orientation grid.Orientation
	Owner       grid.Owner
	Health      int
	Energy      int
	Coord       grid.Coord
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%s)", p.Owner, p.Moniker)
}

func (p *Player) ToCell() grid.Cell {
	return grid.Cell{
		Type:        grid.Player,
		Orientation: p.Orientation,
		Owner:       p.Owner,
	}
}

func (p *Player) Alive() bool {
	return p.Health > 0
}

func (p *Player) Hit(damage int) (alive bool) {
	p.Health -= damage
	if p.Health < 0 {
		p.Health = 0
	}
	return p.Alive()
}

type Laser struct {
	Coord       grid.Coord
	Owner       grid.Owner
	Orientation grid.Orientation
	Lifetime    int
}

func (l *Laser) ToCell() grid.Cell {
	return grid.Cell{
		Type:        grid.Laser,
		Owner:       l.Owner,
		Orientation: l.Orientation,
	}
}

type Battery struct {
	Coord grid.Coord
}

func (e *Battery) ToCell() grid.Cell {
	return grid.Cell{
		Type: grid.Battery,
	}
}

type playerAction struct {
	player  *Player
	command Command
	statech chan TurnState
}

type Game struct {
	mtx        sync.Mutex
	config     *Config
	renderer   renderer.Renderer
	turn       int
	players    []*Player
	grid       *grid.Grid
	actionsch  chan playerAction
	rand       *math_rand.Rand
	lasers     []*Laser
	explosions []grid.Coord
	batteries  []Battery
}

func NewGame(config *Config, renderer renderer.Renderer, done_callback func()) *Game {
	if config == nil {
		config = DefaultConfig()
	}

	logger.Noticef("new game: config=%+v", config)
	g := &Game{
		config:    config,
		renderer:  renderer,
		actionsch: make(chan playerAction),
		rand:      math_rand.New(math_rand.NewSource(time.Now().UnixNano())),
	}

	if config.GridFile != "" {
		var err error
		g.grid, err = grid.LoadFromFile(config.GridFile)
		logger.Errore(err)
	}

	if g.grid == nil {
		g.grid = grid.NewRandom(config.Width, config.Height, config.Walls,
			config.Enclosed)
	}

	go g.run(done_callback)
	return g
}

func (g *Game) Join(moniker string) (id string, state TurnState, err error) {
	id, statech, err := g.join(moniker)
	if err != nil {
		return "", TurnState{}, nil
	}
	state = <-statech
	state.Config = &GameConfig{
		TurnTimeout:        int64(g.config.TurnTimeout),
		ConnectBackTimeout: int64(g.config.ConnectBackTimeout),
		MaxHealth:          g.config.MaxPlayerHealth,
		MaxEnergy:          g.config.MaxPlayerEnergy,
		HealthLoss:         g.config.HealthLoss,
		LaserDamage:        g.config.LaserDamage,
		LaserDistance:      g.config.LaserLifetime,
		LaserEnergy:        g.config.LaserEnergy,
		BatteryPower:       g.config.BatteryPower,
		BatteryHealth:      g.config.BatteryHealth,
	}
	return id, state, nil
}

type GameConfig struct {
	TurnTimeout        int64 `json:"turn_timeout"`
	ConnectBackTimeout int64 `json:"connect_back_timeout"`
	MaxHealth          int   `json:"max_health"`
	MaxEnergy          int   `json:"max_energy"`
	HealthLoss         int   `json:"health_loss"`
	LaserDamage        int   `json:"laser_damage"`
	LaserDistance      int   `json:"laser_distance"`
	LaserEnergy        int   `json:"laser_energy"`
	BatteryPower       int   `json:"battery_power"`
	BatteryHealth      int   `json:"battery_health"`
}

func (g *Game) join(moniker string) (id string, statech <-chan TurnState,
	err error) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if len(g.players) >= g.config.NumPlayers {
		return "", nil, JoinError.New("only %d players allowed",
			g.config.NumPlayers)
	}

	coord, ok := g.randomEmptyCell()
	if !ok {
		panic("grid does not have enough empty cells to place a player!")
	}

	id = newId()
	player := &Player{
		Id:          id,
		Moniker:     moniker,
		Owner:       grid.Owner(len(g.players) + 1),
		Health:      g.config.PlayerHealth,
		Energy:      g.config.PlayerEnergy,
		Coord:       coord,
		Orientation: grid.North,
	}

	g.players = append(g.players, player)

	statech = g.submitAction(player, Join)

	return id, statech, nil
}

func (g *Game) TakeTurn(id string, command Command) (state TurnState,
	err error) {
	statech, err := g.takeTurn(id, command)
	if err != nil {
		return TurnState{}, nil
	}
	return <-statech, nil
}

func (g *Game) takeTurn(id string, command Command) (statech <-chan TurnState,
	err error) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	player := g.findPlayerById(id)
	if player == nil {
		return nil, GameError.New("no such player %q", id)
	}
	return g.submitAction(player, command), nil
}

func (g *Game) submitAction(player *Player, command Command) (
	statech chan TurnState) {

	// chan must be buffered so we don't hang up the run loop.
	statech = make(chan TurnState, 1)

	if g.isStarted() && g.done() {
		statech <- g.turnState(player)
	} else {
		// append the playerAction and signal the run loop
		g.actionsch <- playerAction{
			player:  player,
			command: command,
			statech: statech,
		}
	}

	return statech
}

func (g *Game) isStarted() bool {
	return len(g.players) >= g.config.NumPlayers
}

func (g *Game) turnState(player *Player) TurnState {
	return TurnState{
		Status:      g.playerStatus(player),
		Health:      player.Health,
		Energy:      player.Energy,
		Orientation: player.Orientation,
		Grid:        g.grid.SerializeFor(player.Owner),
	}
}

func (g *Game) done() bool {
	return g.aliveCount() < 2
}

func (g *Game) aliveCount() (count int) {
	for _, player := range g.players {
		if player.Alive() {
			count++
		}
	}
	return count
}

func (g *Game) getWinner() (winner *Player, ok bool) {
	count := 0
	for _, player := range g.players {
		if player.Alive() {
			winner = player
			count++
		}
	}
	if count > 1 {
		winner = nil
	}
	return winner, count <= 1

}

func (g *Game) playerStatus(target *Player) GameStatus {
	alive := g.aliveCount()
	switch {
	case alive == 0:
		return Draw
	case alive > 1:
		return Running
	case target.Alive():
		return Won
	default:
		return Lost
	}
}

func (g *Game) randomEmptyCell() (grid.Coord, bool) {
	nogood := map[grid.Coord]bool{}
	for _, player := range g.players {
		if player.Alive() {
			nogood[player.Coord] = true
		}
	}
	for _, battery := range g.batteries {
		nogood[battery.Coord] = true
	}
	for _, laser := range g.lasers {
		nogood[laser.Coord] = true
	}

	candidates := make([]grid.Coord, 0, g.grid.Width()*g.grid.Height())
	for y := 0; y < g.grid.Height(); y++ {
		for x := 0; x < g.grid.Width(); x++ {
			coord := grid.Coord{X: x, Y: y}
			cell := g.grid.CellAt(coord)
			if cell.Type != grid.Empty || nogood[coord] {
				continue
			}
			candidates = append(candidates, coord)
		}
	}

	if len(candidates) == 0 {
		return grid.Coord{}, false
	}
	return candidates[g.rand.Intn(len(candidates))], true
}

func (g *Game) findPlayerById(id string) *Player {
	for _, player := range g.players {
		if player.Id == id {
			return player
		}
	}
	return nil
}

func (g *Game) findPlayerByOwner(owner grid.Owner) *Player {
	for _, player := range g.players {
		if player.Owner == owner {
			return player
		}
	}
	return nil
}

func (g *Game) run(done_callback func()) {

	// Wait for players to join
	var actions []*playerAction
waiting_for_players:
	for {
		g.renderMessage(renderer.Generic,
			"waiting for players (%d/%d have joined)", len(actions),
			g.config.NumPlayers)
		select {
		case action := <-g.actionsch:
			logger.Noticef("action: %s %s", action.player, action.command)
			actions = append(actions, &action)
			if len(actions) >= g.config.NumPlayers {
				break waiting_for_players
			}
		}
	}

	g.renderMessage(renderer.GameStart, "Ready? Fight!")
	time.Sleep(time.Second)
	g.renderGrid()

	g.mtx.Lock()
	g.sendState(actions)
	g.turn++
	g.mtx.Unlock()

	ticker := time.NewTicker(time.Second / 20)
	defer ticker.Stop()

	for done := false; !done; {
		// collect actions
		start_time := time.Now()
		ignore_actions := false
		actions = actions[:0]
	wait_for_actions:
		for len(actions) < g.aliveCount() {
			select {
			// get an action from a player
			case action := <-g.actionsch:
				if ignore_actions {
					action.command = Noop
				}
				if existing := findPlayerAction(actions, action.player); existing != nil {
					logger.Noticef("%s tried more than one action; self-destruct",
						action.player)
					existing.command = selfDestruct
				} else {
					logger.Noticef("received command %v for %s",
						action.command, action.player)
					actions = append(actions, &action)
				}
			case <-ticker.C:
				elapsed := time.Now().Sub(start_time)
				if g.config.TurnTimeout > 0 &&
					elapsed > g.config.TurnTimeout &&
					!ignore_actions {
					logger.Noticef("%s passed;  ignoring remaining actions from players", g.config.TurnTimeout)
					ignore_actions = true
					g.renderMessage(renderer.Generic, "someone's not responding...")
				}
				switch {
				case g.config.ConnectBackTimeout <= 0:
					if ignore_actions {
						break wait_for_actions
					}
				case elapsed > g.config.ConnectBackTimeout:
					// self-destruct any players that haven't submitted an action
					for _, player := range g.players {
						if !hasPlayerAction(actions, player) {
							logger.Noticef("%s failed to take turn in %s; self-destruct", player, g.config.ConnectBackTimeout)
							actions = append(actions, &playerAction{
								player:  player,
								command: selfDestruct,
							})
						}
					}
				}
			}
		}

		// Apply actions
		g.mtx.Lock()
		turn_ticks := 1
		if g.config.TurnTicks > 0 {
			turn_ticks = g.config.TurnTicks
		}
		for i := 0; !g.done() && i < turn_ticks; i++ {
			g.tick(i == 0, actions)
			g.renderGrid()
			if winner, ok := g.getWinner(); ok {
				if winner == nil {
					g.renderMessage(renderer.GameOver,
						"It's a draw :(")
				} else {
					g.renderMessage(renderer.GameOver,
						fmt.Sprintf("%s wins!", winner))
				}
			}
		}
		g.sendState(actions)
		g.turn++
		done = g.done()
		g.mtx.Unlock()
	}

	done_callback()
}

func (g *Game) sendState(actions []*playerAction) {
	for _, action := range actions {
		if action.statech != nil {
			action.statech <- g.turnState(action.player)
		}
	}
}

func (g *Game) tick(first bool, actions []*playerAction) {
	g.clearObjects()

	logger.Noticef("[turn %d] first=%t", g.turn, first)

	// clear explosions
	g.explosions = g.explosions[:0]

	// Do player actions
	if first {
		for _, player := range g.players {
			if !player.Hit(g.config.HealthLoss) {
				g.newExplosion(player.Coord)
			}
		}
		player_moves := map[grid.Coord][]*Player{}
		for _, pa := range actions {
			player, command := pa.player, pa.command
			logger.Noticef("executing %s for %s", pa.command, pa.player)
			if !player.Alive() {
				continue
			}

			switch command {
			case Noop:
			case selfDestruct:
				player.Health = 0
				g.newExplosion(player.Coord)
			case MoveForward:
				target_cell, target_coord := g.grid.CellRelativeTo(player.Coord,
					player.Orientation)
				if target_cell.Type != grid.Wall {
					player_moves[target_coord] = append(player_moves[target_coord], player)
				}
			case RotateLeft:
				player.Orientation.RotateLeft()
			case RotateRight:
				player.Orientation.RotateRight()
			case FireLaser:
				// add the laser starting at the players coordinates... it
				// will move when lasers are handled with below.  The lifetime
				// is the default lifetime + 1 since it will be decremented
				// below
				if player.Energy >= g.config.LaserEnergy {
					player.Energy -= g.config.LaserEnergy
					g.lasers = append(g.lasers, &Laser{
						Coord:       player.Coord,
						Lifetime:    g.config.LaserLifetime + 1,
						Owner:       player.Owner,
						Orientation: player.Orientation,
					})
				}
			}
		}

		// reconcile the player moves
	player_check:
		for coord, players := range player_moves {
			// make sure another player doesn't already occupy the spot
			for _, player := range g.players {
				if player.Alive() && player.Coord == coord {
					continue player_check
				}
			}

			// if multiple players tried to do the same thing, let's let a random one
			// win
			player := players[g.rand.Intn(len(players))]

			// did they run into a laser that is heading towards them?
			for i, laser := range g.lasers {
				if !(laser.Coord == coord &&
					laser.Orientation == player.Orientation.Opposite()) {
					continue
				}

				g.lasers = append(g.lasers[:i], g.lasers[i+1:]...)
				player.Hit(g.config.LaserDamage)
				g.newExplosion(coord)
				break
			}

			// did they pick up a battery?
			for i, battery := range g.batteries {
				if battery.Coord != coord {
					continue
				}
				g.batteries = append(g.batteries[:i], g.batteries[i+1:]...)
				player.Energy += g.config.BatteryPower
				player.Health += g.config.BatteryHealth
				if player.Energy > g.config.MaxPlayerEnergy {
					player.Energy = g.config.MaxPlayerEnergy
				}
				break
			}
			player.Coord = coord
		}
	}

	// Fire/expire lasers
	laser_moves := map[grid.Coord][]*Laser{}
	for _, laser := range g.lasers {
		laser.Lifetime--
		if laser.Lifetime < 1 {
			continue
		}
		target_cell, target_coord := g.grid.CellRelativeTo(laser.Coord,
			laser.Orientation)
		if target_cell.Type != grid.Wall {
			laser_moves[target_coord] = append(laser_moves[target_coord],
				laser)
		} else {
			// laser hit a wall
			g.newExplosion(target_coord)
		}
	}
	existing_lasers := append([]*Laser(nil), g.lasers...)
	g.lasers = g.lasers[:0]
laser_check:
	for coord, lasers := range laser_moves {
		if len(lasers) != 1 {
			// lasers collided, neither one lives... put an explosion
			g.newExplosion(coord)
			continue
		}
		laser := lasers[0]

		// see if the laser overlaps another laser
		for _, other_laser := range existing_lasers {
			if other_laser.Coord == coord &&
				other_laser.Orientation.Opposite() == laser.Orientation {
				// cull the other laser out of the laser_moves list, since it
				// is now dead.
				for j, laser2 := range laser_moves[laser.Coord] {
					if laser2 == other_laser {
						laser_moves[laser.Coord] = append(
							laser_moves[laser.Coord][:j],
							laser_moves[laser.Coord][j+1:]...)
						break
					}
				}
				continue laser_check
			}
		}

		// see if the laser hits a player
		for _, player := range g.players {
			if !(player.Alive() && player.Coord == coord) {
				continue
			}
			player.Hit(g.config.LaserDamage)
			g.newExplosion(coord)
			continue laser_check
		}

		// see if the laser hits a battery
		for i, battery := range g.batteries {
			if battery.Coord == coord {
				g.batteries = append(g.batteries[:i], g.batteries[i+1:]...)
				g.newExplosion(coord)
				continue laser_check
			}
		}

		// laser didn't hit anything... move it
		laser.Coord = coord
		g.lasers = append(g.lasers, laser)
	}

	// Spawn battery packs
	if first && g.config.BatteryTicks > 0 && g.turn%g.config.BatteryTicks == 0 {
		// Make sure we're not exceeding the max number of batteries
		if g.config.MaxBatteries < 0 ||
			len(g.batteries) < g.config.MaxBatteries {
			if coord, ok := g.randomEmptyCell(); ok {
				g.batteries = append(g.batteries, Battery{
					Coord: coord,
				})
			}
		}
	}
}

func (g *Game) clearObjects() {
	// Clear all tracked objects; they will be placed again after all actions
	// are resolved
	for _, player := range g.players {
		if player.Alive() {
			g.grid.SetCell(player.Coord, grid.EmptyCell)
		}
	}
	for _, laser := range g.lasers {
		g.grid.SetCell(laser.Coord, grid.EmptyCell)
	}
	for _, battery := range g.batteries {
		g.grid.SetCell(battery.Coord, grid.EmptyCell)
	}
	for _, explosion_coord := range g.explosions {
		g.grid.SetCellExploding(explosion_coord, false)
	}
}

func (g *Game) setObjects() {
	// Place all tracked objects in order of rendering priority
	for _, battery := range g.batteries {
		g.grid.SetCell(battery.Coord, battery.ToCell())
	}
	for _, laser := range g.lasers {
		g.grid.SetCell(laser.Coord, laser.ToCell())
	}
	for _, player := range g.players {
		if player.Alive() {
			g.grid.SetCell(player.Coord, player.ToCell())
		}
	}
	// add explosions
	for _, explosion_coord := range g.explosions {
		g.grid.SetCellExploding(explosion_coord, true)
	}

}

func (g *Game) newExplosion(coord grid.Coord) {
	g.explosions = append(g.explosions, coord)
}

func (g *Game) renderMessage(msg_type renderer.MessageType, format string,
	args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if g.renderer != nil {
		logger.Errore(g.renderer.Message(message, msg_type))
	}
}

func (g *Game) renderGrid() {
	g.setObjects()
	if g.renderer == nil {
		return
	}

	var statuses []renderer.PlayerStatus
	for _, player := range g.players {
		statuses = append(statuses, renderer.PlayerStatus{
			Moniker: player.Moniker,
			Health:  ratio(player.Health, g.config.MaxPlayerHealth),
			Energy:  ratio(player.Energy, g.config.MaxPlayerEnergy),
		})
	}
	if len(statuses) >= 2 {
		logger.Errore(g.renderer.SetStatus(statuses[:2]))
	}
	logger.Errore(g.renderer.Update(g.grid.Cells()))
}

func ratio(n, d int) float64 {
	return float64(n) / float64(d)
}

func hasPlayerAction(actions []*playerAction, player *Player) bool {
	return findPlayerAction(actions, player) != nil
}

func findPlayerAction(actions []*playerAction, player *Player) *playerAction {
	for _, action := range actions {
		if action.player == player {
			return action
		}
	}
	return nil
}

func newId() string {
	var id [8]byte
	rand.Read(id[:])
	return hex.EncodeToString(id[:])
}
