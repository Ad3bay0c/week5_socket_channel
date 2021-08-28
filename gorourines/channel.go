package gorourines

import (
	"fmt"
	"time"
)

func ChannelsTesting() {
	fmt.Println("Perform some kind of Channels Testing")

	ch := make(chan interface{})

	go readData(ch)
	go readData(ch)

	_, _ = <- ch, <-ch

	go getData(ch)

	time.Sleep(time.Second)
	fmt.Println("End of Main function")
}

func readData(ch chan interface{}) {
	ch <- "Good morning"
	ch <- 24
	ch <- false
	ch <- "Goodbye"
}

func getData(ch chan interface{}) {
	var input interface{}

	for {
		input = <- ch

		fmt.Println(input)
	}
	close(ch)
}