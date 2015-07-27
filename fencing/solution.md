See [8.3.9 The grammar for path data](http://www.w3.org/TR/SVG/paths.html#PathDataBNF)
for the BNF to parse `<path>` elements.

See [pnpoly](http://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html#Point on an Edge)
for detecting whether a point is inside a polygon.  This test doesn't handle
points exactly on the polygon's perimeter.  Sometime it returns true for
points on the perimeter.  Other times it returns false.

See a [python solution](http://stackoverflow.com/questions/328107/how-can-you-determine-a-point-is-between-two-other-points-on-a-line-segment)
for how to detect if a point is on a line segment.  Iterate through all
edges of a polygon to see if a point lines on the perimeter.
