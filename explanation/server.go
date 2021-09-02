package explanation

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// server running
 // server accept different request

 type server struct {
 	clients		map[net.Addr]*client
 	message		chan string
 }

 type client struct {
 	nick	string
 	conn	net.Conn
 }

 var s = &server{
 	clients: make(map[net.Addr]*client),
 	message: make(chan string),
 }

 func StartServer() {
 	go acceptInput()

 	listener, err := net.Listen("tcp", ":1500")
 	if err!= nil {
 		panic(err)
	}
	log.Println("Server Started at localhost: 1500")
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf(err.Error())
			continue
		}

		go handleRequest(conn)
	}
 }

 func handleRequest(conn net.Conn) {
 	log.Printf("New Connection Made: %v",conn.RemoteAddr())
 	client := &client{
		nick: "unknown",
		conn: conn,
	}

	s.clients[client.conn.RemoteAddr()] = client

	//for {
	//	buf := make([]byte, 1024)
	//
	//	_, err := conn.Read(buf)
	//
	//	//log.Println(string(buf))
	//	//
	//	if err != nil {
	//		fmt.Println("Error Accepting Message")
	//		break
	//	}
	//
	//	s.message <- string(buf)
	//	//fmt.Println(string(buf))
	//
	//	//for _, client := range s.clients{
	//	//	client.conn.Write([]byte(fmt.Sprintf("%s: %s", client.nick, buf)))
	//	//}
	//}

	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			//type a message
			continue
		}

		//s.message <- input

		for _, client := range s.clients{
			client.conn.Write([]byte(fmt.Sprintf("%s: %v", client.nick, input)))
		}

	}

 }

 func acceptInput() {
 	for v := range s.message {
 		for _, client := range s.clients{
 			client.conn.Write([]byte(fmt.Sprintf("%s: %v", client.nick, v)))
		}
	}
 }

func (c *client) receive()  {
	for {
		buf := make([]byte, 1024)

		n, err := c.conn.Read(buf)

		//log.Println(string(buf))
		//
		if err != nil {
			fmt.Println("Error Accepting Message")
			c.conn.Close()
			break
		}
		if n > 0 {
			fmt.Println(string(buf))
			c.conn.Write(buf)
		}

	}
}