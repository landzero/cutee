package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/landzero/cutee/types"
	"landzero.net/x/database/orm"
	"landzero.net/x/net/oauth"
	"landzero.net/x/net/web"
	"landzero.net/x/net/web/cache"
	"landzero.net/x/net/web/session"
)

const twitterRequestTokenPrefix = "twitter.request-token."

func routeLoginTwitter(ctx *web.Context, csm *oauth.Consumer, opt types.Options, eh *ErrorHelper, cch cache.Cache) {
	callbackURL := fmt.Sprintf("%s/login/twitter/callback", opt.Domain)
	token, loginURL, err := csm.GetRequestTokenAndUrl(callbackURL)
	if err != nil {
		eh.InternalServerError(err.Error())
		return
	}
	if err = cch.Put(twitterRequestTokenPrefix+token.Token, token.Secret, 300); err != nil {
		eh.InternalServerError("failed to write Redis")
		return
	}
	eh.Redirect(loginURL)
}

func routeLoginTwitterCallback(
	ctx *web.Context,
	csm *oauth.Consumer,
	eh *ErrorHelper,
	ah *AuthHelper,
	cch cache.Cache,
	sess session.Store,
	db *orm.DB,
) {
	// check denied
	denied := ctx.Query("denied")
	if len(denied) > 0 {
		eh.BadRequest("you've cancelled login")
		cch.Delete(twitterRequestTokenPrefix + denied)
		return
	}
	// verify request token
	verifier := ctx.Query("oauth_verifier")
	token := ctx.Query("oauth_token")
	if len(token) == 0 {
		eh.BadRequest("invalid parameter")
		return
	}
	secret, ok := cch.Get(twitterRequestTokenPrefix + token).(string)
	if !ok {
		eh.InternalServerError("token missing")
		return
	}
	accessToken, err := csm.AuthorizeToken(&oauth.RequestToken{Token: token, Secret: secret}, verifier)
	if err != nil {
		eh.BadRequest("invalid verification code")
		return
	}
	// call Twitter API for credentials verification
	client, err := csm.MakeHttpClient(accessToken)
	if err != nil {
		eh.InternalServerError("failed to create Twitter API client")
		return
	}
	resp, err := client.Get("https://api.twitter.com/1.1/account/verify_credentials.json")
	if err != nil || resp.StatusCode != http.StatusOK {
		eh.InternalServerError("failed to dial Twitter API")
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		eh.InternalServerError("failed to read Twitter API")
		return
	}
	var a types.TAccount
	if err = json.Unmarshal(buf, &a); err != nil {
		eh.InternalServerError("failed to decode Twitter API")
		return
	}
	if len(a.ID) == 0 || len(a.ScreenName) == 0 {
		eh.InternalServerError("failed to decode Twitter API")
		return
	}
	// update database
	u := types.User{}
	if db.Where(map[string]interface{}{
		"account": "twitter-" + a.ID,
	}).Assign(map[string]interface{}{
		"name":   a.ScreenName,
		"avatar": a.AvatarURL,
	}).FirstOrCreate(&u).Error != nil {
		eh.InternalServerError("failed to create user")
		return
	}
	// update session
	ah.SetUser(&u)
	// determine redirection target
	if redirectURL, ok := sess.Get(redirectURLKey).(string); ok && len(redirectURL) > 0 {
		sess.Delete(redirectURL)
		ctx.Redirect(redirectURL)
	} else {
		ctx.Redirect("/")
	}
}

func routeLogout(ctx *web.Context, ah *AuthHelper) {
	ah.SetUser(nil)
	ctx.Redirect("/")
}
