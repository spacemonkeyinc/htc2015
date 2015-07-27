#!/usr/bin/python

import sys, fractions, random
from solve import readCube, rotate, formatCube
from solve import UP, DOWN, LEFT, RIGHT, UNFOLDINGS

if len(sys.argv) <= 1 or sys.argv[1] == "--help":
  print "usage: %s <number-of-copies>" % sys.argv[0]
  sys.exit(1)

copies = int(sys.argv[1])
faces, _ = readCube(sys.stdin)
for i in xrange(copies):
  face = random.choice(faces)
  unfolding = random.choice(UNFOLDINGS)
  initial_dir = random.choice([UP, RIGHT, LEFT, DOWN])
  output_rotation = random.choice([UP, RIGHT, LEFT, DOWN])
  sys.stdout.write(formatCube(face, initial_dir, unfolding, output_rotation))
