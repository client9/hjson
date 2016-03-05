# cson

A [CSON](https://github.com/bevry/cson)
parser, with full unmarshalling support, in Go.

How? It takes an CSON input and converts it to JSON, and let's the
native Golang parser do the hard stuff.

### Not-Working:

 * Multi-line strings with triple-single and triple-double quotes.

### Differences and/or bugs:

 * Supports CoffeeScript multi-line comments