package hw06pipelineexecution

import (
	"fmt"
	"runtime"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	inner := make(Bi, 100)
	wg := sync.WaitGroup{}
	cnt := runtime.NumCPU()

	res := makePipe(inner, &stages)

	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go worker(done, res, out, &wg)
	}

	wg.Add(1)
	go produce(done, in, inner, &wg)

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func makePipe(in In, stages *[]Stage) Out {
	pipe := in
	for _, s := range *stages {
		pipe = s(pipe)
	}

	return pipe
}

func produce(done, in In, inner Bi, wg *sync.WaitGroup) {
	defer close(inner)
	defer wg.Done()

	for v := range in {
		select {
		case <-done:
			fmt.Println("coco1")
			return
		default:
		}

		select {
		case <-done:
			fmt.Println("coco2")

			return
		case inner <- v:
		}
	}
}

func worker(done, res In, out Bi, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range res {
		select {
		case <-done:
			return
		default:
		}

		select {
		case out <- v:
		case <-done:
			return
		}
	}
}
