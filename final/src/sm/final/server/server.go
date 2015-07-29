// Copyright (C) 2015 Space Monkey, Inc.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/jtolds/go-oauth2http/utils"
	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/errors/errhttp"
	"github.com/spacemonkeygo/spacelog"

	"sm/final/game"
)

var (
	badRequestError = errors.NewClass("bad request",
		errhttp.SetStatusCode(http.StatusBadRequest))
	notFoundError = errors.NewClass("not found",
		errhttp.SetStatusCode(http.StatusNotFound))
	methodNotAllowedError = errors.NewClass("method not allowed",
		errhttp.SetStatusCode(http.StatusMethodNotAllowed))
	internalServerError = errors.NewClass("internal server error",
		errhttp.SetStatusCode(http.StatusInternalServerError))

	logger = spacelog.GetLogger()
)

const PlayerIdHeader = "X-SM-PlayerId"
const PlayerMonikerHeader = "X-SM-PlayerMoniker"

type Server struct {
	games *game.Games
}

func New(games *game.Games) *Server {
	return &Server{
		games: games,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Noticef(">>> %s %s", r.Method, r.URL)
	err := s.serveGame(w, r)
	logger.Noticef("<<< %s %s (%v)", r.Method, r.URL, err)
	if err != nil {
		code := errhttp.GetStatusCode(err, 500)
		desc := errhttp.GetErrorBody(err)
		http.Error(w, desc, code)
	}
}

func (s *Server) serveGame(w http.ResponseWriter, r *http.Request) (err error) {
	name, left := utils.Shift(r.URL.Path)
	action, left := utils.Shift(left)
	switch {
	case name == "":
		return notFoundError.New("%s", r.URL.Path)
	case action == "":
		return notFoundError.New("%s", r.URL.Path)
	case left != "":
		return notFoundError.New("%s", r.URL.Path)
	}
	switch r.Method {
	case "POST":
		var state game.TurnState

		command, err := game.CommandFromString(action)
		if err != nil {
			return badRequestError.Wrap(err)
		}

		player_id := r.Header.Get(PlayerIdHeader)

		if command == game.Join {
			moniker := r.Header.Get(PlayerMonikerHeader)
			if moniker == "" {
				return badRequestError.New("missing X-SM-PlayerMoniker header")
			}
			thegame, err := s.games.LookupOrCreate(name)
			if err != nil {
				return internalServerError.Wrap(err)
			}
			player_id, state, err = thegame.Join(moniker)
			if err != nil {
				return badRequestError.Wrap(err)
			}
		} else {
			thegame := s.games.Lookup(name)
			if thegame == nil {
				return badRequestError.New("game %s does not exist", name)
			}
			state, err = thegame.TakeTurn(player_id, command)
			if err != nil {
				return badRequestError.Wrap(err)
			}
		}

		// always return the player id header
		w.Header().Set(PlayerIdHeader, player_id)
		w.Header().Set("Content-Type", "application/json")
		state_bytes, err := json.MarshalIndent(state, "", "\t")
		if err != nil {
			return internalServerError.Wrap(err)
		}
		_, err = w.Write(state_bytes)
		logger.Errore(err)
	default:
		return methodNotAllowedError.New("%s", r.Method)
	}
	return nil
}
