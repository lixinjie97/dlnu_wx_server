package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	cfgFile = flag.String("config", "./config.yaml", "配置文件路径")

	cfg *Config
)

//Config example config
type Config struct {
	Listen string `yaml:"listen"`
	Redis  struct {
		Host        string `yaml:"host"`
		Password    string `yaml:"password"`
		Database    int    `yaml:"database"`
		MaxActive   int    `yaml:"maxActive"`
		MaxIdle     int    `yaml:"maxIdle"`
		IdleTimeout int    `yaml:"idleTimeout"`
	} `yaml:"redis"`
	*OfficialAccountConfig `yaml:"officialAccountConfig"`
	*Oracle                `yaml:"oracle"`
	*Runtime               `yaml:"runtime"`
}

//OfficialAccountConfig 公众号相关配置
type OfficialAccountConfig struct {
	AppID          string `yaml:"appID"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
}

type Oracle struct {
	User                 string `yaml:"user"`
	Password             string `yaml:"password"`
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
	Sid                  string `yaml:"sid"`
}

type Runtime struct {
	Log                  Log `yaml:"log"`
}

type Log struct {
	LogRootPath          string `yaml:"logRootPath"`
	DisableTimestamp     bool   `yaml:"disableTimestamp"`
	TimestampFormat      string `yaml:"timestampFormat"`
}

//GetConfig 获取配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	bytes, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	cfgData := &Config{}
	err = yaml.Unmarshal(bytes, cfgData)
	if err != nil {
		panic(err)
	}
	cfg = cfgData
	return cfg
}
