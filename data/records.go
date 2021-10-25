package data

import "gorm.io/gorm"

type Record struct {
	ID        int      `json:"id" gorm:"primary_key"`
	UserID    int      `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Numbers   []Number `json:"numbers" gorm:"constraint:OnDelete:CASCADE"`
}

type Number struct {
	ID       int    `json:"id" gorm:"primary_key"`
	RecordID int    `json:"record_id"`
	Number   string `json:"number"`
}

func Records(userId int) (records []Record, err error) {
	err = db.Model(&Record{}).Preload("Numbers").Where("user_id = ?", userId).Find(&records).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return records, nil
}

func DeleteRecordByID(userId int, recordId string) error {
	return db.Delete(&Record{}, "user_id = ? and id = ?", userId, recordId).Error
}
