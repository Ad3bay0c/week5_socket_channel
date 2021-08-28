package gorourines

import (
	"fmt"
	"runtime"
	"time"
)

//var data chan<- int : means data can only be sent to channel
//var data <-chan int : means data can only be received from channel

func Example() {
	var c = make(chan int)

	go source(c)
	go suck(c)
	time.Sleep(3 * time.Second)
}

func source(c chan<- int) {
	for i :=1; i<=10; i++ {
		c <- i
	}
	defer close(c)
	//c <- 24
	//c <- 25
}
func suck(c <-chan int) {
	for v := range c  {
		fmt.Println(v)
	}

	//for   {
	//	fmt.Println(<-c)
	//}
	//fmt.Println(<- c)
	//fmt.Println(<- c)
}

func SelectTesting(){
	runtime.GOMAXPROCS(2)
	c1 := make(chan int)
	c2 := make(chan int)

	go pump1(c1)
	go pump2(c2)

	go suck1(c1, c2)
	time.Sleep(1e9)
}
func pump1(c chan<- int) {
	for i:=2; ; i++ {
		c <- i * 2
	}
}
func pump2(c chan<- int) {
	for i := 2; ; i++ {
		c <- i * 5
	}
}
func suck1(c <-chan int, c2 <-chan int) {
	for {
		select {
		case v := <- c:
			fmt.Println("Receiving Data from Pump 1: ", v)
		case v := <- c2:
			fmt.Println("Receiving Data from Pump 2: ", v)
		}
	}

}