package apiModel

type InRaceAccess struct {
	RaceID           string `json:"uid"`
	UserID           string `json:"cUID"`
	IsMedia          bool   `json:"isMedia"`
	IsMediaConfirmed bool   `json:"isMediaConfirmed"`
	IsCompetitor     bool   `json:"isCompetitor"`
	IsOrg            bool   `json:"isOrg"`
}
