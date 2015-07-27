# Contact!

Jodie Foster needs your help! She's receiving datasheets from outer space
again. They're definitely arriving as sides of a cube, but this time the
trouble is it seems like there's very few actually interesting cubes. The
vast majority of the cubes look like repeats.

## The problem

Your program will receive a set of unfolded cube descriptions in a terrible
format. Your job will be to parse all of these terrible cube descriptions and
then output only the first instance of each unique cube you find, in the format
and order you received it.

A cube face here is simply an `N`x`N` square of case-sensitive arbitrary
alpha-numeric characters. Here's an example of a cube face:

```
BAA
ABA
ABB
```

We can rotate this cube face clockwise 90 degrees, 180 degrees, and 270
degrees:

```
AAB  BBA  AAB
BBA  ABA  ABB
BAA  AAB  BAA
```

The above are all equivalent cube faces. If we were to take 6 of these faces
we could make a cube in three dimensions. Unfolded, such a cube might look like
this:

```
   BAA
   ABA
   ABB
BAABAABAABBA
ABAABAABAABA
ABBABBABBAAB
   BAA
   ABA
   ABB
```

With the individual faces separated this looks like:

```
    BAA
    ABA
    ABB

BAA BAA BAA BBA
ABA ABA ABA ABA
ABB ABB ABB AAB

    BAA
    ABA
    ABB
```

There's lots of ways to unfold any given cube. Your job is to understand the
relationships between the cube faces making up the cube and determine which
cubes (after refolding and rotations) you haven't seen before that the unfolded
cubes in the input describe.

## Example

### STDIN

```
   D
BCAC
A

    BCDB
    BBAB
  AABA
  AACA
BBDB
ABCC

CABC
BBAC
  DBABAD
  BBDDCA
      DD
      DD

AB
 CD
  AC

    DB
    AB
BACBAACC
BBBBBCBD
      AA
      AA
```

### STDOUT

```
   D
BCAC
A

    BCDB
    BBAB
  AABA
  AACA
BBDB
ABCC

CABC
BBAC
  DBABAD
  BBDDCA
      DD
      DD
```
