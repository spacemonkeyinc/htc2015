#!/usr/bin/python

import sys, json, argparse

sys.setrecursionlimit(10000)

parser = argparse.ArgumentParser()
parser.add_argument("--minimal-group-size")
args = parser.parse_args()
minimal_group_size = int(args.minimal_group_size)

class node(object):
  def __init__(self):
    self.neighbors = set()
    self.index = None
    self.lowlink = None
    self.onStack = False
  def add(self, neighbor):
    self.neighbors.add(neighbor)

G = {}
for a, b in json.loads(sys.stdin.read()):
  if a not in G: G[a] = node()
  if b not in G: G[b] = node()
  G[a].add(b)

index = 0
S = []

components = []

def strongconnect(v_name):
  global index
  v = G[v_name]
  v.index = index
  v.lowlink = index
  index += 1
  S.append(v_name)
  v.onStack = True

  for w_name in v.neighbors:
    w = G[w_name]
    if w.index is None:
      strongconnect(w_name)
      v.lowlink = min(v.lowlink, w.lowlink)
    elif w.onStack:
      v.lowlink = min(v.lowlink, w.index)

  if v.lowlink == v.index:
    component = []
    while True:
      w_name = S.pop()
      w = G[w_name]
      w.onStack = False
      component.append(w_name)
      if w_name == v_name:
        break
    if len(component) >= minimal_group_size:
      components.append(tuple(sorted(component)))


for v_name, v in G.iteritems():
  if v.index is None:
    strongconnect(v_name)

components.sort()
print json.dumps(components)
