package routes

import (
	"net/http"

	"landzero.net/x/net/web"
)

func routeIndex(ctx *web.Context) {
	ctx.Data["Nav_Index"] = "active"
	ctx.HTML(http.StatusOK, "index")
}

func routeAbout(ctx *web.Context) {
	ctx.Data["Nav_About"] = "active"
	ctx.HTML(http.StatusOK, "about")
}
