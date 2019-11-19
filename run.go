package main

import (
	"github.com/dovadi/dbconfig"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"spider/crawler"
	"os"
	"time"
)


func main()  {
	// logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
			TimeFormat: time.RFC3339,
		},
	)

	// DB
	connectionString := dbconfig.PostgresConnectionString("config/settings.json",  "disable")

	log.Info().Msg(connectionString)
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}else{
		log.Info().Msg("Connected to DB!")
	}

	service := crawler.NewCrawlerHandler(db)
	service.Start()
	defer db.Close()

}
