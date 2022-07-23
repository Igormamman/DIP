package models

type AppConfig struct {
	MinioPort                string `ini:"minio_port"`
	AccessKey                string `ini:"accessKeyID"`
	AccessSecret             string `ini:"secretAccessKey"`
	BucketName               string `ini:"bucketName"`
	Host                     string `ini:"host"`
}
