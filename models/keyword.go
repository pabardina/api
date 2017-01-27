package models

import (
	"github.com/jinzhu/gorm"
)

type Keyword struct {
	gorm.Model
	Label  string `json:"label,omitempty" gorm:"not null;unique"`
	Tweets []Tweet
}

func CreateKeyword(db *gorm.DB, keyword *Keyword) error {
	db.NewRecord(keyword)
	db.Create(&keyword)
	return nil
}

func GetKeywords(db *gorm.DB) ([]Keyword, error) {
	keywords := []Keyword{}
	db.Preload("Tweets").Find(&keywords)
	return keywords, nil
}
