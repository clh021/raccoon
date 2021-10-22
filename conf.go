package raccoon

import (
	"fmt"

	"github.com/spf13/viper"
)

type Step struct {
	Tag     string   `yaml:"tag"`
	Webroot string   `yaml:"webroot"`
	Webaddr string   `yaml:"webaddr"`
	Exec    string   `yaml:"exec"`
	Param   string   `yaml:"param"`
	Envs    []string `yaml:"envs"`
}
type Config struct {
	Step []Step `yaml:"step"`
}

func getConf() *[]Step {
	viper.SetConfigName("raccoon")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GetProgramPath())
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	return &c.Step
}

// type (s *Config) WebServerCnt() int{

// 	return le
// }
