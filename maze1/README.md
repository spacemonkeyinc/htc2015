# Bowie's in a maze

Imagine you're in your house one night, minding your own business, when
suddenly your little brother disappears. Right then, David Bowie shows up and
wisks you off to a gigantic maze. He tells you if you don't solve the maze as
fast as possible, your brother will turn into a muppet.

It's cool, you tell yourself. This isn't actually 1986. You have Google Maps.
So you load up your phone and get a map.

## The problem

Given a map of a maze in the following format, find the shortest route from
start to finish in the least amount of time possible. If it helps, every test
case you will be given will have exactly one possible solution description.

### Input format

Maps are simply grids given to you as newline-delimited rows consisting of four
different types of characters. A cell can either be `S`, `F`, `X`, or `-`, and
a row ends with a `\n`.

* `S` means start, this is your starting cell. There will only be one of these.
* `F` means finish, this is where you are trying to get to. There will also
  only be one of these.
* `-` is a navigable cell. You can move to it if it is horizontally or
  vertically adjacent to your current cell (no diagonal moves).
* `X` is a filled cell, you can't move here.

The map will be delivered over standard in, and standard in will be closed when
the map is over.

### Example map

```
----
SXX-
XF--
```

## Output format

Once you've figured out the shortest path, your program needs to tell someone.
Your program should use a string of characters representing steps for your
path, newline separated. A step can either be `N`, `S`, `E`, `W`.

* `N` means north, go up on the map representation.
* `S` means south, go down.
* `E` means east, go right.
* `W` means west, go left.

### Example output

For the previous example map, the example output should be:

```
N
E
E
E
S
S
W
W
```
