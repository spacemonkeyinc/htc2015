#!/usr/bin/python

import sys

# from http://aspn.activestate.com/ASPN/Cookbook/Python/Recipe/252178
def all_perms(str):
  if len(str) <= 1:
    yield str
  else:
    for perm in all_perms(str[1:]):
      for i in range(len(perm)+1):
        yield perm[:i] + str[0:1] + perm[i:]


SHAPES = [[(1, 0, 0), (1, 0, 0), (0, 1, 0), (-1, 0, 0)]]
SHAPE_TREE = {}
for shape in SHAPES:
  for i in (1, -1):
    for j in (1, -1):
      for k in (1, -1):
        for p in all_perms((0, 1, 2)):
          next_shape = []
          for step in shape:
            next_shape.append((i * step[p[0]], j * step[p[1]], k * step[p[2]]))
          tree_section = SHAPE_TREE
          for step in next_shape:
            if not tree_section.has_key(step):
              tree_section[step] = {}
            tree_section = tree_section[step]
OPPONENTS = ["red"]
MARKER = "blue"

class Cell(object): pass
class OutOfBounds(object): pass

class Board(object):

  def __init__(self, cube_size=3):
    self.cube_size = cube_size
    self.places = [[[Cell() for k in range(self.cube_size)]
                    for i in range(self.cube_size)]
                   for j in range(self.cube_size)]
    for i in range(self.cube_size):
      for j in range(self.cube_size):
        for k in range(self.cube_size):
          self.places[i][j][k].board_pos = (i, j, k)
          self.places[i][j][k].marker = " "

  def __getitem__(self, pos):
    for var in pos:
      if var < 0 or var >= self.cube_size:
        return OutOfBounds()
    i, j, k = pos
    if self.places[i][j][k].marker == " ":
      return None
    else:
      return self.places[i][j][k].marker

  def __setitem__(self, pos, marker):
    for var in pos:
      if var < 0 or var >= self.cube_size:
        raise Exception("out of bounds: %r" % (pos,))
    i, j, k = pos
    if self.places[i][j][k].marker != " ":
      raise Exception("non-empty position: %r" % (pos,))
    self.places[i][j][k].marker = marker

def _count_empty_spots(board, marker, pos, tree_portion, empty_spots, depth):
  if board[pos] is not None and board[pos] != marker:
    return []
  if board[pos] is None:
    empty_spots = empty_spots + [pos]
  if len(tree_portion) == 0:
    if len(empty_spots) == 0:
      raise Exception("My understanding of game rules differs from "+
                      "the driver apparently. This game should be over.")
    else:
      return [empty_spots]
  possible_shape_lists = []
  least_remaining_spots = float('inf')
  for step in tree_portion.keys():
    i, j, k = pos
    i += step[0]
    j += step[1]
    k += step[2]
    rv = _count_empty_spots(
        board, marker, (i, j, k), tree_portion[step], empty_spots, depth + 1)
    if not rv:
      continue
    remaining_spots = len(rv[0])
    if remaining_spots > least_remaining_spots:
      continue
    if remaining_spots < least_remaining_spots:
      least_remaining_spots = remaining_spots
      possible_shape_lists = []
    possible_shape_lists.extend(rv)
  return possible_shape_lists

def _find_best_possible_win(board, marker):
  best_empty_spot_lists = []
  for i in range(board.cube_size):
    for j in range(board.cube_size):
      for k in range(board.cube_size):
        empty_spot_lists = _count_empty_spots(
            board, marker, (i, j, k), SHAPE_TREE, [], 0)
        if not empty_spot_lists:
          continue
        if (best_empty_spot_lists and
            len(empty_spot_lists[0]) > len(best_empty_spot_lists[0])):
          continue
        if (not best_empty_spot_lists or
            len(empty_spot_lists[0]) < len(best_empty_spot_lists[0])):
          best_empty_spot_lists = []
        best_empty_spot_lists.extend(empty_spot_lists)
  return best_empty_spot_lists

def get_next_position(board):
  best_spot_lists = _find_best_possible_win(board, MARKER)
  for opponent in OPPONENTS:
    defense_spot_lists = _find_best_possible_win(board, opponent)
    if not defense_spot_lists:
      continue
    if not best_spot_lists or (
        len(defense_spot_lists[0]) < len(best_spot_lists[0])):
      best_spot_lists = defense_spot_lists
  best_spots = set()
  for spot_list in best_spot_lists:
    for spot in spot_list:
      best_spots.add(spot)
  if len(best_spots) != 1:
    if len(best_spots) < 10:
      raise Exception("unexpected board state: %r" % best_spots)
    raise Exception("unexpected board state (too many possible options)")
  return list(best_spots)[0]

def getOne(input):
  seen_cube = False
  for line in input:
    line = line.strip()
    if not line:
      if seen_cube:
        return
      continue
    seen_cube = True
    yield line

def readCube(input):
  lines = []
  for line in getOne(input):
    lines.append(line)
  if not lines:
    return None
  return lines

def parseCube(lines):
  board = Board(cube_size=len(lines))
  for i, line in enumerate(lines):
    j = 0
    k = 0
    for char in line:
      if char != " ":
        if char != "_":
          board[(i, j, k)] = {"X": "blue", "O": "red"}[char]
        j += 1
      else:
        j = 0
        k += 1
  return board

def printCube(board):
  for i in range(board.cube_size):
    for k in range(board.cube_size):
      for j in range(board.cube_size):
        sys.stdout.write({"blue": "X", "red": "O", None: "_"}[board[(i, j, k)]])
      if k < board.cube_size - 1:
        sys.stdout.write(" ")
    sys.stdout.write("\n")
  sys.stdout.write("\n")

def runCube(board):
  board[get_next_position(board)] = MARKER
  printCube(board)

def main():
  while True:
    lines = readCube(sys.stdin)
    if lines is None:
      return
    runCube(parseCube(lines))

if __name__ == "__main__":
  main()
