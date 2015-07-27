# The NSA needs more friends

You successfully helped the NSA find friends! Congratulations!

Emboldened by your success, the NSA has decided to try and widen its net even
more and has asked you to make some changes to your friend-finding program.

## The problem

Instead of using a social network database where friend relationships are
commutative and undirected, they would like you to use a social network
database with directed relationships, where users "follow" other users instead.
In this case, just because user A follows user B, user B doesn't necessarily
follow user A.

They still want you to find good groups of friends, but the requirement is
relaxed; not every user needs to be following every other user in the group
directly. Instead, it is sufficient if every user is following every other user
in the group transitively and indirectly. To be specific, if user A is
following user B, and user B is following user C, then we can say user A is
indirectly following user C. There's no limit to this chain of transitive
relationships, but every user in a group must have a chain of directed
transitive "following" relationships to every other user in the group.

## Input format

Your program is still expected to take a `--minimal-group-size=N` argument,
where `N` will be the smallest allowed amount of people in a group.

Your program should then read from `stdin` a JSON array of following
relationships, where a following relationship is a size-2 array of names,
the first element being the follower, and the second element being the
followee.

### Example

```
./run --minimal-group-size=4
```

STDIN:

```
[["Edward Snowden", "Laura Poitras"],
 ["James Comey", "Keith Alexander"],
 ["Keith Alexander", "Michael Hayden"],
 ["Michael Hayden", "James Clapper"],
 ["James Clapper", "James Comey"],
 ["Keith Alexander", "Jeff Wendling"]]
```

## Output format

You should then output a JSON array of arrays, where the internal arrays are
lists of names of friends that are all friends with each other. The internal
arrays should be at least of length `N`.

### Example

STDOUT:

```
[["James Comey", "Keith Alexander", "Michael Hayden",
  "James Clapper"]]
```
