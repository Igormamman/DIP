package api

import (
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"photoservice/backend/api/apiModel"
	"photoservice/backend/services"
)

func webhookRouter (ctx *macaron.Context, taskResponse apiModel.ClientResultResponse){
	log.Println("LOG_INFO: WEBHOOK RESPONSE",taskResponse)
	for _,task := range taskResponse.Data{
		services.SaveResult(task)
	}
	ctx.Resp.WriteHeader(http.StatusOK)
	return
}


// task creates via hook on photo update in models.db file
func createTask (ctx *macaron.Context, fileInfo apiModel.NewFile){
	dbImage, _ := services.DB.GetPhotoByPUID(fileInfo.Token)
	if dbImage != nil{
		dbImage.IsDetected=false;
		services.DB.SavePhoto(*dbImage)
/*
		services.DB.DeleteMlTags(*dbImage)
		taskBatch := apiModel.TaskBatch{Data: []apiModel.Task{{fmt.Sprintf("https://api-dip.duckdns.org/file/get/%s/%s?type=ml&token=iojqwejknasdfjknxcvasdfkqwermn", fileNamePrefix, dbImage.PhotoName), dbImage.PUID}}}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(taskBatch)
		if err != nil {
			log.Println("LOG_ERROR: JSON ENCODE ERROR",err.Error())
			fmt.Println(err)
		}
		mlReq, err := http.NewRequest("POST", "https://dip.tujhs.rocks/api/tasks/new", &buf)
		mlReq.Header.Set("content-type", "application/json")
		if err != nil {
			log.Println("LOG_ERROR: HTTP POST ERROR",err.Error())
			fmt.Println(err)
		}
		mlReq.Header.Set("Authorization", "OANQ7uVfrSSZWzrE6g7qZSGPJdP/g/8j8fFV5GE5IgFU4IewqJ8EVa4mqpthgWJ0W0RkK6ZvuWRzOEw975FxAFXo3yTb0LPFS36YcoPR5APNa0z8E7eL2RG1GD6Xs4dOcESihgP1o6ObUH9CfG7T42coPP2XRRaCkqnGm183MqU=")
		_, _ = http.DefaultClient.Do(mlReq)*/
	}
	ctx.Status(200)
	return
}