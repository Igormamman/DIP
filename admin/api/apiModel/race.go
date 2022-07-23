package apiModel

type InRace struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Date string `json:"date"`
	City string `json:"city"`
}

type OutRace struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Date int64  `json:"date"`
	City string `json:"city"`
}
