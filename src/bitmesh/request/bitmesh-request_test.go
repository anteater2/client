package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetConnectionRequestString(testing *testing.T) {
	var duration int32 = 100
	expectedTimeout := int32(time.Now().Unix()) + duration
	expectedBlob := ConnectRequest{
		Name:       "rohits2",
		IP:         "127.0.0.1",
		Expiration: expectedTimeout,
	}
	var actualBlob ConnectRequest
	json.Unmarshal(GetConnectionRequestString("rohits2", "127.0.0.1", duration), &actualBlob)
	fmt.Println(string(GetConnectionRequestString("rohits2", "127.0.0.1", duration)))
	if actualBlob != expectedBlob {
		fmt.Println(actualBlob)
		fmt.Println(expectedBlob)
		testing.Fail()
	}
}

func TestGetLocalConnectionRequestString(testing *testing.T) {
	fmt.Println(GetLocalConnectionRequestString("rohits2", duration))
}

func TestSwarmPush(testing *testing.T) {
	var swarm SwarmInfo
	var duration int32 = 100
	expectedTimeout := int32(time.Now().Unix()) + duration
	for i := 0; i < 10; i++ {
		clientName := "peer" + strconv.Itoa(i)
		client := PeerInfo{
			Name:       clientName,
			IP:         "127.0.0.1",
			Expiration: expectedTimeout,
		}
		swarm.AddPeer(client)
		swarm.AddPeer(client) // Double push to test the swarm will reject duplicate resources
	}
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
}

func makeSwarm() SwarmInfo {
	var swarm SwarmInfo
	var duration int32 = 100
	expectedTimeout := int32(time.Now().Unix()) + duration
	for i := 0; i < 10; i++ {
		clientName := "peer" + strconv.Itoa(i)
		client := PeerInfo{
			Name:       clientName,
			IP:         "127.0.0.1",
			Expiration: expectedTimeout,
		}
		swarm.AddPeer(client)
	}
	return swarm
}

func TestSwarmRemoveNonexistent(testing *testing.T) {
	var swarm = makeSwarm()
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
	swarm.RemoveResource("peer100") //Test
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
}
func TestSwarmRemoveOne(testing *testing.T) {
	var swarm = makeSwarm()
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
	swarm.RemoveResource("peer1")
	if len(swarm.Peers()) != 9 {
		testing.Fail()
	}
}
func TestSwarmRemoveAll(testing *testing.T) {
	var swarm = makeSwarm()
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
	swarm.RemoveResource("peer1")
	if len(swarm.Peers()) != 9 {
		testing.Fail()
	}
	swarm.RemoveResource("peer2")
	if len(swarm.Peers()) != 8 {
		testing.Fail()
	}
	swarm.RemoveResource("peer3")
	swarm.RemoveResource("peer4")
	swarm.RemoveResource("peer5")
	swarm.RemoveResource("peer6")
	swarm.RemoveResource("peer7")
	swarm.RemoveResource("peer8")
	swarm.RemoveResource("peer9")
	swarm.RemoveResource("peer0")
	if len(swarm.Peers()) != 0 {
		testing.Fail()
	}
}

func TestSwarmUpdate(testing *testing.T) {
	var swarm = makeSwarm()
	if len(swarm.Peers()) != 10 {
		testing.Fail()
	}
	var duration = 100
	expectedTimeout := strconv.Itoa(int(time.Now().Unix()) + duration)
	jsonString := "[{\"name\": \"peer0\", \"ip\":\"192.168.0.1\", \"expiration\":" + expectedTimeout + "},{\"name\":\"testPeer\",\"ip\":\"0.0.0.0\",\"expiration\":" + expectedTimeout + "}]"
	swarm.UpdateSwarm([]byte(jsonString))
	_, info, _ := swarm.HasResource("peer0")
	if info.IP != "192.168.0.1" {
		testing.Fail()
	}
	if len(swarm.Peers()) != 11 {
		testing.Fail()
	}
}
