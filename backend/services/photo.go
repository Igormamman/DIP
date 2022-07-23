package services

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"photoservice/backend/api/apiModel"
	dbModel "photoservice/backend/models"
	"time"

	"gorm.io/gorm"
)

func (db *DatabaseService) CreateNewPhoto(meta dbModel.UploadInfo) (*dbModel.Photo, []apiModel.ErrorJSON) {
	log.Println("LOG_INFO: CREATING NEW DB PHOTO")
	if meta.Token != "" {
		photo, err := db.GetPhotoByPUID(meta.Token)
		if photo != nil && err == nil {
			fmt.Println("Photo exists")
			return photo, []apiModel.ErrorJSON{{Classification: "DB", Message: "This Photo exist in DB"}}
		}
	}
	dbImageName := fmt.Sprintf("%d-%s", time.Now().Unix(), meta.FileName)
	photo := dbModel.Photo{
		PhotoName:  dbImageName,
		Height:     meta.Height,
		Width:      meta.Width,
		Size:       meta.ImageSize,
		LoadedAt:   time.Now(),
		RaceID:     meta.RaceUID,
		UserID:     meta.UserID,
		IsActive:   false,
		IsDetected: false,
	}
	result := db.db.Create(&photo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB CREATE ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{{Classification: "DB", Message: result.Error.Error()}}
	}
	return &photo, nil
}

func (db *DatabaseService) DeletePhoto(photo dbModel.Photo, unscoped bool) []apiModel.ErrorJSON {
	var result *gorm.DB
	year, month, day := photo.LoadedAt.Date()
	filepath := fmt.Sprintf("%d/%d/%d/%s", year, month, day, photo.PhotoName)
	fmt.Println(filepath)
	if unscoped {
		deleteObjs := MinioClient.ListObjects(context.Background(), Cfg.BucketName, minio.ListObjectsOptions{Recursive: true, Prefix: filepath})
		removeError := MinioClient.RemoveObjects(context.Background(), Cfg.BucketName, deleteObjs, minio.RemoveObjectsOptions{})
		err := <-removeError
		if err.Err == nil {
			if tagErr := db.db.Where(&dbModel.PhotoTag{PhotoID: photo.ID, UserTag: false}).Unscoped().Delete(&dbModel.PhotoTag{}).Error; tagErr != nil {
				return []apiModel.ErrorJSON{{Classification: "DB", Message: tagErr.Error()}}
			}
			result = db.db.Unscoped().Delete(&photo)
			if result.Error != nil {
				log.Println("LOG_ERROR: DB DELETE ERROR", result.Error.Error())
				fmt.Println(result.Error.Error())
				return []apiModel.ErrorJSON{{Classification: "DB", Message: result.Error.Error()}}
			} else {
				return nil
			}
		} else {
			return []apiModel.ErrorJSON{{Classification: "S3", Message: err.Err.Error()}}
		}
	} else {
		if tagErr := db.db.Where(&dbModel.PhotoTag{PhotoID: photo.ID, UserTag: false}).Delete(&dbModel.PhotoTag{}).Error; tagErr != nil {
			return []apiModel.ErrorJSON{{Classification: "DB", Message: tagErr.Error()}}
		}
		result = db.db.Delete(&photo)
		if result.Error != nil {
			log.Println("LOG_ERROR: DB DELETE ERROR", result.Error.Error())
			fmt.Println(result.Error.Error())
			return []apiModel.ErrorJSON{{Classification: "DB", Message: result.Error.Error()}}
		} else {
			return nil
		}
	}
}

