package apiModel

import "time"

type Race struct {
	UID              string `json:"uid"`
	Name             string `json:"name"`
	Date             string `json:"date"`
	City             string `json:"city"`
	Preview          string `json:"titlePicture"`
	IsMedia          bool   `json:"isMedia"`
	IsMediaConfirmed bool   `json:"isMediaConfirmed"`
	IsCompetitor     bool   `json:"isCompetitor"`
	IsOrg            bool   `json:"isOrg"`
}

type OutRaceInfo struct {
	RaceID string    `json:"uid"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	City   string    `json:"city"`
}

type OutRaceAccess struct {
	RaceID           string `json:"uid"`
	UserID           string `json:"cUID"`
	IsMedia          bool   `json:"isMedia"`
	IsMediaConfirmed bool   `json:"isMediaConfirmed"`
	IsCompetitor     bool   `json:"isCompetitor"`
	IsOrg            bool   `json:"isOrg"`
}

type UserInfo struct {
	UserID   string `json:"cUID"`
	PhotoID  string `json:"cPhoto"`
	UserName string `json:"cName"`
}

type GetRace struct {
	Response *Race   `json:"response"`
	Error    *string `json:"errorMessage"`
}

type GetRaces struct {
	Response []Race  `json:"response"`
	Error    *string `json:"errorMessage"`
}

type UserRequest struct {
	Response *UserInfo `json:"response"`
	Error    *string   `json:"errorMessage"`
}
