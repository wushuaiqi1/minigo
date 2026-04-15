package config

import (
	"flag"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	Production = "production" // 生产环境
	Release    = "release"    // 测试环境
)

var (
	globalConfig   *Config // 全局共有配置
	once           sync.Once
	configFileName string // 配置文件名称
)

func init() {
	flag.StringVar(&configFileName, "c", Release, "配置文件名称")
	// 饿汉式加载
	once.Do(func() {
		setViper()
	})
}

type Config struct {
	Server ServerConfig `json:"server"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Env  string `json:"env"`
}

func GetConfigInstance() *Config {
	if globalConfig == nil {
		// 懒汉式加载
		once.Do(func() {
			setViper()
		})
	}
	return globalConfig
}

func setViper() {

	viper.SetConfigName(configFileName)
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")       // 工作目录找
	viper.AddConfigPath("./config") // 工作目录/Config目录找

	viper.SetDefault("server.publicAddr", ":19001")
	viper.SetDefault("server.internalAddr", "127.0.0.1:19002")
	viper.SetDefault("server.readTimeout", "10s")
	viper.SetDefault("server.writeTimeout", "10s")
	viper.SetDefault("database.dsn", "root:password@tcp(127.0.0.1:3306)/horserun?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("database.maxIdleConns", 10)
	viper.SetDefault("database.maxOpenConns", 100)
	viper.SetDefault("database.connMaxLifetime", "1h")
	viper.SetDefault("authCode.codeLength", 16)

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("警告：无法读取配置文件，使用默认值: %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	globalConfig = &config
	logrus.WithFields(logrus.Fields{
		"configFile": configFileName,
		"env":        config.Server.Env,
		"port":       config.Server.Port,
	}).Info("初始化Viper配置Success")
}
