// Copyright (C) 2015 Space Monkey, Inc.

package main

import (
	"os"
	"time"

	"sm/final/grid"
	"sm/final/renderer"
	"sm/final/renderer/sdl"
)

func fatalif(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	sdl.Run(func() {
		err := sdl.BackgroundMusic("./final/sounds/music1.wav")
		if err != nil {
			panic(err)
		}
		r, err := sdl.NewRenderer("yay", 1200, 840, 2, time.Second/8)
		if err != nil {
			panic(err)
		}
		defer r.Close()

		fatalif(r.Message("Ready?", renderer.GameStart))
		time.Sleep(time.Second)
		fatalif(r.SetStatus([]renderer.PlayerStatus{
			{Moniker: "Player1", Health: 1, Energy: .5},
			{Moniker: "Player2", Health: 1, Energy: .5}}))
		cells := grid.NewRandom(24, 16, 7, true).Cells()
		cells[10][10] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.North,
			Owner:       grid.Owner(1)}
		cells[13][11] = grid.Cell{
			Type:  grid.Player,
			Owner: grid.Owner(2)}
		cells[8][10] = grid.Cell{
			Type: grid.Battery}
		fatalif(r.Update(cells))

		for i := 0; i < 2; i++ {
			cells[10-i][10] = grid.Cell{Type: grid.Empty}
			cells[9-i][10] = grid.Cell{
				Type:        grid.Player,
				Orientation: grid.North,
				Owner:       grid.Owner(1)}
			fatalif(r.Update(cells))
		}
		for i := 0; i < 3; i++ {
			cells[7-i][10] = grid.Cell{
				Type:        grid.Laser,
				Orientation: grid.North,
				Owner:       grid.Owner(1)}
			fatalif(r.Update(cells))
			cells[7-i][10] = grid.Cell{Type: grid.Empty}
		}

		cells[8][10] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.East,
			Owner:       grid.Owner(1)}
		fatalif(r.Update(cells))

		for i := 0; i < 3; i++ {
			cells[8][11+i] = grid.Cell{
				Type:        grid.Laser,
				Orientation: grid.East,
				Owner:       grid.Owner(1)}
			fatalif(r.Update(cells))
			cells[8][11+i] = grid.Cell{Type: grid.Empty}
		}

		cells[8][10] = grid.Cell{Type: grid.Empty}
		cells[8][11] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.East,
			Owner:       grid.Owner(1)}
		fatalif(r.Update(cells))
		cells[8][11] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.South,
			Owner:       grid.Owner(1)}
		fatalif(r.Update(cells))

		cells[8][11] = grid.Cell{Type: grid.Empty}
		cells[9][11] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.South,
			Owner:       grid.Owner(1)}
		fatalif(r.Update(cells))

		for i := 0; i < 3; i++ {
			cells[10+i][11] = grid.Cell{
				Type:        grid.Laser,
				Orientation: grid.South,
				Owner:       grid.Owner(1)}
			fatalif(r.Update(cells))
			cells[10+i][11] = grid.Cell{Type: grid.Empty}
		}
		cells[13][11] = grid.Cell{
			Exploding: true,
			Type:      grid.Player,
			Owner:     grid.Owner(2)}
		fatalif(r.Update(cells))
		cells[13][11] = grid.Cell{
			Exploding: false,
			Type:      grid.Player,
			Owner:     grid.Owner(2)}

		cells[9][11] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.West,
			Owner:       grid.Owner(1)}
		fatalif(r.Update(cells))

		for i := 0; i < 3; i++ {
			cells[9][10-i] = grid.Cell{
				Type:        grid.Laser,
				Orientation: grid.West,
				Owner:       grid.Owner(1)}
			fatalif(r.Update(cells))
			cells[9][10-i] = grid.Cell{Type: grid.Empty}
		}

		cells[9][11] = grid.Cell{
			Type:        grid.Player,
			Orientation: grid.North,
			Owner:       grid.Owner(1)}

		fatalif(r.Update(cells))

		r.WaitForQuit()
		os.Exit(0)
	})
}
