#!/usr/bin/python

import random, json, base64

edges = set()

for disjoint in xrange(10):
  vertices = []

  for i in xrange(20000):
    if vertices and random.choice([True, False, False]):
      v1 = random.choice(vertices)
      v2 = random.choice(vertices)
      if v1 != v2:
        edges.add((v1, v2))
    else:
      vertices.append(
          base64.b64encode(str(disjoint) + "-" + str(i)).rstrip("="))

result = []
for v1, v2 in edges:
  edge = [v1, v2]
  random.shuffle(edge)
  result.append(edge)
random.shuffle(result)

print json.dumps(result)
