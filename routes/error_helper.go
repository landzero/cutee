package routes

import (
	"net/http"

	"landzero.net/x/net/web"
)

// ErrorHelper error helper
type ErrorHelper struct {
	ctx *web.Context
}

// InternalServerError render a 500 page
func (eh *ErrorHelper) InternalServerError(str string) {
	eh.ctx.Data["Error"] = str
	eh.ctx.Data["Crid"] = eh.ctx.Crid()
	eh.ctx.HTML(http.StatusInternalServerError, "error")
}

// BadRequest render a 400 page
func (eh *ErrorHelper) BadRequest(str string) {
	eh.ctx.Data["Error"] = str
	eh.ctx.Data["Crid"] = eh.ctx.Crid()
	eh.ctx.HTML(http.StatusBadRequest, "error")
}

// Redirect render a redirect page
func (eh *ErrorHelper) Redirect(str string) {
	eh.ctx.Data["RedirectURL"] = str
	eh.ctx.HTML(http.StatusOK, "redirect")
}

// InjectErrorHelper inject ErrorHelper to web.Context
func InjectErrorHelper() web.Handler {
	return func(ctx *web.Context) {
		eh := ErrorHelper{ctx: ctx}
		ctx.Map(&eh)
		ctx.Next()
	}
}
