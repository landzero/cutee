package routes

import (
	"net/http"

	"landzero.net/x/net/web"
)

func routeProfile(ctx *web.Context) {
	ctx.Data["Nav_Profile"] = "active"
	ctx.HTML(http.StatusOK, "profile")
}
