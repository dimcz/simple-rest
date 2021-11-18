package data

import (
	"fmt"
	"simple-rest/pkg/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewConnection(config *util.Config, logger *util.Logger) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("database.NewConnection error: %v", err)
	}

	if err = db.AutoMigrate(new(User), new(Record), new(Number)); err != nil {
		logger.Fatalf("database.AutoMigrate error: %v", err)
	}
}
