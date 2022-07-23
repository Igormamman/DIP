package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gopkg.in/macaron.v1"
	"image"
	_"image/jpeg"
	_"image/png"
	"io"
	"log"
	"minio/api/apiModel"
	"minio/services"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var fileList = make(map[string]apiModel.NewFile)

func NewFile(ctx *macaron.Context, newFiles []apiModel.NewFile) {
	lock.Lock()
	defer lock.Unlock()
	for _, file := range newFiles {
		fileList[file.FilePath] = file
	}
	ctx.JSON(http.StatusOK, "")
	return
}

func LoadFile(ctx *macaron.Context) {
	formData := apiModel.LoadStruct{}
	readForm, _ := ctx.Req.MultipartReader()
	var formBuf bytes.Buffer
	part, errPart := readForm.NextPart()
	if part.FormName() != "filepath" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		sendError(fileList[formData.FilePath])
		log.Println("LOG_ERROR: FORM ERROR",errPart.Error())
		ctx.JSON(http.StatusBadRequest, "FORM ERROR")
		return
	} else {
		_, _ = io.Copy(&formBuf, part)
		formData.FilePath = string(formBuf.Bytes())
		formBuf.Reset()
	}
	part, errPart = readForm.NextPart()
	if part.FormName() != "file_size" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		log.Println("LOG_ERROR: FORM ERROR",errPart.Error())
		sendError(fileList[formData.FilePath])
		ctx.JSON(http.StatusBadRequest, "FORM ERROR")
		return
	} else {
		_, _ = io.Copy(&formBuf, part)
		formData.FileSize, _ = strconv.Atoi(string(formBuf.Bytes()))
		formBuf.Reset()
	}
	lock.RLock()
	if fileList[formData.FilePath].FileSize != formData.FileSize {
		fmt.Println("ERROR", fileList[formData.FilePath], formData.FileSize)
		log.Println("LOG_ERROR: WRONG FILESIZE", errPart.Error())
		sendError(fileList[formData.FilePath])
		ctx.JSON(http.StatusBadRequest, "WRONG FILE")
		return
	}
	lock.RUnlock()
	part, errPart = readForm.NextPart()
	if part.FormName() != "file_type" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		log.Println("LOG_ERROR: FORM ERROR",errPart.Error())
		sendError(fileList[formData.FilePath])
		ctx.JSON(http.StatusBadRequest, "FORM ERRROR")
		return
	} else {
		_, _ = io.Copy(&formBuf, part)
		formData.FileType = string(formBuf.Bytes())
		formBuf.Reset()
	}
	part, errPart = readForm.NextPart()
	if part.FormName() != "file" || errPart != nil {
		fmt.Println("FORM ERROR", part.FormName())
		log.Println("LOG_ERROR: FORM ERROR",errPart.Error())
		sendError(fileList[formData.FilePath])
		ctx.JSON(http.StatusBadRequest, "FORM ERRROR")
		return
	} else {
		formData.File = part
	}
	var buf bytes.Buffer
	tee := io.TeeReader(formData.File, &buf)
	_,format,error := image.DecodeConfig(tee)
	if error != nil || (format != "png" && format != "jpeg" && format != "jpg"){
		log.Println("LOG_ERROR: IMAGE PARSE ERROR",errPart.Error())
		ctx.JSON(http.StatusBadRequest,"image parse error")
		return
	}else{
		fmt.Println("FILETYPE"+format)
		log.Println("LOG_INFO: GOT IMAGE: filetype =",format)
	}
	_, err := services.MinioClient.PutObject(context.Background(), services.Cfg.BucketName,
		formData.FilePath, io.MultiReader(&buf,formData.File), int64(formData.FileSize), minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("MINIO ERROR", err.Error())
		log.Println("LOG_ERROR: MINIO LOAD ERROR",err.Error())
		sendError(fileList[formData.FilePath])
		ctx.JSON(http.StatusBadRequest,"minioService put error")
	} else {
		if formData.FileType == "ml" {
			sendTask(fileList[formData.FilePath])
		}
		lock.Lock()
		delete(fileList, formData.FilePath)
		lock.Unlock()
		ctx.JSON(http.StatusOK, "ok")
	}
	return
}

