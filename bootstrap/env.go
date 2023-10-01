package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv string `mapstructure:"APP_ENV"`
	DbUser string `mapstructure:"MONGO_USER"`
	DbPass string `mapstructure:"MONGO_PASS"`
	DbHost string `mapstructure:"MONGO_HOST"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Error while reading config file", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env

}
