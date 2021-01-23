package config

import (
	"github.com/fsnotify/fsnotify"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const defaultConfigFile = "config/config.yaml"

var InitViper *viper.Viper
var InitConfig Server

func init() {
	vip := viper.New()
	vip.SetConfigFile(defaultConfigFile)
	err := vip.ReadInConfig()
	if err != nil {
		logger.Errorf("Failed to get config file!")
		panic(err.Error())
	}
	vip.WatchConfig()
	vip.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: %s", e.Name)
		if err := vip.Unmarshal(&InitConfig); err != nil {
			logger.Errorf("Failed to resolve config file!")
			panic(err.Error())
		}
	})
	if err := vip.Unmarshal(&InitConfig); err != nil {
		logger.Errorf("Failed to resolve config file!")
		panic(err.Error())
	}
	InitViper = vip
}

type Server struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Db     Db     `mapstructure:"db" json:"db" yaml:"db"`
	Log    Log    `mapstructure:"log" json:"log" yaml:"log"`
}

type System struct {
	Host             string `mapstructure:"host" json:"host" yaml:"host"`
	AllowIps         string `mapstructure:"allowIps" json:"allowIps" yaml:"allowIps"`
	AppName          string `mapstructure:"appName" json:"appName" yaml:"appName"`
	ConcurrencyQueue int    `mapstructure:"concurrencyQueue" json:"concurrencyQueue" yaml:"concurrencyQueue"`
	ApiSignEnable    bool   `mapstructure:"apiSignEnable" json:"apiSignEnable" yaml:"apiSignEnable"`
	ApiKey           string `mapstructure:"apiKey" json:"apiKey" yaml:"apiKey"`
	ApiSecret        string `mapstructure:"apiSecret" json:"apiSecret" yaml:"apiSecret"`
	AuthSecret       string `mapstructure:"authSecret" json:"authSecret" yaml:"authSecret"`
	EnableTls        bool   `mapstructure:"enableTls" json:"enableTls" yaml:"enableTls"`
	CaFile           string `mapstructure:"caFile" json:"caFile" yaml:"caFile"`
	CertFile         string `mapstructure:"certFile" json:"certFile" yaml:"certFile"`
	KeyFile          string `mapstructure:"keyFile" json:"keyFile" yaml:"keyFile"`
}

type Db struct {
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Database     string `mapstructure:"database" json:"database" yaml:"database"`
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Charset      string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns int    `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
}

type Log struct {
	Path  string `mapstructure:"path" json:"path" yaml:"path"`
	Name  string `mapstructure:"name" json:"name" yaml:"name"`
	Level string `mapstructure:"level" json:"level" yaml:"level"`
}
