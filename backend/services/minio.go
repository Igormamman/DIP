package services

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var MinioClient *minio.Client

func InitMinio() {
	//endpoint := "s3.nl-ams.scw.cloud"
	endpoint := "hb.bizmrg.com"
	accessKeyID := Cfg.AccessKey
	secretAccessKey := Cfg.AccessSecret
	useSSL := true
	// Initialize minioService client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
	MinioClient = minioClient
	if err != nil {
		log.Fatalln("LOG_ERROR:",err.Error())
	}
	log.Printf("LOG_INFO: %#v\n", MinioClient) // minioClient is now setup\
}


