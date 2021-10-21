package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"simple-rest/settings"
)

type Records struct {
	ID    int    `json:"id,omitempty" gorm:"primary_key"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var db *gorm.DB

func Setup(migration bool) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		settings.AppSettings.Postgres.Host,
		settings.AppSettings.Postgres.Username,
		settings.AppSettings.Postgres.Password,
		settings.AppSettings.Postgres.DBName,
		settings.AppSettings.Postgres.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("model.Setup error: %v", err)
	}

	if migration {
		if err = db.AutoMigrate(new(Records)); err != nil {
			log.Fatalf("model.AutoMigrate error: %v", err)
		}
	}
}

func SelectAll(c *gin.Context) {
	var records []Records
	err := db.Find(&records).Error
	if err == gorm.ErrRecordNotFound {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, records)
}

func SelectRecordByID(c *gin.Context) {
	id := c.Param("record")
	var record Records

	err := db.Find(&record).Where("id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, record)
}

func CreateUser(c *gin.Context) {
	var record Records
	if err := c.BindJSON(&record); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.Create(&record).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, record)
	}
}

func DeleteUser(c *gin.Context) {
	id := c.Param("record")
	if err := db.Where("id = ?", id).Delete(&Records{}).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		c.Writer.WriteHeader(http.StatusNoContent)
	}

}

func UpdateUser(c *gin.Context) {
	id := c.Param("record")

	var record Records
	if err := c.BindJSON(&record); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.Model(&Records{}).Where("id = ?", id).Updates(record).Error; err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}
