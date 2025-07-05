package infrastructure

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB() (*gorm.DB, error) {
	log := logrus.New()

	log.Info("Initializing PostgreSQL database connection")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_SSLMODE"),
	)

	log.WithFields(logrus.Fields{
		"host":     viper.GetString("DB_HOST"),
		"port":     viper.GetString("DB_PORT"),
		"user":     viper.GetString("DB_USER"),
		"database": viper.GetString("DB_NAME"),
		"sslmode":  viper.GetString("DB_SSLMODE"),
	}).Debug("Database connection parameters")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to connect to PostgreSQL database")
		return nil, err
	}

	log.Info("Successfully connected to PostgreSQL database")

	sqlDB, err := db.DB()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get underlying sql.DB")
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to ping PostgreSQL database")
		return nil, err
	}

	log.Info("Database connection ping successful")

	return db, nil
}
