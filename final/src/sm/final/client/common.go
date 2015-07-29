// Copyright (C) 2015 Space Monkey, Inc.

package client

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/spacemonkeygo/errors"

	"sm/final/game"
)

var (
	ClientError = errors.NewClass("client error")
	ServerError = errors.NewClass("server error")
)

func errBody(r io.Reader) string {
	var msg [256]byte
	n, _ := r.Read(msg[:])
	return string(msg[:n])
}

func readAndParseState(r io.Reader) (state game.TurnState, err error) {
	json_bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return state, ServerError.Wrap(err)
	}
	err = json.Unmarshal(json_bytes, &state)
	if err != nil {
		return state, ServerError.Wrap(err)
	}
	return state, nil
}
