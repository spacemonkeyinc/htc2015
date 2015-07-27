# The NSA needs friends

The NSA has had a rough two years and would like to find some new friends.

One thing the NSA has had to deal with is that friends can be fickle and
betray you and then flee. The NSA believes that perhaps finding friends from
strong friend groups might reduce the chance of future betrayal.

## The problem

You've been hired by the NSA to use its database of social network information
to help find some new friends. The NSA would like you to scan all of its
social networks and find groups of at least `N`, where inside the group,
everyone is friends with everyone else. For each friend group you find, it
should be the case that you can't add someone to the group to make it larger
without violating the constraint that everyone is friends with everyone else in
the group.

Your program is expected to take a `--minimal-group-size=N` argument, where `N`
will be the smallest allowed friend group size.

Your program should then read from `stdin` a JSON array of friendships, where
a friendship is a size-2 array of names.

You should then output a JSON array of arrays, where the internal arrays are
lists of names of friends in the valid groups you find.

## Example

### Input

```
./run --minimal-group-size=4
```

```
[["Edward Snowden", "Laura Poitras"],
 ["James Comey", "Keith Alexander"],
 ["James Comey", "Michael Hayden"],
 ["James Comey", "James Clapper"],
 ["Keith Alexander", "Michael Hayden"],
 ["Keith Alexander", "James Clapper"],
 ["Michael Hayden", "James Clapper"],
 ["Keith Alexander", "Jeff Wendling"]]
```

## Output

```
[["James Comey", "Keith Alexander", "Michael Hayden",
  "James Clapper"]]
```
