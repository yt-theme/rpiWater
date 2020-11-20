package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

var Cfg Config

type Config struct {
    SOCKETAddr string `yaml:"SOCKETAddr"`
}

func (self *Config) Read () (*Config, error) {
    file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	
	// 解析配置文件
	err = yaml.Unmarshal(file, self)
	if err != nil {
		return nil, err
	}
	
	return self, nil
}

func init () {
    _, err := Cfg.Read()
    if err != nil {
        log.Fatalln(err)
    }
}