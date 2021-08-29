package implementation

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	groups			map[string]*group
	instructions	chan instruction
}
func newServer() *server {
	return &server{
		groups: make(map[string]*group),
		instructions: make(chan instruction),
	}
}
func StartServer() {
	s := newServer()

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
		case REPLY:
			s.replyMessage(ins.user, ins.args)
		case GROUPLIST:
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

	if user.group == nil {
		user.errorMessage(errors.New("please, Join a Group first; *join groupName"))
		return
	}

	msg := strings.Join(args[1:], " ")
	user.group.message(user, msg)
}

func (s *server) groupList(user *user) {
	if len(s.groups) == 0 {
		user.conn.Write([]byte(fmt.Sprintln("$ Empty Group")))
		return
	}
	group := ""
	for _, v := range s.groups {
		group += v.name + ", "
	}

	user.conn.Write([]byte(fmt.Sprintf("$ Groups are: %v\n",group)))
}

func (s *server) quitGroup(user *user)  {
	if user.group == nil {
		user.errorMessage(errors.New("please join a group"))
		return
	}
	user.group.message(user, fmt.Sprintf("%v left the group", user.username))
	delete(user.group.members, user.conn.RemoteAddr())

}

func (s *server) quitConnection(user *user) {
	if user.group != nil {
		s.quitGroup(user)
	}
	user.conn.Close()
	log.Printf("A connection Disconnected: %v", user.conn.RemoteAddr())
}

func (s *server) joinGroup(u *user, args []string) {
	if len(args) < 2 {
		u.errorMessage(errors.New("enter the group to join"))
		return
	}
	name := strings.TrimSpace(args[1])

	g, ok := s.groups[name]
	if !ok {
		g = &group{
			name: name,
			members: make(map[net.Addr]*user),
		}

		s.groups[name] = g
	}

	g.members[u.conn.RemoteAddr()] = u

	u.group = g
	u.group.message(u,fmt.Sprintf("%v joined the group", u.username))

	u.writeMessage(u, "Welcome to "+ name)
}