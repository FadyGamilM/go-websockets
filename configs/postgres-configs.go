package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type PGConfig struct {
	Postgresdb struct {
		Host     string
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}
}

func LoadPostgresConfig(path string) (*PGConfig, error) {
	viper.Reset()
	config := new(PGConfig)

	// tell viper from where to read
	viper.AddConfigPath(path)
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")

	// configure the feature to override the vars from the yaml file via the environment variables
	viper.AutomaticEnv()
	/*
		viper reads the vars from the yaml file as following :
		SERVER.PORT , but we can't define an env variable with the dot notation, so we will define it with _ and replace the default behaviour of viper
	*/
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// read the configs
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("error unmarshling the data from config file : %v \n", err)
	}

	return config, nil
}
