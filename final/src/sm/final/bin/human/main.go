// Copyright (C) 2015 Space Monkey, Inc.

package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/spacemonkeygo/spacelog"

	"sm/final/client"
	"sm/final/game"
)

var (
	host     = flag.String("host", "localhost:8080", "host to connect to")
	gameName = flag.String("game", "yay:screen", "game name")
	moniker  = flag.String("moniker", "human", "moniker to use")

	logger = spacelog.GetLogger()
)

func main() {
	flag.Parse()
	spacelog.Setup("run-tests", spacelog.SetupConfig{
		Output: "stderr",
		Format: "{{.Message}}"})
	err := Main()
	if err != nil {
		logger.Errore(err)
		os.Exit(1)
	}
}

func Main() (err error) {
	logger.Noticef("connecting to %q", *host)

	client := client.New()

	user := make(chan string)
	session, _, err := client.Join(*host, *gameName, *moniker)
	if err != nil {
		return err
	}
	go func() {
		for {
			var command string
			select {
			case command = <-user:
			default:
				command = "noop"
			}
			var state game.TurnState
			var err error
			switch command {
			case "", "noop":
				state, err = session.NoOp()
			case "l", "left":
				state, err = session.RotateLeft()
			case "r", "right":
				state, err = session.RotateRight()
			case "m", "move":
				state, err = session.MoveForward()
			case "f", "fire":
				state, err = session.FireLaser()
			default:
				logger.Errorf("unknown command %q", command)
				continue
			}
			if err != nil || state.Status != "running" {
				os.Exit(0)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		user <- strings.ToLower(string(scanner.Bytes()))
	}
	logger.Errore(scanner.Err())
	return nil
}
