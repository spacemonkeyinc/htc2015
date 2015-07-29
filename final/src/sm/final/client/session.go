// Copyright (C) 2015 Space Monkey, Inc.

package client

import (
	"fmt"
	"net/http"

	"sm/final/game"
	"sm/final/server"
)

type Session struct {
	http_client *http.Client
	url         string
	player_id   string
}

func newSession(http_client *http.Client, url, player_id string) *Session {
	return &Session{
		http_client: http_client,
		url:         url,
		player_id:   player_id,
	}
}

func (s *Session) NoOp() (state game.TurnState, err error) {
	return s.sendCommand(game.Noop)
}

func (s *Session) RotateLeft() (state game.TurnState, err error) {
	return s.sendCommand(game.RotateLeft)
}

func (s *Session) RotateRight() (state game.TurnState, err error) {
	return s.sendCommand(game.RotateRight)
}

func (s *Session) MoveForward() (state game.TurnState, err error) {
	return s.sendCommand(game.MoveForward)
}

func (s *Session) FireLaser() (state game.TurnState, err error) {
	return s.sendCommand(game.FireLaser)
}

func (s *Session) sendCommand(command game.Command) (state game.TurnState,
	err error) {
	url := fmt.Sprintf("%s/%s", s.url, command)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return state, ClientError.Wrap(err)
	}
	req.Header.Set(server.PlayerIdHeader, s.player_id)
	resp, err := s.http_client.Do(req)
	if err != nil {
		return state, ClientError.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return state, ServerError.New("Unable to join: status=%d err=%s",
			resp.StatusCode, errBody(resp.Body))
	}
	state, err = readAndParseState(resp.Body)
	if err != nil {
		return state, err
	}
	return state, nil
}
