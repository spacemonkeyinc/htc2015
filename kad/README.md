# Don't be a Kad

The typical way to measure distance between two numbers on a number line is
to subtract them and take the absolute value. This is also referred to as the
"Euclidean distance" in one dimension. So, the distance between 5 and 8? 3.

This isn't the only way to measure distance. Just as there's multiple ways to
measure distance in two dimensions ("Manhattan distance", "as the crow flies"
(Euclidean in 2 dimensions), etc.), there's multiple "metrics" in one dimension
as well.

One way to compute a distance with a different metric is to simply XOR the two
numbers together. For example, the XOR-distance between 5 and 8 is 13.

## The problem

You will be given a long list of newline-separated numbers over stdin. Then you
will be given a blank line, and then a series of numbers representing queries
that should be answered with the `N` XOR-nearest numbers from the initial long
list (where `N` is provided by the `--request_size` commandline argument). Each
answer should be output on a new line. Numbers won't be 2^64 or larger or less
than zero.

## Example

### Input

```
./run --request_size=3
```

```
139596472382
332835991912435768
20381213957485674
58347458
73819384887567
2

123451234512345
3125
```

### Output

```
73819384887567
58347458
2
2
58347458
139596472382
```
