# Final round

![final-round-tanks](https://raw.githubusercontent.com/SpaceMonkeyInc/htc2015/master/final/screenshot.png)

The documentation for this game is here:
https://github.com/SpaceMonkeyInc/htc2015/tree/master/final/final/docs.md

The game is a Go program, but since it requires asset building we've provided
a Makefile.

You can build everything with `make` (possibly after redownloading the sound 
files listed in `final/final/sounds/CREDITS`).

To run, first launch the server (`bin/final-server`), then launch two bots 
(`bin/circle-bot`, `bin/battery-bot`).
