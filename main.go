package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/client9/xson/hjson"
)

func main() {
	compact := flag.Bool("c", false, "compact")
	tojson := flag.Bool("j", false, "to JSON")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Need input")
	}
	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("error %s", err)
	}
	if *tojson {
		var v interface{}
		r := hjson.New(f)
		err = json.NewDecoder(r).Decode(&v)
		if err != nil {
			log.Fatalf("Unable to decode HJSON: %s", err)
		}
		var rawout []byte
		if *compact {
			rawout, err = json.Marshal(v)
		} else {
			rawout, err = json.MarshalIndent(v, "", "    ")
		}
		if err != nil {
			log.Fatalf("Unable to encode to JSON: %s", err)
		}
		os.Stdout.Write(rawout)
	}
}
