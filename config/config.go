package config

import (
	"log"
	"strings"
	"time"

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

type PasetoTokenConfig struct {
	Paseto struct {
		SymmetricKey             string
		Access_token_expiration  time.Duration
		Refresh_token_expiration time.Duration
	}
}

type ServerConfig struct {
	Server struct {
		Port string
	}
}

func LoadPasetoTokenConfig(path string) (*PasetoTokenConfig, error) {
	viper.Reset()
	config := new(PasetoTokenConfig)

	// tell viper from where to read
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
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
		log.Printf("[PASETO CONFIG] | error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("[PASETO CONFIG] | error unmarshling the data from config file : %v \n", err)
	}

	// Add this after unmarshaling the configuration
	return config, nil
}

func LoadPostgresConfig(path string) (*PGConfig, error) {
	viper.Reset()
	config := new(PGConfig)

	// tell viper from where to read
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
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
		log.Printf("[POSTGRES CONFIG] | error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("[POSTGRES CONFIG] | error unmarshling the data from config file : %v \n", err)
	}

	log.Println("Environment Variable:", viper.Get("PGCONFIG_POSTGRESDB_FADY"))

	// If you've set the environment variable correctly, it will override the value from config.yaml
	log.Printf("Config Value: %s", config.Postgresdb.Host)

	return config, nil
}

func LoadServerConfigs(path string) (*ServerConfig, error) {
	viper.Reset()
	config := new(ServerConfig)

	// tell viper from where to read
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
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
		log.Printf("[SERVER CONFIG] | error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("[SERVER CONFIG] | error unmarshling the data from config file : %v \n", err)
	}

	// Add this after unmarshaling the configuration
	return config, nil
}
