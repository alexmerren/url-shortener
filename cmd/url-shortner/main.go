package main

import (
	"context"
	"log"
	"url-shortner/internal/config"
	"url-shortner/internal/datastore"
	"url-shortner/internal/logger"
	"url-shortner/internal/rest"
)

func main() {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	loglevel := config.Logger.Level
	logger, err := logger.NewZapLogger(loglevel)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	ctx := context.Background()
	dbName := config.Database.Filename
	dbCapacity := config.Database.Capacity
	sqliteInput := &datastore.SqliteDatabaseInput{
		Logger:           logger,
		DatabaseName:     dbName,
		DatabaseCapacity: dbCapacity,
		Ctx:              ctx,
	}
	database, err := datastore.NewSqliteDatabase(sqliteInput)
	if err != nil {
		logger.With("error", err).Fatal("could not setup database")
	}

	restPort := config.REST.Port
	restServerInput := &rest.RESTServerInput{
		Address:   restPort,
		Logger:    logger,
		Datastore: database,
	}
	restServer, _ := rest.NewRESTServer(restServerInput)

	if err := restServer.Start(); err != nil {
		logger.With("error", err).Fatal("could not start HTTP server")
	}

	if err := restServer.Stop(); err != nil {
		logger.With("error", err).Fatal("could not stop HTTP server")
	}
}
