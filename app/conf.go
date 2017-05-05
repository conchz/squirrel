package app

import (
	"gopkg.in/yaml.v2"
	"sync"
)

var lock sync.RWMutex

type Config struct {
	ServerConf struct {
		Port uint16 `yaml:"port"`
	} `yaml:"server"`
	LoggingConf struct {
		Level    string `yaml:"level"`
		FilePath string `yaml:"filePath"`
	} `yaml:"logging"`
	MySQLConf struct {
		Driver       string `yaml:"driver"`
		DataSource   string `yaml:"dataSource"`
		MaxIdleConns int `yaml:"maxIdleConns"`
		MaxOpenConns int `yaml:"maxOpenConns"`
	} `yaml:"mysql"`
	RedisConf struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
}

// http://stackoverflow.com/questions/20240179/nil-detection-in-go
var config *Config

func LoadConfig() *Config {
	if config != nil {
		return config
	}

	lock.RLock()
	defer lock.RUnlock()

	bytes := ConfigFile()
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}
