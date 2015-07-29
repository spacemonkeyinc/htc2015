// Copyright (C) 2015 Space Monkey, Inc.

package renderer

import (
	"sm/final/grid"
)

type PlayerStatus struct {
	Moniker string
	Health  float64
	Energy  float64
}

type MessageType int

const (
	Generic   MessageType = 0
	GameStart MessageType = 1
	GameOver  MessageType = 2
)

type Renderer interface {
	Message(msg string, msgType MessageType) error
	SetStatus(status []PlayerStatus) error
	Update(cells [][]grid.Cell) error
}
