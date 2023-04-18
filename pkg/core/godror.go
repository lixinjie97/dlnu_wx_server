package core

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/gowechat/example/config"
)

var Db *sql.DB

func Setup() {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(`user=%s password=%s connectString="(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=%d)(SID=%s))(CONNECT_DATA=(SID=%s)))"`, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Sid, cfg.Sid)
	DB, err := sql.Open("godror", dsn)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	DB.SetMaxIdleConns(1000)
	DB.SetMaxOpenConns(2000)
	Db = DB
}