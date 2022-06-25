package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/xEasy/sanji"
)

var luffy *sanji.Luffy

func main() {
	luffy = sanji.New("01G6CVYHR8P9HQ3JKG7212ABPG", "default", "localhost", 7775)

	if os.Getenv("MODE") == "W" {
		runLuffy()
	} else {
		// produceJob()
		produceOneJob()
	}
}

func runLuffy() {
	luffy.Process("nqueue", func(ctx *sanji.Context) {
		data := string(ctx.Job.Data)
		fmt.Println("get job", ctx.Job.ID, " from queue ", ctx.Job.Queue, "data: ", data)
	}, 10, 100, true)

	luffy.Process("pqueue", func(ctx *sanji.Context) {
		data := string(ctx.Job.Data)
		fmt.Println("get job", ctx.Job.ID, " from queue ", ctx.Job.Queue, "data: ", data)
	}, 10, 100, true)

	luffy.Run()
}

func produceOneJob() {
	luffy.Enqueue("nqueue", []byte(fmt.Sprintf("nqueue job %d", 1)), 0)
}

func produceJob() {
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		for i := 0; i < 10; i++ {
			luffy.Enqueue("nqueue", []byte(fmt.Sprintf("nqueue job %d", i)), 0)
		}
		for i := 0; i < 10; i++ {
			luffy.EnqueueWithPriority("pqueue", []byte(fmt.Sprintf("priority mid job %d", i)), sanji.PriorityLevelMid, 0)
		}
		wait.Done()
	}()
	go func() {
		for i := 0; i < 10; i++ {
			luffy.EnqueueWithPriority("pqueue", []byte(fmt.Sprintf("priority hi job %d", i)), sanji.PriorityLevelHight, 0)
		}
		wait.Done()
	}()
	wait.Wait()
}
