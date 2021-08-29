package implementation

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type user struct {
	username		string
	conn			net.Conn
	group			*group
}

func (user *user) acceptInput(s *server) {
	user.writeMessage(user,"Type a writeMessage... \n")
	for {
		input, err := bufio.NewReader(user.conn).ReadString('\n')
		if err != nil {
			user.errorMessage(err)
			continue
		}
		if len(input) < 1 {
			user.conn.Write([]byte("Enter a valid command"))
			continue
		}

		input = strings.Trim(input, "\n")
		args := strings.Split(input, " ")
		ins := strings.ToLower(strings.TrimSpace(args[0]))

		switch ins {
		case "*username":
			s.instructions <- instruction{
				id: USERNAME,
				args: args,
				user: user,
			}
		case "*grouplist":
			s.instructions <- instruction{
				id: GROUPLIST,
				args: args,
				user: user,
			}
		case "*join":
			s.instructions <- instruction{
				id: JOIN,
				args: args,
				user: user,
			}
		case "*quit":
			s.instructions <- instruction{
				id: QUIT,
				args: args,
				user: user,
			}
		case "*reply":
			s.instructions <- instruction{
				id: REPLY,
				args: args,
				user: user,
			}
		default:
			user.errorMessage(errors.New("invalid command"))
		}
	}
}
func (user *user) writeMessage(u *user, msg string) {
	user.conn.Write([]byte(fmt.Sprintf("$%v: %s\n",u.username, msg)))
}

func (user *user) errorMessage(err error) {
	user.conn.Write([]byte(fmt.Sprintf("OOPS!!! An Error Occurred: %v\n", err)))
}