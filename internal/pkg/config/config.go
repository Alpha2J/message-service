package config

import (
	"fmt"
	"github.com/spf13/viper"
	"message-service/internal/pkg/helper"
	"os"
)

type ServerConfig struct {
	EnableDebug bool           `mapstructure:"enable_debug"`
	Env         string         `mapstructure:"env"`
	App         AppConfig      `mapstructure:"app"`
	Mysql       MySqlConfig    `mapstructure:"mysql"`
	Business    BusinessConfig `mapstructure:"business"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type MySqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

// custom this with different business
type BusinessConfig struct {
	OwnerToken string `mapstructure:"owner_token"`
}

const (
	ENV_LOCAL = "local"
	ENV_TEST  = "test"
	ENV_PROD  = "prod"
)

var validEnvs = [...]string{ENV_LOCAL, ENV_TEST, ENV_PROD}

var config ServerConfig
var env string

func init() {
	configsBaseDir, _ := helper.FindProjectRoot()
	configsDir := configsBaseDir + "/configs"

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configsDir)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config files, %s", err))
	}

	env = os.Getenv("go_service_env")
	if env == "" {
		env = viper.GetString("env")
	}

	if env == "" {
		env = "local"
	} else {
		isValidEnv := false
		for _, v := range validEnvs {
			if env == v {
				isValidEnv = true
				break
			}
		}
		if !isValidEnv {
			panic(fmt.Sprintf("Invalid env [%s], only supports one of 'local','prod','test'", env))
		}
	}

	envConfigFile := fmt.Sprintf("application.%s.yaml", env)
	viper.SetConfigName(envConfigFile)
	viper.AddConfigPath(configsDir)
	if err := viper.MergeInConfig(); err != nil {
		panic(fmt.Sprintf("Error merging config file of env [%s]", env))
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode config of env [%s] into struct", env))
	}
}

func GetServerConfig() *ServerConfig {
	return &config
}

func GetAppConfig() *AppConfig {
	return &config.App
}

func GetMySqlConfig() *MySqlConfig {
	return &config.Mysql
}

func GetBusinessConfig() *BusinessConfig {
	return &config.Business
}

func GetEnv() string {
	return env
}

func IsDebugging() bool {
	return config.EnableDebug
}

func IsReleaseMode() bool {
	return env == ENV_PROD
}
