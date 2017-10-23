package request

import (
	"encoding/json"
	"fmt"
	"time"
	"bitmesh/connection"
)

// ConnectRequest is a JSON struct encoding for a connection request.
type ConnectRequest struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Expiration int32  `json:"expiration"`
}

// DisconnectRequest is a JSON struct encoding for a disconnection request.
type DisconnectRequest struct {
	Name string `json:"name"`
}

// PeerInfo is a JSON struct that contains information about a peer node.
type PeerInfo struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Expiration int32  `json:"expiration"`
}

// SwarmInfo is a json struct that contains information about the swarm as a whole (i.e. many peers)
type SwarmInfo struct {
	peers []PeerInfo
}

// GetConnectionRequestString accepts a resource name, IP, and duration and creates
// a JSON blob that will register that resource and IP for the specified duration.
func GetConnectionRequestString(name string, ip string, duration int32) []byte {
	var timeout = int32(time.Now().Unix()) + duration
	connection := ConnectRequest{
		Name:       name,
		IP:         ip,
		Expiration: timeout,
	}
	c, _ := json.Marshal(connection) // TODO: Correct error handling
	return c
}

// GetLocalConnectionRequestString accepts a resource name, IP, and duration and creates
// a JSON blob that will register that resource and IP for the specified duration.
func GetLocalConnectionRequestString(name string, duration int32) []byte {
	var timeout = int32(time.Now().Unix()) + duration
	connection.
	connection := ConnectRequest{
		Name:       name,
		IP:         ip,
		Expiration: timeout,
	}
	c, _ := json.Marshal(connection) // TODO: Correct error handling
	return c
}

// Prune the swarm to remove expired entries
func (sw SwarmInfo) Prune() {
	var timeout = int32(time.Now().Unix())
	for i := range sw.peers {
		if sw.peers[i].Expiration <= timeout {
			sw.peers = append(sw.peers[:i], sw.peers[i+1:]...)
		}
	}
}

// AddPeer adds a peer to the swarm info struct.
func (sw *SwarmInfo) AddPeer(p PeerInfo) {
	peerExists, _, index := sw.HasResource(p.Name)
	if peerExists {
		sw.removeIndex(index)
	}
	sw.peers = append(sw.peers, p)
}

// RemovePeer removes a peer from the swarm info struct.
func (sw *SwarmInfo) removeIndex(i int) {
	sw.peers = append(sw.peers[:i], sw.peers[i+1:]...)
}

// RemovePeer removes a peer from the swarm info struct.
// Returns true if a peer was removed, else false.
func (sw *SwarmInfo) RemovePeer(p PeerInfo) bool {
	sw.Prune()
	for i := range sw.peers {
		if sw.peers[i] == p {
			sw.peers = append(sw.peers[:i], sw.peers[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveResource removes a resource from the swarm info struct.
// Returns true if a peer was removed, else false.
func (sw *SwarmInfo) RemoveResource(resource string) bool {
	for i := range sw.peers {
		if sw.peers[i].Name == resource {
			sw.peers = append(sw.peers[:i], sw.peers[i+1:]...)
			return true
		}
	}
	return false
}

// Peers returns a copy of the backing peer list of the swarm.
// NOTE: THIS METHOD CANNOT MUTATE THE SWARM
func (sw SwarmInfo) Peers() []PeerInfo {
	return sw.peers
}

// HasResource checks the SwarmInfo for a resource (a PeerInfo name).
// Returns (true, PeerInfo, index) if the resource exists and (false, <undefined>, -1) if not.
// NOTE: THIS METHOD CANNOT MUTATE THE SWARM
func (sw SwarmInfo) HasResource(resource string) (bool, PeerInfo, int) {
	for i, element := range sw.peers {
		if element.Name == resource {
			return true, element, i
		}
	}
	var junk PeerInfo
	return false, junk, -1
}

// UpdateSwarm updates the SwarmInfo from a JSON blob of peers.
func (sw *SwarmInfo) UpdateSwarm(jsonBlob []byte) {
	var newPeers []PeerInfo
	err := json.Unmarshal(jsonBlob, &newPeers)
	if err != nil {
		fmt.Println("error:", err) // TODO: Real error handling
	}
	for _, peer := range newPeers {
		sw.AddPeer(peer)
	}
	sw.Prune()
}
