package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

const countWorkers = 10
const countJobs = 100000
const countErrors = 10
const chanSize = 0

func startWorker(workerNum int, in <-chan func() error, chError chan error, quit chan bool) {
	println("start worker")
	println(workerNum)
	for input := range in {
		wg.Add(1)
		select {
		case <-quit:
			wg.Done()
			println("stop worker")
			return
		default:
			println("start job")
			chError <- input()
		}
		wg.Done()
		println("stop job")
	}

}
func startError(in chan func() error, chError chan error, quit chan bool) {
	var count int
	for err := range chError {
		fmt.Println(err)
		count++
		if count == countErrors {
			println("close channels")
			close(quit)
			break
		}
	}
	go func() {
		wg.Wait()
		close(chError)
	}()
	fmt.Println("after max error")
	for err := range chError {
		fmt.Println(err)
	}

	fmt.Println("all done")
}

func main() {
	runtime.GOMAXPROCS(0)
	worketInput := make(chan func() error, chanSize)
	chError := make(chan error, 1)
	quit := make(chan bool)

	var sliceWork []func() error
	for i := 1; i < countJobs; i++ {
		work := work(i)
		sliceWork = append(sliceWork, work)
	}

	for z := 0; z < countWorkers; z++ {
		go startWorker(z, worketInput, chError, quit)
	}

	go startError(worketInput, chError, quit)

	for _, f := range sliceWork {
		select {
		case <-quit:
			return
		default:
			println("add job")
			worketInput <- f
		}

	}
	close(worketInput)

}

func work(x int) func() error {
	return func() error {
		for i := 1; i > 0; i++ {
			if i%x*1000 == 0 {
				time.Sleep(1000 * time.Millisecond)
				return fmt.Errorf("fail work: %d", x)
			}
		}
		return nil
	}
}
