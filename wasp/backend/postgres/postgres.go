package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	client *gorm.DB
}

func New(postgresHost string, user string, password string, dbName string, port int, sslMode string, timeZone string) Postgres {
	dbUri := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s TimeZone=%s", postgresHost, port, user, dbName, sslMode, password, timeZone)
	dbClient, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	postgres := Postgres{client: dbClient}
	return postgres
}

func NewFromEnv() Postgres {
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

	postgres := Postgres{client: dbClient}
	return postgres
}
