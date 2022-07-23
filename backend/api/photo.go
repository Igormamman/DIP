package api

import (
	"bytes"
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gopkg.in/macaron.v1"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"photoservice/backend/api/apiModel"
	dbModel "photoservice/backend/models"
	"photoservice/backend/services"
	"strconv"
	"time"
)

const (
	baseImage = 1
	preview   = 2
	watermark = 3
	ML        = 4
)

func getPhotoCountRouter(ctx *macaron.Context) {
	params, _ := url.ParseQuery(ctx.Req.URL.RawQuery)
	competitor := ""
	raceUID := ""
	detected := true
	if value, exist := params["competitor"]; exist && len(value) > 0 {
		if value[0] != "" && value[0] != "undefined" {
			competitor = value[0]
		}
	}
	if value, exist := params["raceUID"]; exist && len(value) > 0 {
		if value[0] != "" && value[0] != "undefined" {
			raceUID = value[0]
		}
	}
	if value, exist := params["detected"]; exist && len(value) > 0 {
		if value[0] != "false" {
			detected = true
		} else {
			detected = false
		}
	}
	ctx.JSON(http.StatusOK, services.DB.GetPhotoCount(competitor, raceUID, detected))
}

func newTaskRouter(ctx *macaron.Context, taskData apiModel.TaskData) {
	userAccess := getUserAccess(ctx, taskData.RaceID)
	if userAccess != nil && !(userAccess.IsOrg || userAccess.IsMediaConfirmed) {
		ctx.JSON(400, apiModel.ErrorJSON{Classification: "Access", Message: "No access to operation"})
		return
	}
	userID := userAccess.UserID
	if userID != "" {
		fmt.Println(userID)
		for _, file := range taskData.FileData {
			file.UserID = userID
			file.SID = ctx.GetCookie("connect.sid")
			file.RaceID = taskData.RaceID
			apiError := services.DB.NewTaskFile(file)
			if apiError != nil {
				ctx.JSON(400, apiError)
			}
		}
		ctx.JSON(200, "")
	} else {
		ctx.JSON(400, "USER DON'T EXIST")
	}
}

