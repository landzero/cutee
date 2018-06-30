package routes

import (
	"net/http"

	"github.com/landzero/cutee/types"
	"landzero.net/x/database/orm"
	"landzero.net/x/net/web"
	"landzero.net/x/net/web/cache"
)

const healthCheckKey = "_health_check"

func routeNotImplemented(ctx *web.Context) {
	ctx.PlainText(http.StatusOK, []byte("// TODO: not implemented"))
}

func routeHealthCheck(ctx *web.Context, db *orm.DB, che cache.Cache) {
	// check redis via cache service
	defer che.Delete(healthCheckKey)
	if err := che.Put(healthCheckKey, "anything", 10); err != nil {
		ctx.PlainText(http.StatusInternalServerError, []byte("BAD REDIS:"+err.Error()))
		return
	}
	// check database via SELECT COUNT(*) FROM users
	var count int
	if err := db.Model(&types.User{}).Count(&count).Error; err != nil {
		ctx.PlainText(http.StatusInternalServerError, []byte("BAD DATABASE:"+err.Error()))
		return
	}
	ctx.PlainText(http.StatusOK, []byte("OK"))
}

// Mount mount all routes
func Mount(w *web.Web) {
	w.Get("/_health_check", routeHealthCheck)
	w.Use(InjectErrorHelper())
	w.Use(InjectAuthHelper())
	w.Get("/", routeIndex)
	w.Get("/about", routeAbout)
	w.Get("/login/twitter", routeLoginTwitter)
	w.Get("/login/twitter/callback", routeLoginTwitterCallback)
	w.Post("/logout", routeLogout)
	w.Get("/profile", RequireLoggedIn(), routeProfile)
}
