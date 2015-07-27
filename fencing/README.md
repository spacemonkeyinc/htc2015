# Fencing

Now that the NSA monitors all written and electronic communications and
can break any encryption, people only whisper their secrets to each other directly.
Using a 2 trillion dollar budget, the NSA has built a device like Arthur C. Clarke's
[wormholes](https://en.wikipedia.org/wiki/The_Light_of_Other_Days)
or Philip K. Dick's [time scoop](https://en.wikipedia.org/wiki/Paycheck_%28short_story%29)
that can watch and listen to events at any location on earth.

Nonetheless, due in part to some massive whistleblowing, the NSA has agreed
that they will only use their new technology at specific court-authorized GPS
locations inside suspected enemy compounds.

The courts need help determining which GPS coordinates are actually inside
enemy compounds.

## The problem

Your job is to evaluate a series of points and tell the NSA if each point is
inside (true) or outside (false) an enemy compound.

The courts are big supporters of open standards and send you
[SVG](http://www.w3.org/TR/SVG11/) maps. The maps use `<path>` elements rather
than `<polygon>` elements and the perimeter isn't drawn consistently clockwise
or counterclockwise. The perimeters are always simple but not necessarily
convex polygons. Enemy compounds are polygons: they have straight lines between
vertices and do not contain holes. The enemy compound elements have
`id="compound"`. There may be other elements in the SVG document that are not
part of the enemy compound. There will only be one enemy compound per map.

Read an SVG perimeter map (a single `<svg>` element), a blank line, and
a list of (x, y) points. For each point, print `true` if the point is inside or
directly on the perimeter and `false` if the point is outside the perimeter.

## Example

![Enemy perimeter](/static/images/fencing.png)

* *light blue* - enemy compound
* *green* - point in enemy compound (`true`, ok to surveil)
* *red* - point not in enemy compound (`false`, not ok to surveil)

### Input

```
<svg version="1.1" xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink">
  <path id="compound" d="M14,9.722L67.891,23.929L109,9.564L109,58L67.894,37.383L14,58Z"/>
  <path id="kindergarten" d="m60,10l10,0,l0,10l-10,0z"/>
</svg>

65.895, 46.402
32.06, 57.681
84.692, 56.929
62.135, 71.214
119.278, 39.635
77.925, 13.319
6.496, 33.62
39.579, 12.568
120.782, 9.56
20.031, 50.914
17.775, 17.079
29.805, 34.372
65.144, 32.869
102.737, 20.086
102.737, 48.658
89.204, 35.124
```

### Output

```
false
false
false
false
false
false
false
false
false
true
true
true
true
true
true
true
```