func getUserAccess(ctx *macaron.Context, raceID string) *apiModel.InRaceAccess {
	cookie, err := ctx.Req.Cookie("connect.sid")
	if err != nil || cookie.Value == "" {
		return nil
	}
	requestURL := fmt.Sprintf("http://userService:4005/getRaceAccess?ruid=%s", raceID)
	req, err := http.NewRequest("GET", requestURL, nil)
	fmt.Println("connect.sid:", cookie.Value)
	req.AddCookie(cookie)
	if err != nil {
		log.Println("LOG_ERROR: HTTP GET ERROR", err.Error())
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	j := apiModel.InRaceAccess{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	fmt.Println("J:", j, j.RaceID, &j.IsOrg)
	if err == nil || err == io.EOF {
		return &j
	} else {
		return nil
	}
}

func getUserID(ctx *macaron.Context) *string {
	cookie, err := ctx.Req.Cookie("connect.sid")
	if err != nil || cookie.Value == "" {
		return nil
	}
	req, err := http.NewRequest("GET", "http://userService:4005/getUID", nil)
	fmt.Println("connect.sid:", cookie.Value)
	req.AddCookie(cookie)
	if err != nil {
		log.Println("LOG_ERROR: HTTP GET ERROR", err.Error())
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	j := struct {
		UserID string `json:"user_id"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	fmt.Println("J:", j, j.UserID, &j.UserID)
	if (err == nil || err == io.EOF) && j.UserID != "" {
		return &j.UserID
	} else {
		return nil
	}
}

func newUploadMetaRouter(ctx *macaron.Context, uploadMeta apiModel.UploadMeta) {
	var userID *string
	userID = getUserID(ctx)
	if userID != nil {
		uploadMeta.UserID = *userID
		uploadInfo, apiError := services.DB.UpdateTaskInfo(uploadMeta)
		if apiError != nil || uploadInfo == nil {
			ctx.JSON(400, apiError)
		} else {
			dbImage, apiError := services.DB.CreateNewPhoto(*uploadInfo)
			if apiError != nil {
				if apiError[0].Message != "This Photo exist in DB" || dbImage == nil {
					ctx.JSON(400, apiError)
					return
				} else {
					var responseData []apiModel.PathToUpload
					responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, ""})
					responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "preview"})
					responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "watermark"})
					responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "ml"})
					ctx.JSON(http.StatusOK, responseData)
					return
				}
			} else {
				uploadInfo.Token = dbImage.PUID
				services.DB.SaveUploadInfo(*uploadInfo)
				var responseData []apiModel.PathToUpload
				responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, ""})
				responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "preview"})
				responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "watermark"})
				responseData = append(responseData, apiModel.PathToUpload{uploadInfo.Token, "ml"})
				ctx.JSON(http.StatusOK, responseData)
				return
			}
		}
	} else {
		ctx.JSON(http.StatusBadRequest, "USER DON'T EXIST")
		return
	}
}

func newPhotoRouter(ctx *macaron.Context) {
	readForm, _ := ctx.Req.MultipartReader()
	var file *multipart.Part
	var formBuf bytes.Buffer
	var filetype string
	var token string
	var task *dbModel.UploadInfo
	part, errPart := readForm.NextPart()
	if part.FormName() != "token" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		fmt.Println("FORM ERROR", errPart)
		//	log.Println("LOG_ERROR: FORM ERROR", errPart.Error())
		ctx.JSON(http.StatusBadRequest, "FORM ERROR")
		return
	} else {
		_, _ = io.Copy(&formBuf, part)
		token = string(formBuf.Bytes())
		formBuf.Reset()
	}
	part, errPart = readForm.NextPart()
	fmt.Println(part)
	if part.FormName() != "file_type" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		fmt.Println("FORM ERROR", errPart)
		//	log.Println("LOG_ERROR: FORM ERROR", errPart.Error())
		ctx.JSON(http.StatusBadRequest, "FORM ERROR")
		return
	} else {
		_, _ = io.Copy(&formBuf, part)
		filetype = string(formBuf.Bytes())
		formBuf.Reset()
	}
	if filetype != "ml" && filetype != "" && filetype != "preview" && filetype != "watermark" {
		ctx.JSON(http.StatusBadRequest, "wrong filetype ()")
		return
	}
	part, errPart = readForm.NextPart()
	if part.FormName() != "file" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		log.Println("LOG_ERROR: FORM ERROR", errPart.Error())
		ctx.JSON(http.StatusBadRequest, "FORM ERROR1")
		return
	} else {
		file = part
		/*tokenString := fmt.Sprintf("%s%s%s", file.FileName(), "", "-1")
		token := tools.StringGenerateSHA256(tokenString)
		fmt.Println(file.FileName(), "-1", "", token)*/
		task, _ = services.DB.GetUploadInfoByToken(token)
		if task == nil {
			fmt.Println("TASK NIL")
			ctx.JSON(400, "task doesn't exist")
			return
		}
	}
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	_, format, error := image.DecodeConfig(tee)
	if error != nil || (format != "png" && format != "jpeg" && format != "jpg") {
		log.Println("LOG_ERROR: IMAGE PARSE ERROR", errPart.Error())
		ctx.JSON(http.StatusBadRequest, "image parse error")
		return
	} else {
		fmt.Println("FILETYPE" + format)
		log.Println("LOG_INFO: GOT IMAGE: extension =", format)
	}
	// TO DO TO DO TO DO TO DO TO DO
	dbImage, _ := services.DB.GetPhotoByPUID(task.Token)
	year, month, day := dbImage.LoadedAt.Date()
	fileNamePrefix := fmt.Sprintf("%d/%d/%d", year, month, day)
	switch filetype {
	case "":
		uploadPath := fmt.Sprintf("%s/%s/%s", fileNamePrefix, dbImage.PhotoName, dbImage.PhotoName)
		_, err := services.MinioClient.PutObject(context.Background(), services.Cfg.BucketName,
			uploadPath, io.MultiReader(&buf, file), int64(task.ImageSize), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("MINIO ERROR", err.Error())
			log.Println("LOG_ERROR: MINIO LOAD ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, "minioService put error")
		}
		task.MainUploaded = true
		services.DB.SaveUploadInfo(*task)
		break
	case "ml":
		uploadPath := fmt.Sprintf("%s/%s/ml%s", fileNamePrefix, dbImage.PhotoName, dbImage.PhotoName)
		_, err := services.MinioClient.PutObject(context.Background(), services.Cfg.BucketName,
			uploadPath, io.MultiReader(&buf, file), int64(task.MLSize), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("MINIO ERROR", err.Error())
			log.Println("LOG_ERROR: MINIO LOAD ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, "minioService put error")
		}
		task.MLUploaded = true
		services.DB.SaveUploadInfo(*task)
		break
	case "watermark":
		uploadPath := fmt.Sprintf("%s/%s/watermark%s", fileNamePrefix, dbImage.PhotoName, dbImage.PhotoName)
		/*	options := minio.PutObjectOptions{
				UserMetadata: make(map[string]string),
			}
			options.UserMetadata["X-amz-acl"] = "public-read"*/
		_, err := services.MinioClient.PutObject(context.Background(), services.Cfg.BucketName,
			uploadPath, io.MultiReader(&buf, file), int64(task.WatermarkSize), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("MINIO ERROR", err.Error())
			log.Println("LOG_ERROR: MINIO LOAD ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, "minioService put error")
		}
		task.WatermarkUploaded = true
		services.DB.SaveUploadInfo(*task)
		break
	case "preview":
		uploadPath := fmt.Sprintf("%s/%s/resized%s", fileNamePrefix, dbImage.PhotoName, dbImage.PhotoName)
		/*	options := minio.PutObjectOptions{
				UserMetadata: make(map[string]string),
			}
			options.UserMetadata["X-amz-acl"] = "public-read"*/
		_, err := services.MinioClient.PutObject(context.Background(), services.Cfg.BucketName,
			uploadPath, io.MultiReader(&buf, file), int64(task.ResizedSize), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("MINIO ERROR", err.Error())
			log.Println("LOG_ERROR: MINIO LOAD ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, "minioService put error")
		}
		task.ResizedUploaded = true
		services.DB.SaveUploadInfo(*task)
		break
	}
	ctx.JSON(http.StatusOK, "ok")

	return
}

func GetImageType(params url.Values) int8 {
	var imageType int8
	if value, exist := params["type"]; exist {
		switch value[0] {
		case "preview":
			imageType = preview
		case "ml":
			imageType = ML
		case "watermark":
			imageType = watermark
		}
	} else {
		imageType = baseImage
	}
	return imageType
}

func getPhotoMeta(ctx *macaron.Context) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if UID, exist := params["UUID"]; exist {
		photo, apiError := services.DB.GetPhotoByPUID(UID[0])
		if apiError != nil || photo == nil {
			ctx.JSON(http.StatusBadRequest, apiError)
			return
		} else {
			tags := services.DB.GetPhotoTags(*photo)
			photoMeta := apiModel.PhotoMeta{UUID: photo.PUID, RUID: photo.RaceID, Competitors: tags, Height: photo.Height, Width: photo.Width}
			ctx.JSON(http.StatusOK, photoMeta)
			return
		}
	} else {
		log.Println("LOG_WARNING: PHOTO DOESN'T EXIST", ctx.Params("name"))
		ctx.JSON(400, "Photo doesn't exist")
		return
	}
}

func getPhotoPreviewUrls(ctx *macaron.Context) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if value, exist := params["photographer"]; exist && value[0] == "true" {
		getPhotographerPreviews(ctx, params)
		return
	} else {
		getPreviewsByParam(ctx, params)
		return
	}
}

func getPreviewsByParam(ctx *macaron.Context, params url.Values) {
	competitor := ""
	raceUID := ""
	detected := true
	offsetValue, offsetExist := params["offset"]
	limitValue, limitExist := params["limit"]
	if value, exist := params["competitor"]; exist && len(value) > 0 {
		if value[0] != "" && value[0] != "undefined" {
			competitor = value[0]
		}
	}
	if value, exist := params["raceUID"]; exist && len(value) > 0 {
		if value[0] != "" && value[0] != "undefined" {
			raceUID = value[0]
		}
	}
	if value, exist := params["detected"]; exist && len(value) > 0 {
		if value[0] != "false" {
			detected = true
		} else {
			detected = false
		}
	}
	var count int64
	var photos []dbModel.Photo
	if limitExist && offsetExist {
		limit, err := strconv.Atoi(limitValue[0])
		if err != nil {
			log.Println("LOG_ERROR: PARAM CONVERT ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		offset, err := strconv.Atoi(offsetValue[0])
		if err != nil {
			log.Println("LOG_ERROR: PARAM CONVERT ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		var apiError []apiModel.ErrorJSON
		if detected == true {
			fmt.Printf("Get Detected")
			photos, apiError = services.DB.GetDetectedPhotoFromDB(competitor, raceUID, limit, offset)
			fmt.Println("batch count", len(photos))
			count = services.DB.GetPhotoCount(competitor, raceUID, true)
			fmt.Println(count)
		} else {
			fmt.Printf("Get Undetected")
			photos, apiError = services.DB.GetUndetectedPhotoFromDB(raceUID, limit, offset)
			fmt.Println("batch count", len(photos))
			count = services.DB.GetPhotoCount("", raceUID, false)
			fmt.Println("count:", count)
		}
		if apiError != nil && len(photos) > 0 {
			fmt.Println(apiError)
			ctx.JSON(http.StatusBadRequest, apiError)
			return
		}
	}
	var requestUrls []string
	for _, photo := range photos {
		fmt.Println(photo.PUID)
		requestUrls = append(requestUrls, photo.PUID)
	}
	//	"http://localhost:4004/photo/%s/%s?type=resized"
	//"https://minio-dip.duckdns.org/photo/%s/%s?type=resized"
	if len(requestUrls) > 0 {
		ctx.JSON(http.StatusOK, apiModel.OutPhotoPreviewData{Count: count, PreviewURL: requestUrls})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
}

func getPhotographerPreviews(ctx *macaron.Context, params url.Values) {
	fmt.Println("GET PHOTOGRAPHER PREVIEWS")
	raceUID := ""
	offsetValue, offsetExist := params["offset"]
	limitValue, limitExist := params["limit"]
	if value, exist := params["raceUID"]; exist && len(value) > 0 {
		if value[0] != "" && value[0] != "undefined" {
			raceUID = value[0]
		}
	}
	var count int64
	var photos []dbModel.Photo
	if limitExist && offsetExist {
		limit, err := strconv.Atoi(limitValue[0])
		if err != nil {
			log.Println("LOG_ERROR: PARAM CONVERT ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		offset, err := strconv.Atoi(offsetValue[0])
		if err != nil {
			log.Println("LOG_ERROR: PARAM CONVERT ERROR", err.Error())
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		var apiError []apiModel.ErrorJSON
		userID := getUserID(ctx)
		if userID == nil {
			ctx.JSON(http.StatusBadRequest, "user don't exist")
			return
		}
		photos, apiError = services.DB.GetUserPhotosByOffset(limit, offset, *userID, raceUID)
		if apiError != nil && len(photos) > 0 {
			fmt.Println(apiError)
			ctx.JSON(http.StatusBadRequest, apiError)
			return
		}
		count = services.DB.GetUserPhotoCount(*userID, raceUID)
		var requestUrls []string
		for _, photo := range photos {
			requestUrls = append(requestUrls, photo.PUID)
		}
		//	"http://localhost:4004/photo/%s/%s?type=resized"
		//"https://minio-dip.duckdns.org/photo/%s/%s?type=resized"
		if len(requestUrls) > 0 {
			ctx.JSON(http.StatusOK, apiModel.OutPhotoPreviewData{Count: count, PreviewURL: requestUrls})
			return
		} else {
			ctx.JSON(http.StatusOK, nil)
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, "wrong params")
		return
	}
}

func deletePhotoRouter(ctx *macaron.Context) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if UID, exist := params["UUID"]; exist {
		photo, apiError := services.DB.GetPhotoByPUID(UID[0])
		if apiError != nil || photo == nil {
			ctx.JSON(http.StatusBadRequest, apiError)
			return
		} else {
			userID := getUserID(ctx)
			if userID != nil {
				if photo.UserID == *userID {
					apiError := services.DB.DeletePhoto(*photo, true)
					if apiError != nil {
						ctx.JSON(400, apiError)
						return
					} else {
						ctx.JSON(200, "OK")
						return
					}
				} else {
					fmt.Println("PHOTO USER_ID:", photo.UserID)
					fmt.Println("CURRENT USER_ID:", *userID)
					ctx.JSON(400, "photo belongs to other user")
					return
				}
			} else {
				ctx.JSON(400, "user not authorized")
				return
			}
		}
	}
}

func GetFile(ctx *macaron.Context) {
	var UUID string
	var filepath string
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if ctx.Req.Header.Get("If-Modified-Since") != "" {
		ctx.Resp.WriteHeader(http.StatusNotModified)
		return
	}
	if value, exist := params["UUID"]; exist {
		UUID = value[0]
	}
	fileType, exist := params["type"]
	if !exist {
		ctx.JSON(400, "хе хе хе")
		return
	}
	switch fileType[0] {
	case "resized", "watermark":
		dbImage, _ := services.DB.GetPhotoByPUID(UUID)
		year, month, day := dbImage.LoadedAt.Date()
		filepath = fmt.Sprintf("%d/%d/%d/%s/%s%s", year, month, day, dbImage.PhotoName, fileType[0], dbImage.PhotoName)
	case "ml":
		if token, tokenExist := params["token"]; tokenExist && token[0] == "iojqwejknasdfjknxcvasdfkqwermn" {
			dbImage, _ := services.DB.GetPhotoByPUID(UUID)
			year, month, day := dbImage.LoadedAt.Date()
			filepath = fmt.Sprintf("%d/%d/%d/%s/%s%s", year, month, day, dbImage.PhotoName, fileType[0], dbImage.PhotoName)
		} else {
			ctx.JSON(400, "хе хе хе")
			return
		}
	case "":
		dbImage, _ := services.DB.GetPhotoByPUID(UUID)
		if dbImage == nil {
			ctx.JSON(400, "хе хе хе")
			return
		}
		year, month, day := dbImage.LoadedAt.Date()
		filepath = fmt.Sprintf("%d/%d/%d/%s/%s%s", year, month, day, dbImage.PhotoName, "ml", dbImage.PhotoName)
		if filepath != "" {
			returnURL, err := services.MinioClient.PresignedGetObject(context.Background(), services.Cfg.BucketName, filepath, time.Duration(3600)*time.Second, nil)
			if err == nil {
				ctx.JSON(200, struct{ URL string }{returnURL.String()})
				return
			} else {
				fmt.Println(err)
				ctx.JSON(400, "хе хе хе")
				return
			}
		}
	default:
		ctx.JSON(400, "хе хе хе")
		return
	}
	fmt.Println(filepath)
	if filepath != "" {
		file, error := services.MinioClient.GetObject(context.Background(), services.Cfg.BucketName, filepath, minio.GetObjectOptions{})
		stat, statError := file.Stat()
		if error != nil || file == nil {
			fmt.Println(error.Error())
			log.Println("LOG_ERROR: MINIO GET ERROR", error.Error())
			ctx.JSON(400, "file doesn't exist")
			return
		}
		if statError != nil {
			fmt.Println(stat.Size)
			fmt.Println(statError.Error())
			log.Println("LOG_ERROR: MINIO GET ERROR", statError.Error())
			ctx.JSON(400, "file doesn't exist")
			return
		} else {
			defer file.Close()
			http.ServeContent(ctx.Resp, ctx.Req.Request, filepath, time.Unix(-530031600, 0), file)
			return
		}
	}
}
