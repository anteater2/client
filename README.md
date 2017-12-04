# Bitmesh Client-Side Library
This is the bitmesh-client side library, which allows users to make a DHT object by giving the address of a Bitmesh chord ring.  The DHT object has Get() and Put() defined, allowing for string:string lookups and puts.
## Get
Get calls FindSuccessor() on the ring, and then calls GetKey() on the node responsible.  Failure to find the key currently causes a panic.
## Put 
Put calls FindSuccessor() on the ring, and then calls PutKey() on the node responsible.
## Testing
A dockerfile is set up to run ```go test``` on ```connectionlib_test.go```.  The ```go test``` is run from inside a script (```docker-test```) because docker does not like entrypoints with arguments.

The purpose of the dockerfile is to run the tests on the same VLAN as the chord ring, and stop the RPC library from dying when IP addresses change.