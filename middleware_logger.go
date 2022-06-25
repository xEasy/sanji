package sanji

import (
	"strings"
	"time"
)

func logger_middleware(ctx *Context) {
	prefix := strings.Join([]string{"Namespace:", ctx.Job.Namespace, "Queue:", ctx.Job.Queue, "JID:", ctx.Job.ID}, " ")
	start := time.Now()
	Logger.Println(prefix, "start")
	ctx.Next()
	Logger.Println(prefix, "done: ", time.Since(start))
}
