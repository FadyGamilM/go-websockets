package postgres

import (
	"database/sql"
	"log"
	"net/url"

	"github.com/FadyGamilM/go-websockets/config"
)

func SetupPostgresConnection() (*sql.DB, error) {
	configs, err := config.LoadPostgresConfig("./config")
	if err != nil {
		log.Println("error trying to load config variables", err)
		return nil, err
	}
	// construct the conn string
	dsn := url.URL{
		Scheme: "postgres",
		Host:   configs.Postgresdb.Host,
		User:   url.UserPassword(configs.Postgresdb.User, configs.Postgresdb.Password),
		Path:   configs.Postgresdb.Dbname,
	}

	q := dsn.Query()
	q.Add("sslmode", configs.Postgresdb.Sslmode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		log.Println("error trying to open a postgres connection", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("app has ping, but db server didn't pong :( \n : %v \n", err)
		return nil, err
	}

	log.Println("successfully connected to the postgres instance on port : ", configs.Postgresdb.Host)
	return db, nil
}
