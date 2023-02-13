package wasp

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Database interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (interface{}, error)
}

type BaseModel interface {
	GetTable() string
	GetId() int64
	GetRedisDb() int64
	GetRedisKey() string
}

type BasePersistenceModal struct {
	ID int64 `gorm:"type:int;primary_key"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Initializing default
	//	var mode string
	limit := 2
	page := 1
	sort := "created_at asc"
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break

		}
	}
	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

}
