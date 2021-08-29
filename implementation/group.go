package implementation

import (
	"net"
	"strings"
)

type group struct {
	name 	string
	members	map[net.Addr]*user
}

func (g *group) message(user *user, msg []string) {
	message := strings.Join(msg, " ")

	for key, u := range g.members {
		if user.conn.RemoteAddr() != key {
			u.writeMessage(message)
		}
	}
}
