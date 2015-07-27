# Utah pentominoes

Utah is an unbelievable state. Not only does Utah contain 5 wonderful national
parks, enable fantastic skiing and mountain biking, host world-renowned
movie festivals, and offer a world-leading standard of living and
livability, it has a pretty interesting shape. Let's make a game out of
the Utah-shaped pentomino.

## The problem

The game board is a `n`x`n`x`n` *cube*. Two players take turns filling in one
of each of the cells in this cube, like some 3D Tic-Tac-Toe. But instead of
getting `n` cells in row, each player is trying to be the first to make a
Utah-shaped pentomino out of their filled cells. Any orientation, rotation, or
reflection is allowed, with the exception of diagonal (slanted) moves. A
diagonal Utah-pentomino is not an allowed shape.

Here are some examples of layouts a player might try to play for.

```
X__ ___ ___
XX_ ___ ___
XX_ ___ ___
```

Above is a cube. The left 3x3 square is the bottom of the cube, the middle 3x3
square is the middle layer of the cube, and the right 3x3 square is the top
layer of the cube. On the bottom layer, the Xs form a Utah-shaped pentomino.

Here is the cube in the sample image on this page:

```
___ _O_ ___
___ OO_ ___
XX_ XXO X__
```

In this cube, the Utah-shaped pentomino is along the front side of the cube.

Another one:

```
___ ___ ___
__X _XX _XX
___ ___ ___
```

This cube has an upside-down Utah-shaped pentomino down the center right of the
cube.

An example of a diagonal shape that does not contribute to winning:


```
X__ _X_ __X
X__ _X_ ___
___ ___ ___
```

Clear as mud?

To try and help, we've constructed a visualization program. When you clone the
project repo, you'll find `visualizer.py` in addition to the normal things.
(If for some reason this didn't work for you, you can
<a href="/static/files/visualizer.py">download it here</a>.)
You'll need `python` and `python-visual`
([http://vpython.org/](http://vpython.org/)) to be able to click around.
This visualizer is what generated the sample image on this page. Left-click
marks a box, right-click and drag rotates the view, and middle-click and drag
zooms.

We'll be giving you game states over `stdin`. Each game state will be exactly
one move away from either you or your opponent completing a Utah pentomino.
You will be next to play and need to make the optimum move. Your task is to
write a program that takes an arbitrary amount of these final-stage game
states and outputs the next game state after making your move. Your program
plays for the player placing `X`s (the other player is `O`).

## Example

### Input

```
./run
```

Input game states will be `NxNxN` cubes, newline separated, and `stdin` will be
closed when no more game states need to be sent.

```
X__ ___ ___
XX_ _O_ ___
OXO ___ ___

___ OOO ___
___ _OX ___
___ _X_ _X_
```

### Output

For output, just output the next game board state with your new move placed.

```
XX_ ___ ___
XX_ _O_ ___
OXO ___ ___

___ OOO ___
___ XOX ___
___ _X_ _X_
```
