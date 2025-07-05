package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/edumes/golang-api-rest/seeds"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := infrastructure.GetColoredLogger()

	logger.Info("Starting Seeds CLI")

	var seedType = flag.String("type", "all", "Type of seed to run (all, users, projects, project-items)")
	flag.Parse()

	logger.Info("Loading configuration")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Failed to read .env file, using environment variables")
	}
	viper.AutomaticEnv()

	logger.WithFields(logrus.Fields{
		"db_host": viper.GetString("DB_HOST"),
		"db_port": viper.GetString("DB_PORT"),
		"db_name": viper.GetString("DB_NAME"),
	}).Info("Configuration loaded successfully")

	logger.Info("Initializing database connection")
	db, err := infrastructure.NewPostgresDB()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to connect to database")
	}

	logger.Info("Database connection established successfully")

	seeder := seeds.NewSeeder(db)

	ctx := context.Background()

	switch *seedType {
	case "all":
		logger.Info("Running all seeds")
		if err := seeder.RunAll(ctx); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to run all seeds")
		}
	case "users":
		logger.Info("Running user seeds")
		if err := seeder.RunUsers(ctx); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to run user seeds")
		}
	case "projects":
		logger.Info("Running project seeds")
		if err := seeder.RunProjects(ctx); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to run project seeds")
		}
	case "project-items":
		logger.Info("Running project item seeds")
		if err := seeder.RunProjectItems(ctx); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Failed to run project item seeds")
		}
	default:
		logger.WithFields(logrus.Fields{
			"seed_type": *seedType,
		}).Fatal("Invalid seed type")
	}

	logger.Info("Seeds completed successfully")
	fmt.Println("Seeds completed successfully!")
}
