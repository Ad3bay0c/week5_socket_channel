package explanation

import (
	"bufio"
	"fmt"
	"net"
)

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:1500")
	if err != nil {
		panic(err)
	}
	c := &client{
		nick: "unknown",
		conn: conn,
	}
	go c.receive()

	for {
		fmt.Print("type in  a message: ")
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}

		//s.message <- input
		_, err = conn.Write([]byte(input))

		if err != nil {
			fmt.Println(err)
			break
		}
		//
		//if n == 0 {
		//	continue
		//}

	}
}
