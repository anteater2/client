package bitmesh

import (
	"errors"
	"log"
	"time"

	"github.com/anteater2/chord-node/key"
	"github.com/anteater2/chord-node/structs"
	"github.com/anteater2/rpc/rpc"
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
		RPCGetKey = RPCCaller.Declare("", structs.GetKeyResponse{}, 5*time.Second)
		RPCPutKey = RPCCaller.Declare(structs.PutKeyRequest{}, true, 5*time.Second)
		RPCFindSuccessor = RPCCaller.Declare(key.NewKey(1), structs.RemoteNode{}, 6*time.Second)
		RPCIsAlive = RPCCaller.Declare(true, true, 1*time.Second)
	}
	return table
}

// Put puts a key-value pair into dht.
func (dht *DHT) Put(k string, value string) error {
	hostInterf, err := RPCFindSuccessor(dht.constr, key.Hash(k, 1024))
	if err != nil { // I hate this language
		log.Fatal(err)
	}
	host := hostInterf.(structs.RemoteNode)
	rv, err2 := RPCPutKey(ConStr(&host), structs.PutKeyRequest{
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
func (dht *DHT) Get(key string) (string, error) {
	hostInterf, err := RPCFindSuccessor(dht.constr, key)
	if err != nil { // I hate this language
		panic("IDK3")
	}
	host := hostInterf.(structs.RemoteNode)
	rvInterf, err2 := RPCGetKey(ConStr(&host), key)
	if err2 != nil { // I hate this language
		panic("IDK4")
	}
	rv := rvInterf.(structs.GetKeyResponse)
	if rv.Error {
		panic("Could not get key!")
	}
	return string(rv.Data), nil
}

// Gets a connection string for a structs.RemoteNode.
func ConStr(node *structs.RemoteNode) string {
	return node.Address + ":2001"
}
