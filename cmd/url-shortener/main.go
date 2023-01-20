package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-shortner/internal/config"
	"url-shortner/internal/datastore"
	"url-shortner/internal/logger"
	"url-shortner/internal/rest"
)

func RunApplication() error {
	filesystem := config.NewFilesystem()
	config := config.NewConfiguration("config.yaml", filesystem)

	logLevel, err := config.GetString("logger.level")
	logEncoding, err := config.GetString("logger.encoding")
	logger, err := logger.NewZapLogger(logEncoding, logLevel)
	if err != nil {
		log.Print(err)
	}

	ctx := context.Background()
	dbHost, err := config.GetString("database.host")
	dbPort, err := config.GetInt("database.port")
	dbUser, err := config.GetString("database.user")
	dbPassword, err := config.GetString("database.password")
	dbName, err := config.GetString("database.name")
	dbCapacity, err := config.GetInt("database.capacity")
	if err != nil {
		logger.WithError(err).Error("could not get database config")
	}

	database, closeFunc := datastore.NewPostgresUrlStore(
		ctx, logger, dbUser, dbPassword, dbHost, dbName, dbPort, dbCapacity,
	)
	defer closeFunc(ctx)

	restHost, err := config.GetString("rest.host")
	restPort, err := config.GetInt("rest.port")
	if err != nil {
		logger.WithError(err).Error("could not get rest config")
		return nil
	}

	router := rest.NewRouter(logger, database)

	restServer := rest.NewRESTServer(
		logger, database, restHost, restPort, router,
	)

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)

	go restServer.Start()

	<-terminationChannel

	if err := restServer.Stop(); err != nil {
		return err
	}

	if err := logger.Sync(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := RunApplication(); err != nil {
		os.Exit(1)
	}
}
