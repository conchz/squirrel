package app

import "gopkg.in/yaml.v2"

type Config struct {
	ServerConf struct {
		Port uint16 `yaml:"port"`
	} `yaml:"server"`
	LoggingConf struct {
		Level    string `yaml:"level"`
		FilePath string `yaml:"filePath"`
	} `yaml:"logging"`
	MySQLConf struct {
		Driver     string `yaml:"driver"`
		DataSource string `yaml:"dataSource"`
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

	bytes := ConfigFile()
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}
