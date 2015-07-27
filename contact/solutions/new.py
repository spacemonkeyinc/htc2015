#!/usr/bin/python

import sys, random, string

if len(sys.argv) <= 1 or sys.argv[1] == "--help":
  print "usage: %s <side-width> [<alphabet>]" % sys.argv[0]
  sys.exit(1)

side = int(sys.argv[1])

def char():
  if len(sys.argv) < 3:
    alphabet = string.ascii_letters + string.digits
  else:
    alphabet = sys.argv[2]
  return random.choice(alphabet)

for _ in xrange(side):
  for _ in xrange(side * 3):
    sys.stdout.write(char())
  sys.stdout.write("\n")
for _ in xrange(side):
  for _ in xrange(side * 2):
    sys.stdout.write(" ")
  for _ in xrange(side * 3):
    sys.stdout.write(char())
  sys.stdout.write("\n")
