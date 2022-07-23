package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"photoservice/backend/api/apiModel"
	"time"
)

type Photo struct {
	gorm.Model
	PUID      string `gorm:"column:UID;type:uuid;default:uuid_generate_v4()"`
	PhotoName string `gorm:"column:photo_id;unique"`
	//Competitors *string   `gorm:"column:competitors"`
	Size       int       `gorm:"column:size;type:bigint"`
	Height     int       `gorm:"column:height"`
	Width      int       `gorm:"column:width"`
	LoadedAt   time.Time `gorm:"column:loaded_at"`
	RaceID     string    `gorm:"column:race_id;index:idx_race"`
	UserID     string    `gorm:"column:user_id"`
	IsActive   bool      `gorm:"column:is_active;default:false"`
	IsDetected bool      `gorm:"column:is_detected;default:false"`
}

type PhotoTag struct {
	gorm.Model
	PhotoID uint   `gorm:"column:photo_id;index:idx_photoID"`
	RaceID  string `gorm:"column:race_id;index:idx_tag,priority:1"`
	Photo   Photo  `gorm:"foreignKey:PhotoID"`
	Tag     string `gorm:"column:tag;index:idx_tag,priority:2"`
	UserTag bool   `gorm:"column:user_tag"`
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

func (uploadInfo *UploadInfo) AfterUpdate(tx *gorm.DB) (err error) {
	if uploadInfo.MLUploaded && uploadInfo.MainUploaded && uploadInfo.ResizedUploaded && uploadInfo.WatermarkUploaded {
		var photo Photo
		tx.Model(&Photo{}).Where(&Photo{PUID: uploadInfo.Token}).Find(&photo)
		photo.IsActive = true
		tx.Save(&photo)
		tx.Delete(&uploadInfo)
	}
	return
}

func (photo *Photo) AfterUpdate(tx *gorm.DB) (err error) {
	if (!photo.IsDetected) && (photo.IsActive) {
		if err := tx.Where(&PhotoTag{PhotoID: photo.ID,UserTag: false}).Unscoped().Delete(&PhotoTag{}).Error;err!=nil{
			return err
		}
		taskBatch := apiModel.TaskBatch{Data: []apiModel.Task{{fmt.Sprintf("https://api-dip.duckdns.org/file/get?UUID=%s&type=ml&token=iojqwejknasdfjknxcvasdfkqwermn", photo.PUID), photo.PUID}}}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(taskBatch)
		if err != nil {
			log.Println("LOG_ERROR: JSON ENCODE ERROR", err.Error())
			fmt.Println(err)
		}
		mlReq, err := http.NewRequest("POST", "https://dip.tujhs.rocks/api/tasks/new", &buf)
		mlReq.Header.Set("content-type", "application/json")
		if err != nil {
			log.Println("LOG_ERROR: HTTP POST ERROR", err.Error())
			fmt.Println(err)
		}
		mlReq.Header.Set("Authorization", "OANQ7uVfrSSZWzrE6g7qZSGPJdP/g/8j8fFV5GE5IgFU4IewqJ8EVa4mqpthgWJ0W0RkK6ZvuWRzOEw975FxAFXo3yTb0LPFS36YcoPR5APNa0z8E7eL2RG1GD6Xs4dOcESihgP1o6ObUH9CfG7T42coPP2XRRaCkqnGm183MqU=")
		_, _ = http.DefaultClient.Do(mlReq)
	}
	return
}
