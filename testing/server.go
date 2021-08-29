package testing

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms		map[string]*room
	commands	chan command
}

func newServer() *server {
	return &server{
		rooms: 		make(map[string]*room),
		commands: 	make(chan command),
	}
}

func Execute() {
	s := newServer()

	go s.run()

	listener, err := net.Listen("tcp", ":1255")
	if err != nil {
		log.Fatalf("Unable to start server %v", err)
	}

	defer listener.Close()
	log.Println("Started server on :1255")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection")
			continue
		}

		go s.newClient(conn)
	}
}
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}
func (s *server) newClient(conn net.Conn) {
	log.Printf("New CLient is connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn: conn,
		nick: "anonymous",
		commands: s.commands,
	}

	c.readInput(s)
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg(fmt.Sprintf("Nick empty, enter a valid nick: /nick username"))
		return
	}
	c.nick = args[1]
	c.msg(fmt.Sprintf("All right i will call you %s", c.nick))
}
func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("enter the room you wanna enter")
		return
	}
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok{
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	c.msg(fmt.Sprintf("welcome to %s", r.name))
}
func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(rooms, ", ")))

}
func (s *server) msg(c *client, args []string) {
	if c.room == nil{
		c.err(errors.New("you must join the room first"))
		return
	}
	c.room.broadcast(c, c.nick +" "+ strings.Join(args[1:], " "))
}
func (s *server) quit(c *client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.msg("sad to see you go :(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}