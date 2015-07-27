#!/usr/bin/python2.5

import visual, sys

class OutOfBounds(Exception): pass

class Board(object):
  def __init__(self, cube_size=3, empty_color="white"):
    self.cube_size = cube_size
    self.empty_color = empty_color
    visual.scene.ambient = .6
    self.places = \
        [[[visual.box(
              pos=(i*2 - self.cube_size + 1,
                   k*2 - self.cube_size + 1,
                   j*2 - self.cube_size + 1),
              size=(1, 1, 1),
              color=getattr(visual.color, self.empty_color))
           for k in range(self.cube_size)]
          for i in range(self.cube_size)]
         for j in range(self.cube_size)]
    for i in range(self.cube_size):
      for j in range(self.cube_size):
        for k in range(self.cube_size):
          self.places[i][j][k].board_pos = (i, j, k)
          self.places[i][j][k].color_name = self.empty_color

  def __getitem__(self, pos):
    for var in pos:
      if var < 0 or var >= self.cube_size:
        raise OutOfBounds()
    i, j, k = pos
    if self.places[i][j][k].color_name == self.empty_color:
      return None
    return self.places[i][j][k].color_name

  def __setitem__(self, pos, color):
    for var in pos:
      if var < 0 or var >= self.cube_size:
        raise OutOfBounds("out of bounds")
    i, j, k = pos
    self.places[i][j][k].color_name = color
    self.places[i][j][k].color = getattr(visual.color, color)

  def get_next_mouseclick_coords(self):
    visual.scene.mouse.events = 0
    while True:
      visual.rate(30)
      if visual.scene.mouse.clicked:
        pick =  visual.scene.mouse.getclick().pick
        if pick.__class__ == visual.box:
          return pick.board_pos

class Player(object):
  def __init__(self, color):
    self.name = color
    self.color = color

  def get_next_position(self, board):
    return board.get_next_mouseclick_coords()

class Driver(object):
  def __init__(self, players, board):
    self.players = players
    self.board = board

  def run(self):
    while True:
      for player in self.players:
        pos = player.get_next_position(self.board)
        self.board[pos] = player.color
        self.print_layout()

  def print_layout(self):
    char = {"blue": "X", "red": "O", None: '_'}
    for i in range(self.board.cube_size):
      for k in range(self.board.cube_size):
        for j in range(self.board.cube_size):
          sys.stdout.write(char[self.board[(i, j, k)]])
        sys.stdout.write(" ")
      sys.stdout.write("\n")
    sys.stdout.write("\n")

if __name__ == "__main__":
  Driver([Player("blue"), Player("red")], Board(cube_size=3)).run()
