# hjson
A HJSON (http://hjson.org) parser and unmarshaller written in Go

A scant 200 lines implements a [HJSON](http://hjson.org) parser, with full unmarshalling
support, in Go.

How? It takes an HJSON input and converts it to JSON, and lets the
native Golang parser do the hard stuff.

### Working:

 * Initial root object '{' ... '}' is not required
 * Keys do not need to be quoted
 * Values do not need to be quoted
 * Commas are optional
 * Comments script-style `#`

### Not-Working:

Mostly due to laziness

 * Comments double-slash  `//`
 * Comments `/* ... */`
 * Multi-line strings

### Differences and/or bugs:

 * Unquoted strings have trailing whitespace removed.
 * Likely to be some encoding issues, please file bugs.