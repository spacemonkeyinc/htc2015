#!/usr/bin/python

import sys, random

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

def readCube(input):
  lines = []
  max_length = 0
  for line in getOne(input):
    lines.append(line)
    if len(line) > max_length:
      max_length = len(line)
  if not lines:
    return None
  return lines

cubes = []
while True:
  cube = readCube(sys.stdin)
  if cube is None:
    break
  cubes.append(cube)
random.shuffle(cubes)

for cube in cubes:
  for line in cube:
    print line
  print
