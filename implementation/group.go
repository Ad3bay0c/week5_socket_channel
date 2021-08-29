package implementation

import "net"

type group struct {
	name 	string
	members	map[net.Addr]*user
}
