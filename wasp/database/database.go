package database

import (
	"context"
	"fmt"
	"os"

	"github.com/berachain/stargazer/wasp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
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
		Gorm:        OpenGorm(),
		RedisClient: redisClient,
	}, nil
}

func OpenGorm() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	dbClient, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./wasp/query", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(dbClient)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(models.EthBlockModel{}, models.TransactionModel{}, &models.EthTxnReceipt{}, &models.EthLog{})

	// Execute the generator
	g.Execute()

	dbClient.AutoMigrate(&models.EthBlockModel{}, &models.TransactionModel{}, &models.EthTxnReceipt{}, &models.EthLog{})
	if err != nil {
		panic(err)
	}

	return dbClient
}

func (db *Database) Get(r *GetRequest, gormFunc func() ([]byte, error)) ([]byte, error) {
	data, err := db.RedisClient.Get(r.Key)
	if err != nil || data == nil {
		fmt.Print("\nI AM IN PSQL\n")
		data, err := gormFunc()
		redisErr := db.RedisClient.Set(r.Key, data)
		if redisErr != nil {
			fmt.Print("shes fked bud")
		}
		return data, err

	}
	return data, err
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
