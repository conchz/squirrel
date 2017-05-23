package app

import (
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

// http://stackoverflow.com/questions/20240179/nil-detection-in-go
var (
	config *Config
	lock   sync.RWMutex
)

func LoadConfig() *Config {
	if config != nil {
		return config
	}

	lock.RLock()
	defer lock.RUnlock()

	bytes := GetAppConfBytes()
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	// Setup server port
	port = config.ServerConf.Port

	return config
}
