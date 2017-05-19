package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
)

type MySQLTemplate struct {
	engine *xorm.Engine
}

const (
	mysqlStoreEngine = "InnoDB"
	mysqlCharset     = "utf8mb4"
)

var (
	mySQLTemplate *MySQLTemplate
	mutex         = new(sync.RWMutex)
)

// Call this function when initializing server
func ConnectToMySQL(config *Config) {
	mysqlConfig := config.MySQLConf

	mutex.Lock()
	defer mutex.Unlock()

	engine, err := xorm.NewEngine(mysqlConfig.Driver, mysqlConfig.DataSource)
	if err != nil {
		panic(err)
	}

	engine.StoreEngine(mysqlStoreEngine)
	engine.Charset(mysqlCharset)
	engine.ShowSQL(mysqlConfig.ShowSQL)
	engine.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	engine.SetMaxOpenConns(mysqlConfig.MaxOpenConns)

	mySQLTemplate = &MySQLTemplate{
		engine: engine,
	}
}

func GetMySQLTemplate() *MySQLTemplate {
	return mySQLTemplate
}

func (t *MySQLTemplate) XormEngine() *xorm.Engine {
	return t.engine
}

func (t *MySQLTemplate) Insert(bean interface{}) (int64, error) {
	return t.engine.Insert(bean)
}

// Retrieve one record from database by id
// bean: a point of struct
func (t *MySQLTemplate) GetById(id int64, bean interface{}) interface{} {
	exists, err := t.engine.Id(id).Get(bean)
	if err != nil || !exists {
		return nil
	}
	return bean
}

// Retrieve one record from database, bean's non-empty fields will be as conditions
// bean: a point of struct
func (t *MySQLTemplate) GetByNonEmptyFields(bean interface{}) interface{} {
	exists, err := t.engine.Get(bean)
	if err != nil || !exists {
		return nil
	}
	return bean
}

func (t *MySQLTemplate) Close() error {
	return t.engine.Close()
}
