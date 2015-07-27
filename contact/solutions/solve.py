#!/usr/bin/python

import sys, fractions, random

UP = 0
RIGHT = 1
DOWN = 2
LEFT = 3

UNFOLDINGS = [
        {UP: {}, DOWN: {}, RIGHT: {RIGHT: {RIGHT: {}}}},
        {RIGHT: {UP: {}, DOWN: {}, RIGHT: {RIGHT: {}}}},

        {DOWN: {}, RIGHT: {UP: {}, RIGHT: {RIGHT: {}}}},
        {DOWN: {}, RIGHT: {RIGHT: {UP: {}, RIGHT: {}}}},
        {DOWN: {}, RIGHT: {RIGHT: {RIGHT: {UP: {}}}}},
        {RIGHT: {UP: {UP: {}, RIGHT: {RIGHT: {}}}}},
        {RIGHT: {DOWN: {}, RIGHT: {UP: {}, RIGHT: {}}}},
        {RIGHT: {DOWN: {RIGHT: {DOWN: {RIGHT: {}}}}}},
        {RIGHT: {DOWN: {RIGHT: {DOWN: {}, RIGHT: {}}}}},
        {RIGHT: {RIGHT: {DOWN: {RIGHT: {RIGHT: {}}}}}},
        {RIGHT: {DOWN: {RIGHT: {RIGHT: {DOWN: {}}}}}},

        {UP: {}, RIGHT: {DOWN: {}, RIGHT: {RIGHT: {}}}},
        {UP: {}, RIGHT: {RIGHT: {DOWN: {}, RIGHT: {}}}},
        {UP: {}, RIGHT: {RIGHT: {RIGHT: {DOWN: {}}}}},
        {RIGHT: {DOWN: {DOWN: {}, RIGHT: {RIGHT: {}}}}},
        {RIGHT: {UP: {}, RIGHT: {DOWN: {}, RIGHT: {}}}},
        {RIGHT: {UP: {RIGHT: {UP: {RIGHT: {}}}}}},
        {RIGHT: {UP: {RIGHT: {UP: {}, RIGHT: {}}}}},
        {RIGHT: {RIGHT: {UP: {RIGHT: {RIGHT: {}}}}}},
        {RIGHT: {UP: {RIGHT: {RIGHT: {UP: {}}}}}}]

def getOne(input):
  seen_cube = False
  for line in input:
    line = line.rstrip()
    if not line:
      if seen_cube:
        return
      continue
    seen_cube = True
    yield line

def getNextNeighbor(face, first, second):
  neighbor = face["neighbors"][first]
  if neighbor is None:
    return None, None
  rotation = face["rotations"][first]
  second = (4 + second - rotation) % 4
  return (neighbor["neighbors"][second],
          (neighbor["rotations"][second] + rotation) % 4)

def setNeighbor(face, dir):
  if face["neighbors"][dir] is not None:
    return True
  cw, cw_rotation = getNextNeighbor(face, (dir + 3) % 4, dir)
  if cw is not None:
    face["neighbors"][dir] = cw
    face["rotations"][dir] = (cw_rotation + 1) % 4
    return True
  ccw, ccw_rotation = getNextNeighbor(face, (dir + 1) % 4, dir)
  if ccw is not None:
    face["neighbors"][dir] = ccw
    face["rotations"][dir] = (ccw_rotation + 3) % 4
    return True
  return False

