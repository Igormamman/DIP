package models

type AppConfig struct {
	DatabaseConnectionString string `ini:"database_connection_string"`
	ApiPort                  string `ini:"api_port"`
	AccessKey                string `ini:"accessKeyID"`
	AccessSecret             string `ini:"secretAccessKey"`
	BucketName               string `ini:"bucketName"`
	WebPort                  string `ini:"web_port"`
	Host                     string `ini:"host"`
	CookieSecret             string `ini:"cookie_secret"`
}
