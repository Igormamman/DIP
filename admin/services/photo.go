package services

import (
	"admin/api/apiModel"
	dbModel "admin/models"
	"gorm.io/gorm"
)

func (db *DatabaseService) GetPhotosByOffset(limit int, offset int) ([]dbModel.Photo, []apiModel.ErrorJSON) {
	var photos []dbModel.Photo
	var result *gorm.DB
	if limit < 0 {
		result = db.db.Offset(offset).Find(&photos)
	} else {
		result = db.db.Limit(limit).Offset(offset).Find(&photos)
	}
	if result.Error != nil {
		return nil, []apiModel.ErrorJSON{
			{Classification: "DB", Message: result.Error.Error()},
		}
	}
	return photos, nil
}
