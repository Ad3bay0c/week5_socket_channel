package implementation

type instruction struct {
	id		int
	args	[]string
	user	*user
}

const (
	USERNAME = iota
	CHATS
	QUIT
	SEND
	JOIN
)