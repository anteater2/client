package bitmesh

import (
	"errors"
	"log"
	"time"

	"github.com/anteater2/bitmesh/chord"
	"github.com/anteater2/bitmesh/chord/key"
	"github.com/anteater2/bitmesh/rpc"
)

var RPCCaller *rpc.Caller

var RPCGetKey rpc.RemoteFunc
var RPCPutKey rpc.RemoteFunc
var RPCFindSuccessor rpc.RemoteFunc
var RPCIsAlive rpc.RemoteFunc

type DHT struct {
	constr string
	maxkey key.Key
}

func Create(address string) DHT {
	table := DHT{
		constr: address + ":2001",
		maxkey: 1024,
	}
	if RPCCaller == nil {
		tmpC, err := rpc.NewCaller(2000)
		if err != nil {
			panic("Could not create RPCCaller!")
		}
		RPCCaller = tmpC
		RPCCaller.Start()
		RPCGetKey = RPCCaller.Declare("", chord.GetKeyResponse{}, 5*time.Second)
		RPCPutKey = RPCCaller.Declare(chord.PutKeyRequest{}, true, 5*time.Second)
		RPCFindSuccessor = RPCCaller.Declare(key.NewKey(1), chord.RemoteNode{}, 6*time.Second)
		RPCIsAlive = RPCCaller.Declare(true, true, 1*time.Second)
	}
	return table
}

// Put puts a key-value pair into dht.
func (dht *DHT) Put(k string, value string) error {
	hostInterf, err := RPCFindSuccessor(dht.constr, key.Hash(k, 1<<10))
	if err != nil { // I hate this language
		log.Fatal(err)
	}
	host := hostInterf.(chord.RemoteNode)
	rv, err2 := RPCPutKey(ConStr(&host), chord.PutKeyRequest{
		KeyString: k,
		Data:      []byte(value),
	})
	if err2 != nil { // I hate this language
		panic("IDK2")
	}
	if rv != true {
		return errors.New("insert failed")
	}
	return nil
}

// Get gets the value corresponding to the key from dht
func (dht *DHT) Get(k string) (string, error) {
	hostInterf, err := RPCFindSuccessor(dht.constr, key.Hash(k, 1<<10))
	if err != nil { // I hate this language
		panic("IDK3")
	}
	host := hostInterf.(chord.RemoteNode)
	rvInterf, err2 := RPCGetKey(ConStr(&host), k)
	if err2 != nil { // I hate this language
		panic("IDK4")
	}
	rv := rvInterf.(chord.GetKeyResponse)
	if rv.Error == false {
		panic("Could not get key!")
	}
	return string(rv.Data), nil
}

// Gets a connection string for a RemoteNode.
func ConStr(node *chord.RemoteNode) string {
	return node.Address + ":2001"
}
