package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"gopkg.in/redis.v5"
	"gopkg.in/yaml.v2"
	"sync"
)

type Config struct {
	ServerConf struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	MySQLConf struct {
		Driver       string `yaml:"driver"`
		DataSource   string `yaml:"dataSource"`
		ShowSQL      bool   `yaml:"showSql"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"mysql"`
	RedisConf struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

type MySQLTemplate struct {
	engine *xorm.Engine
}

const (
	mysqlStoreEngine = "InnoDB"
	mysqlCharset     = "utf8mb4"
)

// http://stackoverflow.com/questions/20240179/nil-detection-in-go
var (
	mutex         sync.RWMutex
	config        *Config
	mySQLTemplate *MySQLTemplate
	redisClient   *redis.Client
)

func LoadConfig() *Config {
	if config != nil {
		return config
	}

	mutex.Lock()
	defer mutex.Unlock()

	bytes := GetAppConfBytes()
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	// Setup server port
	port = config.ServerConf.Port

	return config
}

func ConnectToRedis(config *Config) {
	redisConfig := config.RedisConf
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func CloseRedisClient() error {
	return redisClient.Close()
}

func ConnectToMySQL(config *Config) {
	mysqlConfig := config.MySQLConf

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
