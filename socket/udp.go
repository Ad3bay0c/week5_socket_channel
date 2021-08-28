package socket

import (
	"fmt"
	"net"
)

func UDPServer() {
	service := ":2000"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	startUdp, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	_, err = startUdp.Write([]byte("What's up"))
	checkError(err)

	server, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	buf := make([]byte, 128)
	_, err = server.Read(buf[0:])
	checkError(err)
	fmt.Println(string(buf))
}