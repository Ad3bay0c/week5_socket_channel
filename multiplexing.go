package main

type Request struct {
	a, b	int
	replyC	chan int
}

type binOp func(a, b int) int

func Client(op binOp, R *Request) {
	R.replyC <- op(R.a, R.b)
}

func Server(op binOp, service chan *Request) {
	for {
		req := <- service

		go Client(op, req)
	}
}

func StartServer(op binOp) chan *Request {
	reqChan := make(chan *Request)
	go Server(op, reqChan)
	return reqChan
}