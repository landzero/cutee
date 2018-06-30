package routes

import (
	"github.com/landzero/cutee/types"
	"landzero.net/x/database/orm"
	"landzero.net/x/net/web"
	"landzero.net/x/net/web/session"
)

const accountKey = "Account"
const redirectURLKey = "RedirectURL"

// AuthHelper auth helper
type AuthHelper struct {
	db   *orm.DB
	sess session.Store
	user *types.User
}

func (ah *AuthHelper) setup() {
	account, ok := ah.sess.Get(accountKey).(string)
	if ok {
		u := types.User{}
		if ah.db.First(&u, map[string]interface{}{"account": account}).Error == nil && u.ID > 0 {
			ah.user = &u
		} else {
			ah.sess.Delete(accountKey)
		}
	} else {
		ah.sess.Delete(accountKey)
	}
}

// IsLoggedIn is session logged in
func (ah *AuthHelper) IsLoggedIn() bool {
	return ah.user != nil
}

// User returns the user
func (ah *AuthHelper) User() *types.User {
	return ah.user
}

// SetUser set session user, nil will clear session
func (ah *AuthHelper) SetUser(u *types.User) {
	if u == nil {
		ah.sess.Delete("Account")
	} else {
		ah.sess.Set("Account", u.Account)
	}
}

// InjectAuthHelper inject auth helper
func InjectAuthHelper() web.Handler {
	return func(ctx *web.Context, db *orm.DB, sess session.Store) {
		ah := &AuthHelper{sess: sess, db: db}
		ah.setup()
		ctx.Data["IsLoggedIn"] = ah.IsLoggedIn()
		ctx.Data["User"] = ah.User()
		ctx.Map(ah)
		ctx.Next()
	}
}

// RequireLoggedIn require logged-in
func RequireLoggedIn() web.Handler {
	return func(ctx *web.Context, ah *AuthHelper, sess session.Store) {
		if !ah.IsLoggedIn() {
			sess.Set(redirectURLKey, ctx.Req.URL.String())
			ctx.Redirect("/login/twitter")
		} else {
			ctx.Next()
		}
	}
}
