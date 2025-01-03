package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	// the variavle env mentioned in local.yaml , put the info in env
	// env-default:"production"
	Env         string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// this file must run  , error occur hame app nahi start karna hai the the application shd not run
// config is imp
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	// if config is empty string then inside if
	// wheather is it in CLI
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file") /// name, value, usage
		flag.Parse()

		configPath = *flags

		// if still it is emty
		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	// we only want to check err , we can write in two different file
	// we can use ; n use in same line
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist : %s", configPath)
	}

	// final we will serialize in Config struct
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg) // & address
	if err != nil {
		log.Fatalf("cannot read config file : %s", err.Error())
	}

	return &cfg
}
