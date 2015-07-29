// Copyright (C) 2015 Space Monkey, Inc.

package client

import (
	"fmt"
	"net/http"

	"sm/final/game"
	"sm/final/server"
)

type Client struct {
	http_client *http.Client
}

func New() *Client {
	return &Client{
		http_client: new(http.Client),
	}
}

func (c *Client) Join(host, game, moniker string) (session *Session,
	state game.TurnState, err error) {

	game_url := fmt.Sprintf("http://%s/game/%s", host, game)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/join", game_url), nil)
	if err != nil {
		return nil, state, ClientError.Wrap(err)
	}
	req.Header.Set(server.PlayerMonikerHeader, moniker)
	resp, err := c.http_client.Do(req)
	if err != nil {
		return nil, state, ClientError.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, state, ServerError.New("Unable to join: status=%d err=%s",
			resp.StatusCode, errBody(resp.Body))
	}
	state, err = readAndParseState(resp.Body)
	if err != nil {
		return nil, state, err
	}
	return newSession(c.http_client, game_url, resp.Header.Get(server.PlayerIdHeader)), state, nil
}
