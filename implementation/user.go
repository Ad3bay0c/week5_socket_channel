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
	conn 			net.Conn
	chat 			*chat
}

func (user *user) acceptInput(s *server) {
	user.writeMessage(user,"Type a Message... \n")
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
		case "*chats":
			s.instructions <- instruction{
				id: CHATS,
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
		case "*send":
			s.instructions <- instruction{
				id: SEND,
				args: args,
				user: user,
			}
		default:
			user.errorMessage(
				errors.New(
					"invalid command; Choose from these commands: \n\t *username 'your name'-> to set your username," +
						"\n\t *chats -> to list all available chat chats," +
						"\n\t *join 'chats'-> to join a chat," +
						"\n\t *send 'message'-> to send a message to a chat," +
						"\n\t *quit -> to disconnect from a connection",
					))
		}
	}
}
func (user *user) writeMessage(u *user, msg string) {
	user.conn.Write([]byte(fmt.Sprintf("$%v: %s\n",u.username, msg)))
}

func (user *user) errorMessage(err error) {
	user.conn.Write([]byte(fmt.Sprintf("OOPS!!! An Error Occurred: \n%v\n", err)))
}