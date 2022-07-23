package services

import (
	"fmt"
	"log"
	"photoservice/backend/api/apiModel"
	dbModel "photoservice/backend/models"
)

func (db *DatabaseService) NewTaskFile(info apiModel.FileData) []apiModel.ErrorJSON {
	task, apiError := db.GetUploadInfo(info.FileName, info.UserID, info.RaceID)
	fmt.Println("TAAAAAAAAAAAAAAAAAAAAAASK",info.FileName,info.UserID,info.RaceID,info.SID)
	if task != nil || apiError != nil {
		fmt.Println(apiError, task)
		return nil
	}
	result := db.db.Create(&dbModel.UploadInfo{FileName: info.FileName, UserID: info.UserID, SID: info.SID, RaceUID: info.RaceID})
	if result.Error != nil {
		log.Println("LOG_ERROR: DB DELETE ERROR", result.Error.Error())
		fmt.Println(result.Error.Error())
		return []apiModel.ErrorJSON{{Classification: "DB", Message: result.Error.Error()}}
	}
	return nil
}

func (db *DatabaseService) UpdateTaskInfo(meta apiModel.UploadMeta) (*dbModel.UploadInfo, []apiModel.ErrorJSON) {
	var uploadInfo dbModel.UploadInfo
	fmt.Println(meta.PhotoName, meta.UserID)
	if meta.UserID == "" || meta.PhotoName == "" {
		log.Println("LOG_ERROR: DB FIND ERROR")
		return nil, nil
	}
	result := db.db.Where(&dbModel.UploadInfo{FileName: meta.PhotoName, UserID: meta.UserID, RaceUID: meta.RaceID}).Find(&uploadInfo)
	if result.Error != nil {
		log.Println("LOG_ERROR: DB FIND ERROR", result.Error.Error())
		fmt.Println(result.Error.Error())
		return nil, nil
	} else if result.RowsAffected == 0 {
		log.Println("LOG_ERROR: Task don't exist")
		fmt.Println("TASK DOESN'T EXIST")
		return nil, nil
	} else {
		uploadInfo.Width = meta.ImageWidth
		uploadInfo.Height = meta.ImageHeight
		uploadInfo.ImageSize = meta.ImageSize
		uploadInfo.ResizedSize = meta.ResizedSize
		uploadInfo.MLSize = meta.MLSize
		uploadInfo.WatermarkSize = meta.WatermarkSize
		return &uploadInfo, nil
	}
}

func (db *DatabaseService) GetUploadInfo(Filename string, UserID string, RaceID string) (*dbModel.UploadInfo, []apiModel.ErrorJSON) {
	var uploadInfo dbModel.UploadInfo
	result := db.db.Where(&dbModel.UploadInfo{FileName: Filename, UserID: UserID, RaceUID: RaceID}).Find(&uploadInfo)
	if result.Error != nil {
		log.Println("LOG_ERROR: DB FIND ERROR", result.Error.Error())
		fmt.Println(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		fmt.Println("NIIIIIIIIIIIIIIIIIIIIIIIIIIL",Filename,UserID,RaceID)
		return nil, nil
	}
	return &uploadInfo, nil
}

func (db *DatabaseService) GetUploadInfoByToken(Token string) (*dbModel.UploadInfo, []apiModel.ErrorJSON) {
	var uploadInfo dbModel.UploadInfo
	result := db.db.Where(&dbModel.UploadInfo{Token: Token}).Find(&uploadInfo)
	if result.Error != nil {
		log.Println("LOG_ERROR: DB FIND ERROR", result.Error.Error())
		fmt.Println(result.Error.Error())
	}
	if result.RowsAffected == 0 {
		fmt.Println("NIIIIIIIIIIIIIIIIIIIIIIIIIIL",Token)
		return nil, nil
	}
	return &uploadInfo, nil
}

func (db *DatabaseService) SaveUploadInfo(info dbModel.UploadInfo) {
	db.db.Save(&info)
}
