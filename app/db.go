package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	mysqlStoreEngine = "InnoDB"
	mysqlCharset     = "utf8mb4"
)

var engine *xorm.Engine

func ConnectToMySQL(config *Config) {
	mysqlConfig := config.MySQLConf

	var err error
	engine, err = xorm.NewEngine(mysqlConfig.Driver, mysqlConfig.DataSource)
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(false)
	engine.StoreEngine(mysqlStoreEngine)
	engine.Charset(mysqlCharset)
	engine.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	engine.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
}

func GetXormEngine() *xorm.Engine {
	return engine
}

func Insert(bean interface{}) (int64, error) {
	return engine.Insert(bean)
}

func FindById(id int64, bean interface{}) interface{} {
	engine.Id(id).Get(bean)
	return bean
}
