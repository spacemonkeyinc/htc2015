#!/usr/bin/python

import sys, json, argparse

parser = argparse.ArgumentParser()
parser.add_argument("--minimal-group-size")
args = parser.parse_args()
minimal_group_size = int(args.minimal_group_size)

n = {}
for a, b in json.loads(sys.stdin.read()):
  if a not in n: n[a] = set()
  if b not in n: n[b] = set()
  n[a].add(b)
  n[b].add(a)

answers = []
def bk(r, p, x):
  if not p and not x:
    if len(r) >= minimal_group_size:
      answers.append(tuple(sorted(list(r))))
  for v in p:
    bk(r | set([v]), p & n[v], x & n[v])
    p = p - set([v])
    x = x | set([v])


bk(set(), set(n.keys()), set())
answers.sort()
print json.dumps(answers)
