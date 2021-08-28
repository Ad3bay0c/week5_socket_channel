package gorourines

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func BufferedChannels() {
	ch := make(chan interface{}, 2)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	ch <- 2
	go GetData(ch, wg)
	go ReadData(ch, wg)
	wg.Wait()
	//time.Sleep(3 * time.Second)
	fmt.Println("Done")
}

func GetData(ch chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Chat 2: \"%v\": \n", <-ch)
		//fmt.Println(<-ch)
		fmt.Print("Chat 1 : ")
		input, _ := inputReader.ReadString('\n')
		ch <- strings.TrimLeft(input, "\n")
	}
	//fmt.Println(<-ch)
	//ch <- "Good Morning"
}

func ReadData(ch chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Chat 1: \"%v\": \n", <-ch)
		fmt.Print("Chat 2 : ")
		input, _ := inputReader.ReadString('\n')
		ch <- strings.TrimLeft(input, "\n")
	}
	//fmt.Println(<-ch)

	//ch <- "Hope you're good"
}
