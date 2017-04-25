package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func Connect() *xorm.Engine {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:test1234@localhost:3306/test?charset=utf8")
	if err != nil {
		panic(err)
	}

	return engine
}
