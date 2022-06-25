package sanji

func recovery_middleware(ctx *Context) {
	defer Recovery(func(err any) {
		ctx.acknowledge = false
		panic(err)
	})

	ctx.Next()
}
