# cson

A [CSON](https://github.com/bevry/cson)
parser, with full unmarshalling support, in Go.

How? It takes an CSON input and converts it to JSON, and lets the
native Golang parser do the hard stuff.

### Differences and/or bugs:

 * Supports CoffeeScript multi-line comments
 * Multi-line strings may have some minor edge case bugs with whitespace