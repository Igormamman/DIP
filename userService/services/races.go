package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"photoservice/userService/api/apiModel"
	dbModel "photoservice/userService/models"
	"time"
)

func (db *DatabaseService) GetRaceFromDB(UID string) (*dbModel.Race, []apiModel.ErrorJSON) {
	var race *dbModel.Race
	result := db.db.Where(&dbModel.Race{RaceID: UID}).First(&race)
	if result.Error != nil {
		fmt.Println("RACE NOT FOUND")
		return nil, []apiModel.ErrorJSON{{Classification: "DB", Message: "RACE NOT FOUND"}}
	}
	return race, nil
}

func (db *DatabaseService) NewRace(race dbModel.Race) *apiModel.ErrorJSON {
	result := db.db.Save(&race)
	if result.Error == nil {
		return nil
	}else{
		return &apiModel.ErrorJSON{Classification: "DB",Message: "ERROR SAVE NEW RACE"}
	}
}

func GetRaces() {
	fmt.Println("Getting Races")
	var races []apiModel.Race
	dateFrom := "2021-01-01"
	dateTo := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://marshalone.ru/api/races?date_from=%s&date_to=%s&user_uid=undefined", dateFrom, dateTo)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Errored when sending request to the server 2")
		return
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&races)
	for _, race := range races {
		result := DB.db.Where("race_id = ?", race.UID).
			Limit(1).
			Find(&dbModel.Race{})
		exist := result.RowsAffected > 0
		if exist {
			continue
		} else {
			date, err := time.Parse("2006-01-02", race.Date)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(race.Date)
			fmt.Println(date)
			dbRace := dbModel.Race{
				RaceID: race.UID,
				City:   race.City,
				Name:   race.Name,
				Date:   date,
			}
			DB.db.Create(&dbRace)
		}
	}
	//	fmt.Println("Errored when sending request to the server4")
	var racesEnd []dbModel.Race
	DB.db.Find(&racesEnd)
}
