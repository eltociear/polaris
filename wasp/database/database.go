package database

import (
	"context"
	"fmt"
	"os"

	"github.com/berachain/stargazer/wasp/utils"
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

func (db *Database) Get(r *GetRequest, gormFunc func() ([]byte, error)) ([]byte, error) {
	data, err := db.RedisClient.Get(r.Key)
	if err != nil {
		byteData, err := gormFunc()
		if err != nil {
			panic(err)
		}
		return byteData, nil
	}
	byteData, err := utils.GetBytes(data)
	if err != nil {
		panic(err)
	}
	return byteData, nil
}

func (db *Database) Set(r *SetRequest, gormFunc func() error) error {
	err := gormFunc()
	if err != nil {
		return err
	}
	err = db.RedisClient.Set(r.Key, r.Value)
	return err
}

type GetRequest struct {
	RedisDb int64
	Key     string
}

type SetRequest struct {
	RedisDb int64
	Key     string
	Value   []byte
}