func sendTask(fileInfo apiModel.NewFile) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(fileInfo)
	if err != nil {
		fmt.Println(err)
	}
	mlReq, err := http.NewRequest("POST", "http://backend:4000/task/create", &buf)
	mlReq.Header.Set("content-type", "application/json")
	if err != nil {
		log.Println("LOG_ERROR: HTTP GET ERROR",err.Error())
		fmt.Println(err)
	}
	_, err = http.DefaultClient.Do(mlReq)
}

func sendError(fileInfo apiModel.NewFile) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(fileInfo)
	if err != nil {
		fmt.Println(err)
	}
	mlReq, err := http.NewRequest("POST", "http://backend:4000/photo/delete", &buf)
	mlReq.Header.Set("content-type", "application/json")
	if err != nil {
		log.Println("LOG_ERROR: HTTP POST ERROR",err.Error())
		fmt.Println(err)
	}
	_, err = http.DefaultClient.Do(mlReq)
}

func GetFile(ctx *macaron.Context) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		log.Println("LOG_ERROR: QUERY PARSE ERROR",err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if ctx.Req.Header.Get("If-Modified-Since") != "" {
		ctx.Resp.WriteHeader(http.StatusNotModified)
		return
	}
	var filepath string
	if value, exist := params["type"]; exist {
		filepath = fmt.Sprintf("%s/%s/%s/%s/%s%s",
			ctx.Params("year"), ctx.Params("month"), ctx.Params("day"),
			ctx.Params("name"), value[0], ctx.Params("name"))
	} else {
		filepath = fmt.Sprintf("%s/%s/%s/%s/%s",
			ctx.Params("year"), ctx.Params("month"), ctx.Params("day"),
			ctx.Params("name"), ctx.Params("name"))
	}
	fmt.Println(filepath)
	if filepath != "" {
		file, error := services.MinioClient.GetObject(context.Background(), services.Cfg.BucketName, filepath, minio.GetObjectOptions{})
		stat,statError := file.Stat()
		if error != nil || file == nil {
			sendError(apiModel.NewFile{FileName: ctx.Params("name")})
			fmt.Println(error.Error())
			log.Println("LOG_ERROR: MINIO GET ERROR",error.Error())
			ctx.JSON(400, "file doesn't exist")
			return
		}
		if statError  != nil{
			sendError(apiModel.NewFile{FileName: ctx.Params("name")})
			fmt.Println(stat.Size)
			fmt.Println(statError.Error())
			log.Println("LOG_ERROR: MINIO GET ERROR",statError.Error())
			ctx.JSON(400, "file doesn't exist")
			return
		} else {
			defer file.Close()
			http.ServeContent(ctx.Resp, ctx.Req.Request, filepath, time.Unix(-530031600, 0), file)
			return
		}
	}
}

func DeleteFile(ctx *macaron.Context) {
	var filepath string
	filepath = fmt.Sprintf("%s/%s/%s/%s",
		ctx.Params("year"), ctx.Params("month"), ctx.Params("day"),
		ctx.Params("name"))
	fmt.Println(filepath)
	if filepath != "" {
		deleteObjs := services.MinioClient.ListObjects(context.Background(), services.Cfg.BucketName, minio.ListObjectsOptions{Recursive: true, Prefix: filepath})
		removeError := services.MinioClient.RemoveObjects(context.Background(), services.Cfg.BucketName, deleteObjs, minio.RemoveObjectsOptions{})
		err := <-removeError
		if err.Err != nil {
			log.Println("LOG_ERROR: MINIO DELETE ERROR",err.Err.Error())
			fmt.Println(err.Err)
			ctx.JSON(400, "error")
		}
		ctx.JSON(200, "ok")
		return
	}
}
