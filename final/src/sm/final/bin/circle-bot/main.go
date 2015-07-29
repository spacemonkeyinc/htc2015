// Copyright (C) 2015 Space Monkey, Inc.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	game = flag.String("game", "http://localhost:8080/game/yay:screen",
		"game basepath")
	moniker = flag.String("moniker", "circle-bot", "moniker")

	playerId string
)

type state struct {
	Status string `json:"status"`
	Health int    `json:"health"`
	Energy int    `json:"energy"`
	Coord  struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"coord"`
	Orientation string `json:"orientation"`
	Grid        string `json:"grid"`
}

func safePost(action string) (rv *state) {
	req, err := http.NewRequest("POST", *game+"/"+action, nil)
	if err != nil {
		panic(err)
	}
	if playerId != "" {
		req.Header.Set("X-Sm-Playerid", playerId)
	} else {
		req.Header.Set("X-Sm-Playermoniker", *moniker)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}
	if playerId == "" {
		playerId = resp.Header.Get("X-Sm-Playerid")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	rv = &state{}
	err = json.Unmarshal(body, rv)
	if err != nil {
		panic(err)
	}
	if rv.Status != "running" {
		os.Exit(0)
	}
	return rv
}

func main() {
	flag.Parse()

	r := rand.New(rand.NewSource(time.Now().Unix()))

	safePost("join")
	for {
		safePost("right")
		if r.Intn(4) == 0 {
			safePost("fire")
		}
		if r.Intn(2) == 0 {
			safePost("move")
			safePost("move")
		}
	}
}
