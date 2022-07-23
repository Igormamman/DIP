package services

import (
	"fmt"
	"gorm.io/gorm"
	"photoservice/userService/api/apiModel"
	dbModel "photoservice/userService/models"
	"time"
)

func (db *DatabaseService) CreateUserCookie(userCookie dbModel.UserCookie) []apiModel.ErrorJSON {
	result := db.db.Save(&userCookie)
	if result.Error != nil{
		return []apiModel.ErrorJSON{{Classification: "DB",Message: result.Error.Error()}}
	}
	return nil
}

func (db *DatabaseService) CreateUserAccess(userAccess dbModel.UserAccess) []apiModel.ErrorJSON {
	result := db.db.Save(&userAccess)
	if result.Error != nil{
		return []apiModel.ErrorJSON{{Classification: "DB",Message: result.Error.Error()}}
	}
	return nil
}


func (db *DatabaseService) GetUserCookie(cookie string) (*dbModel.UserCookie,[]apiModel.ErrorJSON) {
	userCookie := dbModel.UserCookie{}
	result := db.db.Where(&dbModel.UserCookie{Cookie: cookie}).First(&userCookie)
	if result.Error != nil{
		fmt.Println("ERROR DB",result.Error.Error())
		return nil,[]apiModel.ErrorJSON{{Classification: "DB",Message: result.Error.Error()}}
	}else{
		if userCookie.CreatedAt.Add(time.Minute*5).Before(time.Now()){
			db.db.Unscoped().Delete(&userCookie)
			return nil,[]apiModel.ErrorJSON{{Classification: "DB",Message: "DB record expired"}}
		}
	}
	return &userCookie,nil
}

func (db *DatabaseService) GetUserAccess(cookie string, raceID string) (*dbModel.UserAccess,[]apiModel.ErrorJSON) {
	UserAccess := dbModel.UserAccess{}
	result := db.db.Where(&dbModel.UserAccess{Cookie: cookie, RaceUID: raceID}).First(&UserAccess)
	if result.Error != nil{
		fmt.Println("ERROR DB",result.Error.Error())
		return nil,[]apiModel.ErrorJSON{{Classification: "DB",Message: result.Error.Error()}}
	}else{
		if UserAccess.CreatedAt.Add(time.Minute*5).Before(time.Now()){
			db.db.Unscoped().Delete(&UserAccess)
			return nil,[]apiModel.ErrorJSON{{Classification: "DB",Message: "DB record expired"}}
		}
	}
	return &UserAccess,nil
}

func (db *DatabaseService) DeleteStaleUserCookie() (){
	fmt.Println("RUNNING USER COOKIE HOOK")
	var deleteData []dbModel.UserCookie
	var result *gorm.DB
	result = db.db.Where("created_at < ?", time.Now().Add(-time.Minute*5)).Find(&deleteData)
	if result.Error != nil{
		fmt.Println("DELETE HOOK ERROR")
		fmt.Println(result.Error.Error())
	}else{
		fmt.Println("ROWS TO DELETE:",result.RowsAffected)
		if result.RowsAffected != 0{
			db.db.Unscoped().Delete(&deleteData)
		}
	}
	return
}

func (db *DatabaseService) DeleteUserAccessByCookie(cookie string) (){
	var deleteData []dbModel.UserAccess
	var result *gorm.DB
	result = db.db.Where(&dbModel.UserAccess{Cookie: cookie}).Find(&deleteData)
	if result.Error != nil{
		fmt.Println("DELETE User ACCESS ON LOGOUT ERROR")
		fmt.Println(result.Error.Error())
	}else{
		fmt.Println("ROWS TO DELETE:",result.RowsAffected)
		if result.RowsAffected != 0{
			db.db.Unscoped().Delete(&deleteData)
		}
	}
	return
}

func (db *DatabaseService) DeleteUserCookie(cookie string) (){
	var deleteData []dbModel.UserCookie
	var result *gorm.DB
	result = db.db.Where(&dbModel.UserCookie{Cookie: cookie}).Find(&deleteData)
	if result.Error != nil{
		fmt.Println("DELETE User COOKIE ON LOGOUT ERROR")
		fmt.Println(result.Error.Error())
	}else{
		fmt.Println("ROWS TO DELETE:",result.RowsAffected)
		if result.RowsAffected != 0{
			db.db.Unscoped().Delete(&deleteData)
		}
	}
	return
}

func (db *DatabaseService) DeleteStaleUserAccess() (){
	fmt.Println("RUNNING USER ACCESS HOOK")
	var deleteData []dbModel.UserAccess
	result := db.db.Where("created_at < ?", time.Now().Add(-time.Minute*5)).Find(&deleteData)
	if result.Error != nil{
		fmt.Println("DELETE HOOK ERROR")
		fmt.Println(result.Error.Error())
	}else{
		fmt.Println("ROWS TO DELETE:",result.RowsAffected)
		if result.RowsAffected != 0{
			db.db.Unscoped().Delete(&deleteData)
		}
	}
	return
}