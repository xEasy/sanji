## Sanji

Sanji is a Job queue base on the [Lmstfy](https://github.com/bitleak/lmstfy)

#### Usage

look at [example](https://github.com/xEasy/sanji/tree/main/example)

```golang
  // make a new processor instance, get token by lmstfy admin api.
	chef = sanji.New("01G6CVYHR8P9HQ3JKG7212ABPG", "default", "localhost", 7775)

  // [Consumer]
  // add job processing middleware
  chef.Use(JobFunc)

  // func (s *Luffy) Process(queue string, job JobFunc, ttr uint32, concurrency int, is_priority bool)
  // queue: the job queue
  // job: how to processing job
  // ttr: job processing ttr
  // concurrency: max concurrency worker
  // is_priority: need priority queue feature?
	chef.Process("nqueue", func(ctx *sanji.Context) {
		data := string(ctx.Job.Data)
		fmt.Println("get job", ctx.Job.ID, " from queue ", ctx.Job.Queue, "data: ", data)
	}, 10, 100, true)
  
  // [producer]
  // func (s *Luffy) Enqueue(queue string, data []byte, delay uint32) (jobID string, err error)
  // queue: the job processing queue
  // data: job data
  // delay: whether to delay processing, 0 means process now.
	chef.Enqueue("nqueue", []byte(fmt.Sprintf("nqueue job %d", 1)), 0)
  
  // enqueue job with priority
  chef.EnqueueWithPriority(queue string, data []byte, priority PriorityLevel, delay uint32)
```

#### Feature

- Processing middleware
- Delay Job
- Auto retry if processing job fail with panic
- Priority task by level
- Concurrency workers
- Graceful shutdown
