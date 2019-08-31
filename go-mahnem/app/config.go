package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var _config *AppConfig

// DatabaseConfig application config
type DatabaseConfig struct {
	URL      string
	Port     int16
	Database string
	Username string
	Password string
}

// SiteConfig application config
type SiteConfig struct {
	URL      string //mahnem location
	Login    string
	Password string
}

// AppConfig is main application configuration
type AppConfig struct {
	Db   DatabaseConfig
	Site SiteConfig
}

// DbConfig - database config
func (cfg AppConfig) DbConfig() *DatabaseConfig {
	return &cfg.Db
}

// SiteConfig - web site configuration
func (cfg AppConfig) SiteConfig() *SiteConfig {
	return &cfg.Site
}

func init() {

	log.Println("init [config.go]")

	var err error
	_config, err = readAppConfiguration()
	if _config == nil || err != nil {
		log.Fatalf("Can't read config file, %s\n", err.Error())
	}
}

func readAppConfiguration() (*AppConfig, error) {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return populateConfig()
}

func populateConfig() (*AppConfig, error) {

	var cfg AppConfig

	if err := fetchConfig("db", &cfg.Db); err != nil {
		return nil, err
	}

	if err := fetchConfig("site", &cfg.Site); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchConfig(section string, ptr interface{}) error {

	vals := viper.GetStringMap(section)
	if vals == nil {
		return fmt.Errorf("No %s section in config", section)
	}

	err := mapstructure.Decode(vals, ptr)
	if err != nil {
		return fmt.Errorf("Error during encoding config->%s, %s", section, err.Error())
	}

	return nil
}

// GetAppConfig returns application configuration
func GetAppConfig() *AppConfig {
	return _config
}
