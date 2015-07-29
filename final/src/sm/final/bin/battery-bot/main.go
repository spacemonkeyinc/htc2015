// Copyright (C) 2015 Space Monkey, Inc.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	game = flag.String("game", "http://localhost:8080/game/yay:screen",
		"game basepath")
	moniker = flag.String("moniker", "battery-bot", "moniker")
	random  = rand.New(rand.NewSource(time.Now().Unix()))

	playerId string
	config   *configStruct
)

type configStruct struct {
	LaserEnergy int `json:"laser_energy"`
}

type state struct {
	Status      string `json:"status"`
	Health      int    `json:"health"`
	Energy      int    `json:"energy"`
	Orientation string `json:"orientation"`
	Grid        string `json:"grid"`
	GridParsed  [][]rune
	Config      *configStruct `json:"config"`
}

func parseGrid(grid string) (rv [][]rune) {
	for _, line := range strings.FieldsFunc(grid,
		func(r rune) bool { return r == '\n' }) {
		var row []rune
		for _, char := range line {
			row = append(row, char)
		}
		rv = append(rv, row)
	}
	return rv
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
	if rv.Config != nil {
		config = rv.Config
	}
	if config == nil {
		panic("no config")
	}
	rv.GridParsed = parseGrid(rv.Grid)
	if rv.Status != "running" {
		os.Exit(0)
	}
	return rv
}

func dist(m1, n1, m2, n2 int) float64 {
	m_dst := float64(m1) - float64(m2)
	n_dst := float64(n1) - float64(n2)
	return math.Sqrt(m_dst*m_dst + n_dst*n_dst)
}

func objectCoords(grid [][]rune, obj rune, closest_m, closest_n int) (m, n int) {
	best_m, best_n := -1, -1
	for m, row := range grid {
		for n, cell := range row {
			if cell == obj {
				if closest_m == -1 {
					return m, n
				}
				if best_m == -1 ||
					dist(best_m, best_n, closest_m, closest_n) >= dist(
						m, n, closest_m, closest_n) {
					best_m, best_n = m, n
				}
			}
		}
	}
	return best_m, best_n
}

func main() {
	flag.Parse()

	s := safePost("join")
	coords := func(obj rune) (int, int) {
		return objectCoords(s.GridParsed, obj, -1, -1)
	}
	closestCoords := func(obj rune, m, n int) (int, int) {
		return objectCoords(s.GridParsed, obj, m, n)
	}
	fire := func() {
		s = safePost("fire")
	}
	move := func() {
		s = safePost("move")
	}
	faceUp := func() {
		switch s.Orientation {
		case "south":
			s = safePost("left")
			s = safePost("left")
		case "east":
			s = safePost("left")
		case "north":
		case "west":
			s = safePost("right")
		default:
			panic("uhoh")
		}
	}
	faceDown := func() {
		switch s.Orientation {
		case "north":
			s = safePost("right")
			s = safePost("right")
		case "east":
			s = safePost("right")
		case "south":
		case "west":
			s = safePost("left")
		default:
			panic("uhoh")
		}
	}
	faceLeft := func() {
		switch s.Orientation {
		case "north":
			s = safePost("left")
		case "east":
			s = safePost("right")
			s = safePost("right")
		case "south":
			s = safePost("right")
		case "west":
		default:
			panic("uhoh")
		}
	}
	faceRight := func() {
		switch s.Orientation {
		case "north":
			s = safePost("right")
		case "south":
			s = safePost("left")
		case "west":
			s = safePost("left")
			s = safePost("left")
		case "east":
		default:
			panic("uhoh")
		}
	}
	random_plan_m, random_plan_n := -1, -1
	for {
		my_m, my_n := coords('X')
		if random_plan_m == my_m && random_plan_n == my_n {
			random_plan_m, random_plan_n = -1, -1
		}

		moveTowards := func(m, n int) {
			options := [](func() bool){
				func() bool {
					if m < my_m {
						faceUp()
						move()
						return true
					}
					return false
				},
				func() bool {
					if m > my_m {
						faceDown()
						move()
						return true
					}
					return false
				},
				func() bool {
					if n < my_n {
						faceLeft()
						move()
						return true
					}
					return false
				},
				func() bool {
					if n > my_n {
						faceRight()
						move()
						return true
					}
					return false
				},
			}
			for _, i := range random.Perm(4) {
				if options[i]() {
					return
				}
			}
		}

		other_m, other_n := coords('O')
		if my_m == -1 || other_m == -1 {
			panic("no people")
		}

		if s.Energy >= config.LaserEnergy {
			if other_m == my_m {
				if other_n < my_n {
					faceLeft()
					fire()
					continue
				}
				if other_n > my_n {
					faceRight()
					fire()
					continue
				}
			}
			if other_n == my_n {
				if other_m < my_m {
					faceUp()
					fire()
					continue
				}
				if other_m > my_m {
					faceDown()
					fire()
					continue
				}
			}
		}

		battery_m, battery_n := closestCoords('B', my_m, my_n)
		if battery_m != -1 {
			moveTowards(battery_m, battery_n)
			continue
		}

		if random_plan_m == -1 {
			random_plan_m = random.Intn(len(s.GridParsed))
			random_plan_n = random.Intn(len(s.GridParsed[0]))
		}
		moveTowards(random_plan_m, random_plan_n)
	}
}