func (db *DatabaseService) GetPhotoByName(name string) (*dbModel.Photo, []apiModel.ErrorJSON) {
	var photo dbModel.Photo
	result := db.db.Where(dbModel.Photo{PhotoName: name}).First(&photo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return &photo, nil
}

func (db *DatabaseService) GetPhotoByID(photoID int) (*dbModel.Photo, []apiModel.ErrorJSON) {
	var photo dbModel.Photo
	result := db.db.Where("id = ?", photoID).First(&photo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return &photo, nil
}

func (db *DatabaseService) GetPrevPUID(photoID int, count int) ([]string, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	var uuids []string
	result := db.db.Where("id<?", photoID).Order("id desc").Limit(count).Find(&photos)
	if result.Error == nil && len(photos) > 0 {
		for _, photo := range photos {
			uuids = append(uuids, photo.PUID)
		}
		return uuids, nil
	} else {
		return nil, []apiModel.ErrorJSON{{Classification: "db", Message: result.Error.Error()}}
	}
}

func (db *DatabaseService) GetNextPUID(photoID int, count int) ([]string, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	var uuids []string
	result := db.db.Where("id>?", photoID).Order("id asc").Limit(count).Find(&photos)
	if result.Error == nil && len(photos) > 0 {
		for _, photo := range photos {
			uuids = append(uuids, photo.PUID)
		}
		return uuids, nil
	} else {
		return nil, []apiModel.ErrorJSON{{Classification: "db", Message: result.Error.Error()}}
	}
}

func (db *DatabaseService) GetPhotoByPUID(UID string) (*dbModel.Photo, []apiModel.ErrorJSON) {
	var photo dbModel.Photo
	result := db.db.Where(&dbModel.Photo{PUID: UID}).First(&photo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return &photo, nil
}

func (db *DatabaseService) GetPhotoCount(competitor string, raceUID string, detected bool) int64 {
	if detected {
		if (competitor == "") && (raceUID != "") {
			result := db.db.Where(dbModel.Photo{RaceID: raceUID}).Joins("left join photo_tags on photo_tags.photo_id = photos.id").
				Where("photo_tags.photo_id IS NOT NULL").Where("photo_tags.deleted_at is NULL").Distinct().Find(&[]dbModel.Photo{})
			return result.RowsAffected
		} else if (competitor != "") && (raceUID != "") {
			result := db.db.Where(dbModel.PhotoTag{RaceID: raceUID, Tag: competitor}).Find(&[]dbModel.PhotoTag{})
			return result.RowsAffected
		}
	} else {
		result := db.db.Where(dbModel.Photo{RaceID: raceUID}).Joins("left join photo_tags on photo_tags.photo_id = photos.id").
			Where("photo_tags.photo_id IS NULL").Where("photo_tags.deleted_at is NULL").Find(&[]dbModel.Photo{})
		return result.RowsAffected
	}
	return 0
}

func (db *DatabaseService) GetUserPhotoCount(userID string, raceUID string) int64 {
	result := db.db.Where(dbModel.Photo{RaceID: raceUID, UserID: userID}).Find(&[]dbModel.Photo{})
	return result.RowsAffected
}

func (db *DatabaseService) GetDetectedPhotoFromDB(competitor string, raceUID string, limit int, offset int) ([]dbModel.Photo, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	var result *gorm.DB
	if (competitor == "") && (raceUID != "") {
		result = db.db.Where(dbModel.Photo{RaceID: raceUID, IsDetected: true}).Joins("left join photo_tags on photo_tags.photo_id = photos.id").
			Where("photo_tags.photo_id IS NOT NULL").Where("photo_tags.deleted_at is NULL").Distinct().Limit(limit).Offset(offset).Find(&photos)
	} else if (competitor != "") && (raceUID != "") {
		var photoTags []dbModel.PhotoTag
		result = db.db.Preload("Photo").Where(dbModel.PhotoTag{RaceID: raceUID, Tag: competitor}).Where("photo_tags.deleted_at is NULL").Limit(limit).Offset(offset).Find(&photoTags)
		for _, tag := range photoTags {
			fmt.Println(tag.RaceID, tag.ID, tag.Photo.ID)
			photos = append(photos, tag.Photo)
		}
	}
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return photos, nil
}

func (db *DatabaseService) GetUndetectedPhotoFromDB(raceUID string, limit int, offset int) ([]dbModel.Photo, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	var result *gorm.DB
	fmt.Println("RACE", raceUID)
	if raceUID != "" {
		result = db.db.Where(dbModel.Photo{RaceID: raceUID}).Joins("left join photo_tags on photo_tags.photo_id = photos.id").
			Where("photo_tags.photo_id IS NULL").Where("photo_tags.deleted_at is NULL").Distinct().Limit(limit).Offset(offset).Find(&photos)
	}
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return photos, nil
}

func DeletePhoto(dbImage dbModel.Photo) []apiModel.ErrorJSON {
	error := DB.DeletePhoto(dbImage, false)
	if error != nil {
		fmt.Println(error)
		log.Println("LOG_ERROR: DB DELETE ERROR", error)
		return error
	}
	return nil
}

func (db *DatabaseService) GetPhotosByOffset(limit int, offset int) ([]dbModel.Photo, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	result := db.db.Limit(limit).Offset(offset).Find(&photos)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return photos, nil
}

func (db *DatabaseService) GetUserPhotosByOffset(limit int, offset int, userID string, raceID string) ([]dbModel.Photo, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	result := db.db.Where(dbModel.Photo{UserID: userID, RaceID: raceID}).Limit(limit).Offset(offset).Find(&photos)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		log.Println("LOG_ERROR: DB GET ERROR", result.Error.Error())
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return photos, nil
}

func (db *DatabaseService) SavePhoto(photo dbModel.Photo) {
	DB.db.Save(&photo)
}

/*func (db *DatabaseService) AppendTag(tag string, dbImage dbModel.Photo) []apiModel.ErrorJSON {
	competitors := dbImage.Competitors
	if validateTag(tag) {
		sb := strings.Builder{}
		if competitors != nil {
			sb.WriteString(*competitors)
			sb.WriteString(tag)
			sb.WriteString(",")
			*dbImage.Competitors = sb.String()
		} else {
			sb.WriteString(",")
			sb.WriteString(tag)
			sb.WriteString(",")
		}
	}
	DB.db.Save(dbImage)
	return nil
}*/

/*func (db *DatabaseService) DeleteTag(tag string, dbImage dbModel.Photo) []apiModel.ErrorJSON {
	competitors := dbImage.Competitors
	if competitors != nil {
		*competitors = strings.Replace(*competitors, fmt.Sprintf(",%s,", tag), ",", 1)
		if *competitors == "," || competitors == nil {
			competitors = nil
		}
		*dbImage.Competitors = *competitors
	} else {
		return nil
	}
	DB.db.Save(dbImage)
	return nil
}*/

func validateTag(tag string) bool {
	return true
}
