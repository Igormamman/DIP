package models

type AppConfig struct {
	DatabaseConnectionString string `ini:"database_connection_string_userService"`
	UserPort                 string `ini:"user_port"`
	AccessKey                string `ini:"accessKeyID"`
	Host                     string `ini:"host"`
	CookieSecret             string `ini:"cookie_secret"`
}
