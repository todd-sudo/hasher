package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBName             string `env:"DB_NAME" env-default:"hasher_database.db"`
	Limit              int    `env:"LIMIT" env-default:"5"`
	PrivatePemFileName string `env:"PRIVATE_PEM_FILE_NAME" env-default:"private.pem"`
	SizeRSAKey         int    `env:"SIZE_RSA_KEY" env-default:"10000"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Go-rshok hasher"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
