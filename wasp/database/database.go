package database

import (
	"context"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	ctx         context.Context
	RedisClient RedisClient
	Gorm        *gorm.DB
}

func NewDatabase() (*Database, error) {
	redisClient, err := newRedisClient()
	if err != nil {
		return nil, err
	}

	return &Database{
		ctx:         context.Background(),
		Gorm:        openGorm(),
		RedisClient: redisClient,
	}, nil
}

func openGorm() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		// os.Getenv("POSTGRES_DB"),
	)

	// dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
	dbClient, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return dbClient
}
