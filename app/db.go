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

func ConnectToMySQL(config *Config) error {
	mysqlConfig := config.MySQLConf

	var err error
	engine, err = xorm.NewEngine(mysqlConfig.Driver, mysqlConfig.DataSource)
	if err != nil {
		return err
	}

	engine.ShowSQL(true)
	engine.StoreEngine(mysqlStoreEngine)
	engine.Charset(mysqlCharset)
	engine.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	engine.SetMaxOpenConns(mysqlConfig.MaxOpenConns)

	return nil
}

func GetXormEngine() *xorm.Engine {
	return engine
}
