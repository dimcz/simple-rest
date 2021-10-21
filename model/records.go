package model

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"simple-rest/logging"
	"simple-rest/settings"
)

type Records struct {
	ID    int    `json:"id,omitempty" gorm:"primary_key"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var db *gorm.DB
var logger *logging.Logger

func Setup(migration bool, l *logging.Logger) {
	logger = l
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		settings.AppSettings.Postgres.Host,
		settings.AppSettings.Postgres.Username,
		settings.AppSettings.Postgres.Password,
		settings.AppSettings.Postgres.DBName,
		settings.AppSettings.Postgres.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("model.Setup error: %v", err)
	}

	if migration {
		if err = db.AutoMigrate(new(Records)); err != nil {
			logger.Fatalf("model.AutoMigrate error: %v", err)
		}
	}
}

func SelectAll() (records []Records, err error) {
	err = db.Find(&records).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return records, nil
}

func SelectRecordByID(id interface{}) (*Records, error) {
	var record Records
	err := db.Find(&record).Where("id = ?", id).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &record, nil
}

func CreateUser(record Records) (id int, err error) {
	if err := db.Create(&record).Error; err != nil {
		logger.Error(err)
		return 0, err
	} else {
		return record.ID, nil
	}
}

func DeleteUser(id interface{}) error {
	err := db.Where("id = ?", id).Delete(&Records{}).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func UpdateUser(id interface{}, record Records) error {
	err := db.Model(&Records{}).Where("id = ?", id).Updates(record).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}
