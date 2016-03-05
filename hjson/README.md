# hjson
A HJSON (http://hjson.org) parser and unmarshaller written in Go

[![Build Status](https://travis-ci.org/client9/hjson.svg?branch=master)](https://travis-ci.org/client9/hjson) [![Go Report Card](http://goreportcard.com/badge/client9/hjson)](http://goreportcard.com/report/client9/hjson) [![GoDoc](https://godoc.org/github.com/client9/hjson?status.svg)](https://godoc.org/github.com/client9/hjson) [![Coverage](http://gocover.io/_badge/github.com/client9/hjson)](http://gocover.io/github.com/client9/hjson) [![license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/client9/hjson/master/LICENSE)

A scant 200-lines implements a HJSON parser, with full unmarshalling support, in Go.

How? It takes an HJSON input and converts it to JSON, and let's the
native Golang parser do the hard stuff.

### Working:

 * Initial root object '{' ... '}' is not required
 * Keys do not need to be quoted
 * Values do not need to be quoted
 * Commas are optional

### Not-Working:

Mostly due to laziness

 * Commennts script-style  `#`
 * Comments double-slash  `//`
 * Comments `/* ... */`
 * Multi-line strings

### Differences and/or bugs:

 * Unquoted strings have trailing whitespace removed.
 * Likely to be some encoding issues