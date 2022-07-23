package services

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"time"
	_ "time/tzdata"
)

var Scheduler gocron.Scheduler

func InitScheduler() {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
		fmt.Println("FAILED TO PARSE LOCATION")
	}
	Scheduler = *gocron.NewScheduler(location)

	UserCookieJobInit()
	UserAccessJobInit()
	RaceJobInit()
	Scheduler.StartAsync()
}

func UserCookieJobInit() {
	_, err := Scheduler.Every(1).Minute().StartAt(time.Now().Add(time.Second * 1)).Do(DB.DeleteStaleUserCookie)
	//_, err := Scheduler.Every(2).Second().StartAt(time.Now()).Do(getRaces)
	if err != nil {
		fmt.Println(err)
	}
}

func UserAccessJobInit() {
	_, err := Scheduler.Every(1).Minute().StartAt(time.Now().Add(time.Second * 1)).Do(DB.DeleteStaleUserAccess)
	//_, err := Scheduler.Every(2).Second().StartAt(time.Now()).Do(getRaces)
	if err != nil {
		fmt.Println(err)
	}
}

func RaceJobInit(){
	_, err := Scheduler.Every(1).Day().At("12:00").StartAt(time.Now().Add(time.Second*5)).Do(GetRaces)
	//_, err := Scheduler.Every(2).Second().StartAt(time.Now()).Do(getRaces)
	if err != nil {
		fmt.Println(err)
	}
}
