package models

import (
	"time"

	"gorm.io/gorm"
)

type UserCookie struct {
	gorm.Model
	Cookie string `gorm:"column:Cookie"`
	PUID   string `gorm:"column:UID"`
}

type UserAccess struct {
	gorm.Model
	Cookie                  string `gorm:"column:Cookie"`
	PUID                    string `gorm:"column:UID"`
	RaceUID                 string `gorm:"column:RaceUID"`
	IsOrg                   bool   `gorm:"column:isOrg"`
	IsComp                  bool   `gorm:"column:isComp"`
	IsMedia                 bool   `gorm:"column:isMedia"`
	IsMediaConfirmed        bool   `gorm:"column:isMediaConfirmed"`
	IsRegistrationConfirmed bool   `gorm:"column:isMediaConfirmed"`
}

type Race struct {
	gorm.Model
	RaceID string    `gorm:"race_id;unique"`
	Name   string    `gorm:"name"`
	Date   time.Time `gorm:"date"`
	City   string    `gorm:"city"`
}
