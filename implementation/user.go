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
	instructions	chan<- instruction
}

func (user *user) acceptInput() {
	user.writeMessage("Type a writeMessage...")
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
			user.instructions <- instruction{
				id: USERNAME,
				args: args,
				user: user,
			}
		case "*grouplist":
			user.instructions <- instruction{
				id: GROUPLIST,
				args: args,
				user: user,
			}
		case "*join":
			user.instructions <- instruction{
				id: JOIN,
				args: args,
				user: user,
			}
		case "*quit":
			user.instructions <- instruction{
				id: QUIT,
				args: args,
				user: user,
			}
		case "*reply":
			user.instructions <- instruction{
				id: REPLY,
				args: args,
				user: user,
			}
		default:
			user.errorMessage(errors.New("invalid command"))
		}
	}
}
func (user *user) writeMessage(msg string) {
	user.conn.Write([]byte(fmt.Sprintf("%v: %s",user.username, msg)))
}

func (user *user) errorMessage(err error) {
	user.conn.Write([]byte(fmt.Sprintf("OOPS!!! An Error Occurred: %v", err)))
}