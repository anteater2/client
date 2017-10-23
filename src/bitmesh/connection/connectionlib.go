package connection

import (
	"bitmesh/request"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// GetConnectionRequestString accepts a resource name, IP, and duration and creates
// a JSON blob that will register that resource and IP for the specified duration.
func RegisterToTracker(endpoint string, name string, duration int32) []byte {
	//http.Post(endpoint, "text/json", &buf)
	request.GetLocalConnectionRequestString(name, duration)
	var timeout = int32(time.Now().Unix()) + Timeout
	connection := ConnectRequest{
		Name:       "rohits2",
		IP:         "127.0.0.1",
		Expiration: timeout,
	}
	c, _ := json.Marshal(connection) // TODO: Correct error handling
	return c
}

// GetIPAddr gets a local IP of this machine.
// Note: This may not be an internet IP.
func GetIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Could not open network interface!")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
