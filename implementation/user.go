package implementation

import "net"

type user struct {
	username		string
	conn			net.Conn
	instruction		chan<- instruction
}
