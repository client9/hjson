# hjson
A HJSON (http://hjson.org) parser and unmarshaller written in Go

A scant 200 lines implements a [HJSON](http://hjson.org) parser, with full unmarshalling
support, in Go.

How? It takes an HJSON input and converts it to JSON, and lets the
native Golang parser do the hard stuff.

### Differences and/or bugs:

 * Unquoted strings have trailing whitespace removed.
 * Strings can be quoted using single quote, e.g. `'foo'`.
 * Multi-line strings can also be started with double quote, e.g. `"""`
 * Likely to be some encoding issues, please file bugs.
 * Likely edge cases with multi-line strings, please file bugs