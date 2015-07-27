# nbd-lang

You're going to write a small programming language! Relax, it will be easy. It
will only support integers, addition/subtraction, ifzero statements, and
println.

## The problem

Here's the grammar for the language:

```
<program>     ::= <statements>
<statements>  ::= <statement>
                | <statement> ";" <statements>
<statement>   ::= <ifzero>
                | <assignment>
                | <addition> | <subtraction>
                | <println>
<ifzero>      ::= "?" <value> "{" <statements> "}"
<assignment>  ::= <variable> "=" <value>
<addition>    ::= <variable> "+=" <value>
<subtraction> ::= <variable> "-=" <value>
<println>     ::= "!" <value>
<value>       ::= <variable> | <integer>
<variable>    ::= [a-zA-Z][a-zA-Z0-9]*
<integer>     ::= [0-9]+
```

Optional whitespace is allowed between each term, but not required anywhere.

## Example

### STDIN

Here's an example program exercising all of the language:

```
x = 3;
y = 2;
y += x;
y -= 2;
x -= y;
? x {
  ! 1;
  y = 0;
  z = 0
};
? y {
  ! 2
};
? z {
  ! 3
}
```

### STDOUT

```
1
2
3
```
