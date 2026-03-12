package config

import (
	flag2 "flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Address string
}
type Config struct {
	//These are called struct tags in go `yaml:"env" env-required:"true"`
	Env              string `yaml:"env" env-required:"true"`
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	HTTPServerConfig `yaml:"http_server"`
}

// MustLoad config parsing is done over here and we can use this config to use anywhere in our program
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag2.String("config", "", "Path to config file")
		flag2.Parse()

		configPath = *flags
		if configPath == "" {
			log.Fatal("config file path required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file %s", err.Error())
	}

	return &cfg
}
