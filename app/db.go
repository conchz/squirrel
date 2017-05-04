package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
)

const (
	timeZone         = "Asia/Shanghai"
	mysqlStoreEngine = "InnoDB"
	mysqlCharset     = "utf8mb4"
)

var engine *xorm.Engine

func init() {
	err := connectToMySQL(LoadConfig())
	if err != nil {
		panic(err)
	}
}

func connectToMySQL(config *Config) error {
	mysqlConfig := config.MySQLConf

	var err error
	engine, err = xorm.NewEngine(mysqlConfig.Driver, mysqlConfig.DataSource)
	if err != nil {
		return err
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}

	engine.TZLocation = location
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
