// Copyright (C) 2015 Space Monkey, Inc.

package main

import (
	"flag"
	"net/http"

	"github.com/jtolds/go-oauth2http/utils"
	"github.com/russross/blackfriday"
	"github.com/spacemonkeygo/spacelog"

	"sm/codecomp/setup/general"
	"sm/final/assets"
	"sm/final/game"
	"sm/final/renderer/sdl"
	"sm/final/server"
)

var (
	endpoint        = flag.String("http.endpoint", ":8080", "host:port to serve from")
	backgroundMusic = flag.Bool("bgmusic", true, "if false, no background music")
	staticPath      = flag.String("http.static", "", "path to file serving")
	logger          = spacelog.GetLogger()
)

func main() { general.Run(Main) }

func docs(w http.ResponseWriter, r *http.Request) {
	doc, err := assets.Asset("final/docs.md")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(header))
	w.Write(blackfriday.MarkdownCommon([]byte(doc)))
	w.Write([]byte(footer))
}

func Main() error {
	logger.Noticef("listening at %q", *endpoint)

	sdl.Run(func() {
		if *backgroundMusic {
			go func() {
				err := sdl.BackgroundMusic("final/sounds/music1.wav")
				if err != nil {
					logger.Noticef("assets missing background music: %v", err)
				}
			}()
		}

		go func() {
			mux := utils.DirMux{
				"game": server.New(game.NewGames(game.DefaultConfig())),
				"":     http.HandlerFunc(docs)}
			if *staticPath != "" {
				mux["static"] = http.FileServer(http.Dir(*staticPath))
			}
			panic(http.ListenAndServe(*endpoint, mux))
		}()
	})

	return nil
}

var (
	header = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Future Tank - hackthe.computer</title>
    <link rel="stylesheet" href="https://hackthe.computer/static/css/bootstrap.css?1">
    <link rel="stylesheet" href="https://hackthe.computer/static/css/bootstrap-theme.css?1">
    <!--[if lt IE 9]>
      <script src="//oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js">
        </script>
      <script src="//oss.maxcdn.com/respond/1.4.2/respond.min.js">
        </script>
    <![endif]-->
    <link rel="stylesheet" href="https://hackthe.computer/static/css/site.css?1" />
  </head>
  <body>
    <div class="container main">`
	footer = `    </div>

    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js">
      </script>
    <script src="https://hackthe.computer/static/js/bootstrap.js?1"></script>
    <script src="https://hackthe.computer/static/js/site.js?1"></script>
  </body>
</html>`
)