def readCube(input):
  lines = []
  max_length = 0
  for line in getOne(input):
    lines.append(line)
    if len(line) > max_length:
      max_length = len(line)
  if not lines:
    return None, None
  if len(lines) > max_length:
    vertical = True
    grid_width = max_length
    grid_height = len(lines)
  else:
    vertical = False
    grid_width = len(lines)
    grid_height = max_length

  side = fractions.gcd(grid_height, grid_width)
  grid = [[None for x in xrange(grid_width/side)]
          for y in xrange(grid_height/side)]
  face_count = 0
  faces = {}
  for m, line in enumerate(lines):
    for n, cell in enumerate(line):
      if cell.isspace():
        continue
      if vertical:
        row = m
        col = n
      else:
        row = n
        col = grid_width - m - 1
      if grid[row/side][col/side] is None:
        face = {
            "id": face_count,
            "neighbors": [None, None, None, None],
            "rotations": [0, 0, 0, 0],
            "desc": [[None for x in xrange(side)] for y in xrange(side)]}
        faces[face_count] = face
        face_count += 1
        grid[row/side][col/side] = face
      grid[row/side][col/side]["desc"][row%side][col%side] = cell

  for m, row in enumerate(grid):
    for n, col in enumerate(row):
      if col is None:
        continue
      if m - 1 >= 0:
        other = grid[m-1][n]
        if other is not None:
          col["neighbors"][UP] = other
          col["rotations"][UP] = 0
          other["neighbors"][DOWN] = col
          other["rotations"][DOWN] = 0
      if n - 1 >= 0:
        other = grid[m][n-1]
        if other is not None:
          col["neighbors"][LEFT] = other
          col["rotations"][LEFT] = 0
          other["neighbors"][RIGHT] = col
          other["rotations"][RIGHT] = 0

  all_good = False
  while not all_good:
    all_good = True
    for face in faces.itervalues():
      if not setNeighbor(face, LEFT):
        all_good = False
      if not setNeighbor(face, RIGHT):
        all_good = False
      if not setNeighbor(face, UP):
        all_good = False
      if not setNeighbor(face, DOWN):
        all_good = False

  return faces.values(), lines

def rotate(grid, rotations):
  for _ in xrange(rotations):
    new_grid = [[None for _ in xrange(len(grid))]
                for _ in xrange(len(grid[0]))]
    for m, row in enumerate(grid):
      for n, cell in enumerate(row):
        new_grid[n][len(grid) - m - 1] = cell
    grid = new_grid
  return grid

def faceDescEqual(desc1, desc2):
  if len(desc1) != len(desc2):
    return False
  for i in xrange(len(desc1)):
    for j in xrange(len(desc1)):
      if desc1[i][j] != desc2[i][j]:
        return False
  return True

def formatCube(face, initial_dir, unfolding, output_rotation):
  sparse_grid = {}
  def fill(face, m, n, dir, next_steps):
    sparse_grid[(m, n)] = (face, dir)
    for next_step, steps_after in next_steps.iteritems():
      if next_step == UP:
        next_m, next_n = m - 1, n
      elif next_step == RIGHT:
        next_m, next_n = m, n + 1
      elif next_step == DOWN:
        next_m, next_n = m + 1, n
      elif next_step == LEFT:
        next_m, next_n = m, n - 1
      new_dir = (4 + next_step - dir) % 4
      next_face = face["neighbors"][new_dir]
      next_rotation = (face["rotations"][new_dir] + dir) % 4
      fill(next_face, next_m, next_n, next_rotation, steps_after)

  fill(face, 0, 0, initial_dir, unfolding)

  min_m, min_n = 0, 0
  max_m, max_n = 0, 0
  for m, n in sparse_grid.iterkeys():
    min_m, min_n = min(min_m, m), min(min_n, n)
    max_m, max_n = max(max_m, m), max(max_n, n)

  grid = [[None for n in xrange((max_n - min_n + 1) * len(face["desc"]))]
          for m in xrange((max_m - min_m + 1) * len(face["desc"]))]
  for (m, n), (face, dir) in sparse_grid.iteritems():
    face_base_m = (m - min_m) * len(face["desc"])
    face_base_n = (n - min_n) * len(face["desc"])
    face_desc = face["desc"]
    face_desc = rotate(face_desc, dir)
    for face_m, row in enumerate(face_desc):
      for face_n, cell in enumerate(row):
        grid[face_base_m + face_m][face_base_n + face_n] = cell

  grid = rotate(grid, output_rotation)

  output = ""
  for row in grid:
    for cell in row:
      if cell is None:
        output += " "
      else:
        output += cell
    output += "\n"
  output += "\n"
  return output

def faceEqual(face1, face2):
  for dir in [UP, RIGHT, LEFT, DOWN]:
    if not faceDescEqual(face1["desc"], rotate(face2["desc"], dir)):
      continue
    if (formatCube(face1, 0, UNFOLDINGS[0], 0) ==
        formatCube(face2, dir, UNFOLDINGS[0], 0)):
      return True
  return False

def cubeEqual(cube1, cube2):
  for face in cube1:
    if faceEqual(face, cube2[0]):
      return True
  return False

def main():
  cubes = []
  while True:
    cube, lines = readCube(sys.stdin)
    if cube is None:
      break
    exists = False
    for other in cubes:
      if cubeEqual(cube, other):
        exists = True
        break
    if not exists:
      for line in lines:
        print line
      print
      cubes.append(cube)

if __name__ == "__main__": main()
