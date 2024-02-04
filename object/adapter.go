package object

import (
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // db = mysql
	_ "github.com/lib/pq"              // db = postgres
	"github.com/spf13/viper"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var engine *xorm.Engine

func InitAdapter() {
	dbUrl := viper.GetString(ENV_DB_URL)
	spits := strings.Split(dbUrl, "://")
	var connStr string
	if spits[0] == "postgres" {
		connStr = dbUrl
	} else if spits[0] == "mysql" {
		connStr = spits[1]
	} else {
		fmt.Println("No YUDAI_DB_URL environment variable provided.")
		os.Exit(10)
		return
	}
	_engine, err := xorm.NewEngine(spits[0], connStr)
	if err != nil {
		fmt.Println("Error when initialize DB.", dbUrl, err)
		os.Exit(11)
		return
	}
	engine = _engine
	dbPrefix := viper.GetString(ENV_DB_PREFIX)
	if dbPrefix != "" {
		// remove last _ before append _
		dbPrefix = strings.TrimSuffix(dbPrefix, "_") + "_"
	}
	tbMapper := names.NewPrefixMapper(names.GonicMapper{}, dbPrefix)
	engine.SetTableMapper(tbMapper)
	// engine.SetColumnMapper(names.GonicMapper{})
	engine.DatabaseTZ = time.UTC

	engine.Sync(new(Token), new(Provider), new(User), new(UserSocial))
}
