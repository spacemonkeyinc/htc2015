// Copyright (C) 2015 Space Monkey, Inc.

package game

import (
	"flag"
	"strings"
	"sync"
	"time"

	"sm/final/renderer"
	"sm/final/renderer/sdl"
)

var (
	renderTime = flag.Duration("renderer.time", time.Second/8,
		"render time per game cycle")
	screenWidth  = flag.Int("renderer.width", 1200, "width of screen")
	screenHeight = flag.Int("renderer.height", 840, "height of screen")
)

type Games struct {
	mtx    sync.Mutex
	games  map[string]*Game
	config *Config
}

func NewGames(config *Config) (
	games *Games) {
	return &Games{
		games:  map[string]*Game{},
		config: config}
}

func (g *Games) Lookup(name string) *Game {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.games[name]
}

func (g *Games) LookupOrCreate(name string) (*Game, error) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	game := g.games[name]
	if game != nil {
		return game, nil
	}

	var screen renderer.Renderer
	if strings.HasSuffix(name, ":screen") {
		sdl_screen, err := sdl.NewRenderer(name[:len(name)-len(":screen")],
			*screenWidth, *screenHeight, *numPlayers, *renderTime)
		if err != nil {
			return nil, err
		}
		go func() {
			sdl_screen.WaitForQuit()
			sdl_screen.Close()
		}()
		screen = sdl_screen
	}

	game = NewGame(g.config, screen, func() {
		g.mtx.Lock()
		delete(g.games, name)
		g.mtx.Unlock()
	})
	g.games[name] = game
	return game, nil
}
