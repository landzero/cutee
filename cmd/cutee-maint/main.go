package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/landzero/cutee/types"
	"landzero.net/x/database/orm"
	_ "landzero.net/x/database/orm/dialects/postgres"
	"landzero.net/x/log"
	"landzero.net/x/net/http/httpext"
	"landzero.net/x/os/osext"
)

var isMigration bool
var isHealthCheck bool

func main() {
	defer osext.DoExit()

	// decode flags
	flag.BoolVar(&isMigration, "migrate", false, "run in migration mode")
	flag.BoolVar(&isHealthCheck, "health-check", false, "run health check")
	flag.Parse()

	// run health check
	if isHealthCheck {
		if httpext.HealthCheck(fmt.Sprintf("http://127.0.0.1:%s/_health_check", os.Getenv("PORT"))) {
			osext.WillExit(0)
		} else {
			osext.WillExit(1)
		}
		return
	}

	// run database migration
	if isMigration {
		var db *orm.DB
		var err error
		if db, err = orm.Open("postgres", os.Getenv("DATABASE_URL")); err != nil {
			log.Println("failed to initialize db,", err)
			osext.WillExit(1)
			return
		}
		db.LogMode(true).AutoMigrate(&types.User{})
	}
}
