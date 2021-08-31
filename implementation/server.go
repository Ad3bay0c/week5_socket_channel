package implementation

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	chats        map[string]*chat
	instructions chan instruction
}

func StartServer() {
	s := &server{
		chats:        make(map[string]*chat),
		instructions: make(chan instruction),
	}

	go s.readClient()

	listener, err := net.Listen("tcp", ":1300")
	checkError(err, "Error Listening to Port")
	defer listener.Close()
	log.Println("Server started at localhost:1300")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection")
			continue
		}

		go s.handleRequest(conn)
	}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Printf("ERR: %s\n", msg)
	}
}

func (s *server) handleRequest(conn net.Conn) {
	log.Printf("New User Connected; %v", conn.RemoteAddr().String())

	newUser := &user{
		username:    "unknown",
		conn:        conn,
	}

	newUser.acceptInput(s)
}

func (s *server) readClient() {
	for ins := range s.instructions{
		switch ins.id {
		case USERNAME:
			s.addUsername(ins.user, ins.args)
		case JOIN:
			s.joinGroup(ins.user, ins.args)
		case SEND:
			s.replyMessage(ins.user, ins.args)
		case CHATS:
			s.groupList(ins.user)
		case QUIT:
			s.quitConnection(ins.user)
		}
	}
}

func (s *server) addUsername(user *user, args []string) {
	if len(args) < 2 {
		user.errorMessage(errors.New("specify the Username; *name Bay"))
		return
	}
	user.username = strings.TrimSpace(args[1])
	user.writeMessage(user, "Username Updated: " + user.username)
}

func (s *server) replyMessage(user *user, args []string) {
	if len(args) < 2 {
		user.errorMessage(errors.New("type a message; *reply hi"))
		return
	}

	if user.chat == nil {
		user.errorMessage(errors.New("please, Join a Group first; *join groupName"))
		return
	}

	msg := strings.Join(args[1:], " ")
	user.chat.message(user, msg)
}

func (s *server) groupList(user *user) {
	if len(s.chats) == 0 {
		user.conn.Write([]byte(fmt.Sprintln("$ Empty Group")))
		return
	}
	group := ""
	for _, v := range s.chats {
		group += v.name + ", "
	}

	user.conn.Write([]byte(fmt.Sprintf("$ Groups are: %v\n",group)))
}

func (s *server) quitGroup(user *user)  {
	if user.chat == nil {
		user.errorMessage(errors.New("please join a chat"))
		return
	}
	user.chat.message(user, fmt.Sprintf("%v left the chat", user.username))
	delete(user.chat.members, user.conn.RemoteAddr())

}

func (s *server) quitConnection(user *user) {
	if user.chat != nil {
		s.quitGroup(user)
	}
	user.conn.Close()
	log.Printf("A connection Disconnected: %v", user.conn.RemoteAddr())
}

func (s *server) joinGroup(u *user, args []string) {
	if len(args) < 2 {
		u.errorMessage(errors.New("enter the chat to join"))
		return
	}
	name := strings.TrimSpace(args[1])

	g, ok := s.chats[name]
	if !ok {
		g = &chat{
			name: name,
			members: make(map[net.Addr]*user),
		}

		s.chats[name] = g
	}

	g.members[u.conn.RemoteAddr()] = u

	u.chat = g
	u.chat.message(u,fmt.Sprintf("%v joined the chat", u.username))

	u.writeMessage(u, "Welcome to "+ name)
}