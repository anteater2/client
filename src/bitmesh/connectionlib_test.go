package bitmesh

import (
	"log"
	"testing"
)

func TestPutGetSimple(*testing.T) {
	tbl := Create("172.0.0.2")
	log.Printf("Created DHT! Putting value...")
	tbl.Put("hello", "world")
	log.Printf("Put value! Getting value...")
	rv, _ := tbl.Get("hello")
	if rv != "world" {
		panic("uh oh")
	}
	log.Printf("Got value!")
}
