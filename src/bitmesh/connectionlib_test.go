package bitmesh

import (
	"log"
	"testing"
)

func TestPutGetSimple(*testing.T) {
	tbl := Create("172.0.0.1")
	log.Printf("Created DHT! Putting value...")
	tbl.Put("hello", "world")
	log.Printf("Put value! Getting value...")
	rv, _ := tbl.Get("hello")
	if rv != "world" {
		panic("FUCK THIS LANGUAGE JESUS")
	}
	log.Printf("Got value!")
}
