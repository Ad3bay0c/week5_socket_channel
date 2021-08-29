package implementation

type server struct {
	groups			map[string]*group
	instructions	chan instruction
}


