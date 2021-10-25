package data

import "gorm.io/gorm"

type User struct {
	ID       int      `json:"id" gorm:"primary_key"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Records  []Record `json:"records" gorm:"constraint:OnDelete:CASCADE"`
}

func CheckAuth(username, password string) (int, error) {
	var user User
	err := db.Select("id").Where(User{Username: username, Password: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	if user.ID > 0 {
		return user.ID, nil
	}

	return 0, nil
}
