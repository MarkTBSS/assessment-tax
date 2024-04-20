package databases

import (
	"fmt"
	"log"
	"sync"

	"github.com/MarkTBSS/assessment-tax/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	*gorm.DB
}

func (db *postgresDatabase) Connect() *gorm.DB {
	return postgresDatabaseInstace.DB
}

var (
	postgresDatabaseInstace *postgresDatabase
	once                    sync.Once
)

func NewPostgresDatabase(conf *config.Database) Database {
	once.Do(func() {
		dsn := fmt.Sprint(conf.DatabaseURL)
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		log.Printf("Connection String : %s", conf.DatabaseURL)
		postgresDatabaseInstace = &postgresDatabase{conn}
	})
	return postgresDatabaseInstace
}
