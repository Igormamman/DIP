package services

import (
	"fmt"
	"log"
	"net/http"
	"photoservice/backend/api/apiModel"
	dbModel "photoservice/backend/models"
	"time"
)

func SetupWebhook () (error){
	time.Sleep(time.Second*2)
	fmt.Println("SETUP WEBHOOK")
	log.Println("LOG_INFO: SETUP WEBHOOK")
	req, _ := http.NewRequest("GET",fmt.Sprintf("https://dip.tujhs.rocks/api/client/setWebhook?url=%s","https://api-dip.duckdns.org/ml-webhook"),nil)
	req.Header.Set("Authorization","OANQ7uVfrSSZWzrE6g7qZSGPJdP/g/8j8fFV5GE5IgFU4IewqJ8EVa4mqpthgWJ0W0RkK6ZvuWRzOEw975FxAFXo3yTb0LPFS36YcoPR5APNa0z8E7eL2RG1GD6Xs4dOcESihgP1o6ObUH9CfG7T42coPP2XRRaCkqnGm183MqU=")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("RECOVERING FROM ERROR")
			log.Println("LOG_WARNING: RECOVERING FROM ERROR")
		}
	}();
	resp, err := http.DefaultClient.Do(req)
	if (err!=nil){
		fmt.Println("WEBHOOK NOT SETUP",resp.Status)
		log.Println("LOG_ERROR: WEBHOOK NOT SETUP, response status:",resp.Status)
		return err
	}else{
		fmt.Println("WEBHOOK SUCCESSFUL SETUP ",resp.Status)
		log.Println("LOG_INFO: WEBHOOK SUCCESSFUL SETUP, response status:",resp.Status)
		return nil
	}
}

func SaveResult (taskResult apiModel.TaskResult) (){
	dbPhoto, err := DB.GetPhotoByPUID(taskResult.Token)
	if err != nil || dbPhoto == nil {
		fmt.Println("DB TOKEN FIND ERROR ")
		log.Println("LOG_WARNING: DB TOKEN FIND ERROR, ERROR:",err)
		return
	}
	for _,tag:=range taskResult.DetectedNums {
		tag := dbModel.PhotoTag{PhotoID: dbPhoto.ID, Tag: tag, UserTag: false,RaceID: dbPhoto.RaceID}
		DB.db.Save(&tag)
	}
	dbPhoto.IsDetected = true
	DB.db.Save(dbPhoto)
	return
}

func (db* DatabaseService) GetPhotoTags(photo dbModel.Photo) []string{
	var tags []dbModel.PhotoTag
	var result []string
	DB.db.Where(&dbModel.PhotoTag{PhotoID: photo.ID}).Find(&tags)
	for _,tag := range tags{
		result = append(result,tag.Tag)
	}
	return result
}
