package salt

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func RBAC(action string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			if ctx.Can(action) {
				next(ctx)
				return
			}
			ctx.Error(WithCode(403), WithMessage("Forbidden"))
		}
	}
}
