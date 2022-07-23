package models

import (
	"gorm.io/gorm"
	"time"
)


type Photo struct {
	gorm.Model
	PUID        string    `gorm:"column:UID;type:uuid;default:uuid_generate_v4()"`
	PhotoName   string    `gorm:"column:photo_id;unique"`
	//Competitors *string   `gorm:"column:competitors"`
	Size        int       `gorm:"column:size;type:bigint"`
	Height      int       `gorm:"column:height"`
	Width       int       `gorm:"column:width"`
	LoadedAt    time.Time `gorm:"column:loaded_at"`
	RaceID      string    `gorm:"column:race_id;index:idx_race"`
	UserID      string    `gorm:"column:user_id"`
	IsActive    bool      `gorm:"column:is_active;default:false"`
	IsDetected  bool      `gorm:"column:is_detected;default:false"`
}

type PhotoTag struct {
	gorm.Model
	PhotoID uint   `gorm:"column:photo_id;index:idx_tag,priority:2;index:idx_photoID"`
	Photo   Photo `gorm:"foreignKey:PhotoID"`
	Tag     string   `gorm:"column:tag;index:idx_tag,priority:1"`
	UserTag bool  `gorm:"column:user_tag"`
}

type UploadInfo struct {
	gorm.Model
	FileName          string `gorm:"column:filename"`
	UserID            string `gorm:"column:user_id"`
	RaceUID           string `gorm:"column:race_id"`
	Width             int    `gorm:"column:width"`
	Height            int    `gorm:"column:height"`
	ImageSize         int    `gorm:"column:image_size;default:0" `
	ResizedSize       int    `gorm:"column:resized_size;default:0"`
	MLSize            int    `gorm:"column:ml_size;default:0"`
	WatermarkSize     int    `gorm:"column:watermark_size;default:0"`
	MainUploaded      bool   `gorm:"column:main_uploaded;default:false"`
	ResizedUploaded   bool   `gorm:"column:resized_uploaded;default:false"`
	MLUploaded        bool   `gorm:"column:ml_uploaded;default:false"`
	WatermarkUploaded bool   `gorm:"column:watermark_uploaded;default:false"`
	SID               string `gorm:"column:sid"`
	Token             string `gorm:"column:token"`
}