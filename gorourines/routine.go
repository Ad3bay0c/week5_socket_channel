package gorourines

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
)

func Basic() {
	fmt.Println(runtime.NumCPU(), runtime.NumGoroutine())
	fmt.Println("Beginning of main")
	//go longwait()
	//time.Sleep(time.Second)
	fmt.Println("end of main function")
	fmt.Println(runtime.NumCPU(), runtime.NumGoroutine())
}
func longwait(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Beginning of longwait")

}
func WaitGroup() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go longwait(wg)

	wg.Wait()

	fmt.Println("All finished")
	numCores := flag.Int("n", 2, "Number of Cores")
	flag.Parse()
	fmt.Println(*numCores)
}