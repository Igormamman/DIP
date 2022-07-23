package api

import (
	service "admin/services"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/macaron.v1"
	"net/http"
)
import apiModel "admin/api/apiModel"


func getAdmin(ctx *macaron.Context){
	var pageData = apiModel.AdminHtmlPageObject{
		UrlParams: []apiModel.UrlParam{
			{ Table: "Photo", Label: "Фото"},
		},
	}
	ctx.Data["pageData"]=pageData
	ctx.HTML(200,"admin")
}

func UpdateAll() {
	offset := 0
	dbImageList, apiError := service.DB.GetPhotosByOffset(100, offset)
	if apiError != nil {
		fmt.Println(apiError)
	}
	for len(dbImageList) > 0 {
		fmt.Println(offset)
		for _, dbImage := range dbImageList {
			fileInfo := apiModel.NewFile{Token: dbImage.PUID}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(fileInfo)
			if err != nil {
				fmt.Println(err)
			}
			mlReq, err := http.NewRequest("POST", "http://backend:4000/ml/task/create", &buf)
			mlReq.Header.Set("content-type", "application/json")
			if err != nil {
				fmt.Println(err)
			}
			_, err = http.DefaultClient.Do(mlReq)
		}
		offset += len(dbImageList)
		dbImageList, apiError = service.DB.GetPhotosByOffset(100, offset);
	}
}