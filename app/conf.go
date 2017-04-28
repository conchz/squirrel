package app

import "gopkg.in/yaml.v2"

type Config struct {
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

func GetConfig() *Config {
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
