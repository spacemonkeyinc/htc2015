# Recover the Source

Hackerman is trying to hack his way back in time so he can ride a laser raptor.
Unfortunately, the original source code enabling him to do so has been
corrupted! Luckily, he backed it up using basic erasure encoding and was able
to recover some of the data with minimal loss.

The erasure code Hackerman was using was a (3, 5)-erasure code, where for every
3 bytes of real data that enter, 5 numbers exit. The first 3 numbers are
just the original bytes and the last 2 are parity numbers that have been
appended.

Hackerman's erasure code utilizes
"[polynomial interpolation](https://en.wikipedia.org/wiki/Polynomial_interpolation)",
but don't let that scare you! It only requires knowledge of algebra to use.

Perhaps you remember polynomials from math class (like `p(x) = 2x^2 + 3`).
For each 3 bytes of data, Hackerman is simply constructing a polynomial that
goes through 3 points representing those 3 bytes. Hackerman is using his
first 3 bytes of data as numbers, such that for byte 0,
`p(0) = <data of byte 0>`. Similarly for byte 1, `p(1) = <data of byte 1>`.
Here, the index of the byte is the `x` coordinate and the integer value of the
byte is the `y` coordinate. The first 3 bytes (indexes 0, 1, and 2), along with
the byte values themselves, give us the coordinates we need to create a
polynomial. We create a polynomial using the coordinates from the data, and
then we oversample it for the indexes 3 and 4 to get 2 parity numbers. Using
the parity numbers we can recreate missing or corrupt data.

## Example

Here is an example of 3 bytes of data that will be encoded. Using the indexes
and the integer values of each byte we plot it on a graph.

```
byte     | x | y
----------------
00000000 | 0 | 0 <-- data
00000010 | 1 | 2 <-- data
00000111 | 2 | 7 <-- data
```

![Data bytes plotted on a graph](/static/images/recover-polynomial.png)

We find a function that can generate these points on a graph (let's call
this function `p(x)` where `x` is the index of the byte). We use this function
to plot indexes `x = 3` and `x = 4`. Those values will be parity numbers.

```
byte     | x | y
----------------
00000000 | 0 | 0  <-- data
00000010 | 1 | 2  <-- data
00000111 | 2 | 7  <-- data
         | 3 | 15 <-- parity number
         | 4 | 26 <-- parity number
```

![Data bytes and oversampled parity numbers
plotted on a graph](/static/images/recover-polynomial-oversampled.png)

Now, with any 3 points of this graph, we can recreate the function `p(x)`
and use it to find the original byte values.

## Generating a polynomial from an arbitrary number of points.

There is a non-scary way of generating a polynomial that goes through every
point in a table of points called
[Lagrangian interpolation](http://mathworld.wolfram.com/LagrangeInterpolatingPolynomial.html).

Unfortunately, I haven't been able to find a
"Lagrange interpolation for beginners" link, but that doesn't mean it's
complicated! You might check out
[William Mueller's explanation](http://wmueller.com/precalculus/families/lagrange.html),
or [MatematicasVisuales' explanation](http://www.matematicasvisuales.com/english/html/analysis/interpolacion/lagrange.html)
(requires Flash), or even the overly detailed
[Wikipedia page on Polynomial interpolation](https://en.wikipedia.org/wiki/Polynomial_interpolation).

I think the key to making this less scary is to write out an entire Lagrange
interpolating polynomial in one line:

![Lagrange interpolation](/static/images/recover-wiki.png)

Here, the `x`s and `y`s with subscripts are your datapoints from your table.

The key observation is that we want `p(x_0) = y_0`, so to make that happen, we
want to multiply `y_0` by something that is 1 when `x = x_0` and multiply by 0
when `x` is any other data point we care about. We want the same property to
hold for `(x_1, y_1)`, etc. The construction above facilitates this property.
Try it out, the piece multiplied by `y_0` is 1 when `x = x_0` and is 0 when
`x` is `x_1, x_2, ...`.

So that's it, that's a Lagrangian interpolating polynomial. That's an equation
that goes through all the points you have, and provided that you only specify
it with the same amount of points every time, it will always simplify to the
same polynomial. Once you have `n` real data points, you construct this
polynomial with your `n` data points, and then you can generate parity data
points from it. You can also use `n` data points where some of the `n` are
parity to generate the original polynomial and therefore the original
coordinates.

You can do this! Good luck!

## Input

Hackerman used 2 parity numbers for every 3 real bytes, but you don't have to
do that. Each time your code is run, command line arguments will be passed for
`--in` indicating how many real bytes your encoding should use and `--out` to
indicate how many total numbers (real + parity) came out.

Your program will receive the same number of lines as the `--out` argument
separated by newlines. Each line will be a base-10 rational number (numerator
over denominator) of the point from the polynomial, or the string `MISSING`
used to indicate a missing or corrupt number. Up to `--out` - `--in`
numbers can be labeled as missing (so, 2 in the (3, 5) case).

### Example Input

```
./run --in=3 --out=5
```

```
0/1
MISSING
7/1
15/1
MISSING
```

## Output

Once you've divined the polynomial used for the block of data, you should
recreate the original bytes and return them (with each byte separated by
newlines).

### Example Output

For the previous example input, the example output should be:

```
00000000
00000010
00000111
```
