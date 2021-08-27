package socket

import (
	"fmt"
	"net"
	"os"
	"time"
)

func tcpClient() {
	tcpclient, err := net.DialTCP("tcp", nil, nil)
	checkError(err)
	_, err = tcpclient.Write([]byte("okay good"))
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v", err)
		os.Exit(1)
	}
}

func RunIt() {
	listen, err := net.Listen("tcp", ":8000")
	checkError(err)

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		checkError(err)

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	err := conn.SetDeadline(time.Now().Add(time.Minute))
	checkError(err)
	for {
		//reader := bufio.NewReader(conn)
		buff := make([]byte, 128)
		//s, err := reader.ReadString('\n')
		len, err := conn.Read(buff)
		checkError(err)

		if len == 0 {
			break
		}

		fmt.Println(string(buff))
	}

}

func TCPServer() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := tcpListen.Accept()
		checkError(err)

		go handle(conn)
	}
}


