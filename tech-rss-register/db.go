package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Subscription struct {
	URL         string
	ChannelName string
}

func setupDB() (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(sqlite.Open("../data/prod.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Subscription{})
	return db, nil
}

// URLがDBになければ作る。
func CreateSubscription(db *gorm.DB, url string, channelName string) {
	var count int64
	tx := db.Begin()
	tx.Model(&Subscription{}).Where("url = ?", url).Count(&count)
	if count == 0 {
		tx.Create(&Subscription{
			URL:         url,
			ChannelName: channelName,
		})
	}
	tx.Commit()
}
