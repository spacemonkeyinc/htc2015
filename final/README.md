# Final round

![final-round-tanks](https://raw.githubusercontent.com/SpaceMonkeyInc/htc2015/master/final/screenshot.png)

The documentation for this game is here:
https://github.com/SpaceMonkeyInc/htc2015/tree/master/final/final/docs.md

The game is a Go program. You'll first have to create the assets file
(`Makefile` in `final/src/sm/final/assets/`), possibly after redownloading the
sound files (listed in `final/final/sounds/CREDITS`), then you can build
everything with:

```
GOPATH=/path/to/htc2015/final go install sm/...
```
