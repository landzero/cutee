package main

import (
	"os"

	"landzero.net/x/os/osext"

	"github.com/landzero/cutee/routes"
	"github.com/landzero/cutee/types"
	"landzero.net/x/database/orm"
	_ "landzero.net/x/database/orm/dialects/postgres"
	"landzero.net/x/log"
	"landzero.net/x/net/oauth"
	"landzero.net/x/net/web"
	"landzero.net/x/net/web/cache"
	_ "landzero.net/x/net/web/cache/redis"
	"landzero.net/x/net/web/i18n"
	"landzero.net/x/net/web/session"
	_ "landzero.net/x/net/web/session/redis"
)

func main() {
	defer osext.DoExit()

	// create Web
	m := web.Modern()
	// use i18n
	m.Use(i18n.I18ner(i18n.Options{
		Directory:   "locales",
		BinFS:       m.IsProduction(),
		Locales:     []string{"en-US", "zh-CN", "ja"},
		LocaleNames: []string{"english", "简体中文", "日本語"},
	}))
	// create options
	opt := types.Options{
		Domain: os.Getenv("DOMAIN"),
	}
	m.Map(opt)
	// create DB
	var db *orm.DB
	var err error
	if db, err = orm.Open("postgres", os.Getenv("DATABASE_URL")); err != nil {
		log.Println("failed to initialize db:", err)
		osext.WillExit(1)
		return
	}
	m.Map(db.LogMode(!m.IsProduction()))
	// create oauth.Consumer
	csm := oauth.NewConsumer(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})
	m.Map(csm)
	// use cache
	m.Use(cache.Cacher(cache.Options{
		Adapter:       "redis",
		AdapterConfig: os.Getenv("REDIS_URL"),
	}))
	// use session
	m.Use(session.Sessioner(session.Options{
		Adapter:       "redis",
		AdapterConfig: os.Getenv("REDIS_URL"),
		Secure:        m.IsProduction(),
	}))
	// mount routes
	routes.Mount(m)
	// run
	m.Run() // run with $HOST, $PORT, $WEB_ENV
}
