package models

import "github.com/jinzhu/gorm"

type Keyword struct {
	gorm.Model
	Label  string  `json:"label" gorm:"not null;unique"`
	Tweets []Tweet `json:"-"`
}

func CreateKeyword(db *gorm.DB, keyword *Keyword) error {
	if err := db.Create(&keyword).Error; err != nil {
		return err
	}
	return nil
}

func GetKeywords(db *gorm.DB) ([]Keyword, error) {
	keywords := []Keyword{}
	db.Preload("Tweets").Find(&keywords)
	return keywords, nil
}
