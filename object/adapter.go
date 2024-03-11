package object

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // db = mysql
	_ "github.com/lib/pq"              // db = postgres
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var engine *xorm.Engine

func InitAdapter() {
	_engine, err := xorm.NewEngine(Conf.DbDriver, Conf.DbConnStr)
	if err != nil {
		fmt.Printf("Initialize database failed: %v\n", err)
		os.Exit(1)
		return
	}

	engine = _engine
	tbMapper := names.NewPrefixMapper(names.GonicMapper{}, Conf.DbTablePrefix)
	engine.SetTableMapper(tbMapper)
	// engine.SetColumnMapper(names.GonicMapper{})
	engine.DatabaseTZ = time.UTC

	_ = engine.Sync(new(Token), new(Provider), new(User), new(UserSocial), new(VerificationRecord))
}
