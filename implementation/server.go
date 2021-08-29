package implementation

import (
	"errors"
	"log"
	"net"
	"strings"
)

type server struct {
	groups			map[string]*group
	instructions	chan instruction
}

func StartServer() {
	s := &server{
		instructions: make(chan instruction),
	}

	listener, err := net.Listen("tcp", ":1300")
	checkError(err, "Error Listening to Port")

	log.Println("Server started at localhost:1300")

	for {
		conn, err := listener.Accept()
		checkError(err, "Error Accepting Request: " + err.Error())

		go s.handleRequest(conn)
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Printf("ERR: %s\n", msg)
	}
}

func (s *server) handleRequest(conn net.Conn) {
	log.Printf("New User COnnected; %v", conn.RemoteAddr().String())

	newUser := &user{
		username:    "unknown",
		conn:        conn,
		instructions: make(chan instruction),
	}

	go newUser.acceptInput()
}

func (s *server) readClient() {
	for ins := range s.instructions{
		switch ins.id {
		case USERNAME:
			s.addUsername(ins.user, ins.args)
		case JOIN:
		case REPLY:
			s.replyMessage(ins.user, ins.args)
		case GROUPLIST:
		case QUIT:
		}
	}
}

func (s *server) addUsername(user *user, args []string) {
	if len(args) < 2 {
		user.errorMessage(errors.New("specify the Username; *name Bay"))
	}
	user.username = strings.TrimSpace(args[1])
	user.writeMessage("Username Updated: " + user.username)
}

func (s *server) replyMessage(user *user, args []string) {
	if len(args) < 2 {
		user.errorMessage(errors.New("type a message; *reply hi"))
	}

	if user.group == nil {
		user.errorMessage(errors.New("please, Join a Group first; *join groupname"))
	}

	user.group.message(args)
}


