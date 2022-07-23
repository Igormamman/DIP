package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"photoservice/userService/api/apiModel"
	dbModel "photoservice/userService/models"
	"photoservice/userService/services"
	"time"

	"gopkg.in/macaron.v1"
)

func getRacesRouter(ctx *macaron.Context) {
	races, err := getUserRaces(ctx)
	if err != nil {
		ctx.JSON(400, err.Error())
	} else {
		ctx.JSON(200, races)
	}
}

func getUserIDRouter(ctx *macaron.Context) {
	fmt.Println("cookie:", ctx.GetCookie("connect.sid"))
	cookie := ctx.GetCookie("connect.sid")
	userID, error := getUserID(cookie)
	if error != nil || userID == "" {
		fmt.Println(error.Error())
		ctx.JSON(400, "USER NOT FOUND")
		return
	} else {
		ctx.JSON(200, struct {
			UserID string `json:"user_id"`
		}{UserID: userID})
		return
	}
}

func LogoutRouter(ctx *macaron.Context) {
	req, err := http.NewRequest("GET", "https://marshalone.ru/api/private/logout", nil)
	if err != nil {
		ctx.JSON(400, nil)
		return
	}
	cookie := http.Cookie{Name: "connect.sid", Value: ctx.GetCookie("connect.sid")}
	req.AddCookie(&cookie)
	services.DB.DeleteUserAccessByCookie(ctx.GetCookie("connect.sid"))
	services.DB.DeleteUserCookie(ctx.GetCookie("connect.sid"))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(400, nil)
	} else {
		ctx.JSON(200, nil)
	}
	return
}

func getUserInfoRouter(ctx *macaron.Context) {
	req, err := http.NewRequest("GET", "https://marshalone.ru/api/private/current_user", nil)
	cookie := http.Cookie{Name: "connect.sid", Value: ctx.GetCookie("connect.sid")}
	req.AddCookie(&cookie)
	if err != nil {
		log.Println("LOG_ERROR: HTTP GET ERROR", err.Error())
		fmt.Println(err)
		ctx.JSON(400, "1")
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(400, "m1 req err")
		return
	}
	if resp.StatusCode != 200 {
		ctx.JSON(400, "m1 bad status")
		return
	}
	defer resp.Body.Close()
	j := apiModel.UserRequest{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	ctx.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if j.Error == nil && j.Response != nil {
		ctx.JSON(200, &j.Response)
		return
	} else {
		setCookieHeader := resp.Header.Get("set-cookie")
		if setCookieHeader != "" {
			fmt.Println("setCOkkie header found")
			ctx.Header().Set("set-cookie", setCookieHeader)
		}
		ctx.JSON(400, "empty response,setting cookie")
		return
	}
}

func getUserID(cookie string) (string, error) {
	if cookie == "" {
		return "", errors.New("empty SID")
	}
	userCookie, apiError := services.DB.GetUserCookie(cookie)
	if apiError != nil || userCookie == nil {
		fmt.Println(apiError)
		userID, err := getUserCookieFromM1(cookie)
		fmt.Println("user_id:", userID, err)
		if userID == "" || err != nil {
			return "", errors.New("USER NOT FOUND")
		} else {
			services.DB.CreateUserCookie(dbModel.UserCookie{Cookie: cookie, PUID: userID})
			return userID, nil
		}
	} else {
		return userCookie.PUID, nil
	}
}

func getUserCookieFromM1(sid string) (string, error) {
	req, err := http.NewRequest("GET", "https://marshalone.ru/api/private/current_user", nil)
	cookie := http.Cookie{Name: "connect.sid", Value: sid}
	req.AddCookie(&cookie)
	if err != nil {
		log.Println("LOG_ERROR: HTTP GET ERROR", err.Error())
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("http get error")
	}
	defer resp.Body.Close()
	j := apiModel.UserRequest{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		fmt.Println(err)
	}
	if j.Error == nil && j.Response != nil {
		fmt.Println("3", j.Response.UserID)
		return j.Response.UserID, nil
	} else {
		if j.Error != nil{
			fmt.Println("4", *j.Error)
		}else{
			fmt.Println("4 response nil")
		}
		return "", errors.New("wrong SID")
	}
}

func getUserRaces(ctx *macaron.Context) ([]apiModel.Race, error) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		return nil, err
	}
	fmt.Println("cookie:", ctx.GetCookie("connect.sid"))
	dateFrom, paramExist := params["date_from"]
	if !paramExist {
		fmt.Println("LOG_ERROR: date_from doesn't exist")
		log.Println("LOG_ERROR: date_from doesn't exist")
		return nil, err
	}
	dateTo, paramExist := params["date_to"]
	if !paramExist {
		fmt.Println("LOG_ERROR: date_to doesn't exist")
		log.Println("LOG_ERROR: date_to doesn't exist")
		return nil, err
	}
	var j apiModel.GetRaces
	if dateTo[0] != "" && dateFrom[0] != "" {
		fmt.Println(fmt.Sprintf("https://marshalone.ru/api/public/races?date_from=%s&date_to=%s", dateFrom[0], dateTo[0]))
		url := fmt.Sprintf("https://marshalone.ru/api/public/races?date_from=%s&date_to=%s", dateFrom[0], dateTo[0])
		req, err := http.NewRequest("GET", url, nil)
		cookie := http.Cookie{Name: "connect.sid", Value: ctx.GetCookie("connect.sid")}
		req.AddCookie(&cookie)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("LOG_ERROR: HTTP LOAD RACES ERROR", err.Error())
			fmt.Println("Errored when sending request to the server2")
			return nil, err
		}
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		ctx.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		err = json.Unmarshal(respBody, &j)
		if err != nil {
			log.Println("LOG_ERROR: M1 Race JSON PARSE ERROR", err.Error())
			return nil, err
		} else {
			if j.Error == nil && j.Response != nil {
				return j.Response, nil
			} else {
				return nil, errors.New(*j.Error)
			}
		}
	} else {
		return nil, errors.New("Wrong data ")
	}
}

func getRaceInfo(raceID string) (*apiModel.OutRaceInfo, error) {
	if raceID == "" {
		return nil, errors.New("empty SID")
	}
	race, apiErr := services.DB.GetRaceFromDB(raceID)
	if race != nil && apiErr == nil {
		return &apiModel.OutRaceInfo{RaceID: race.RaceID, Name: race.Name, Date: race.Date, City: race.City}, nil
	} else {
		fmt.Println(fmt.Sprintf("https://marshalone.ru/api/race?ruid=%s", raceID))
		req, err := http.NewRequest("GET", fmt.Sprintf("https://marshalone.ru/api/public/race?ruid=%s", raceID), nil)
		if err != nil {
			log.Println("LOG_ERROR: HTTP GET ERROR", err.Error())
			fmt.Println(err)
			return nil, errors.New("GET RACE ERROR1")
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.New("http get error1")
		}
		defer resp.Body.Close()
		var j = apiModel.GetRace{}
		fmt.Println(j)
		err = json.NewDecoder(resp.Body).Decode(&j)
		fmt.Println(j)
		if j.Response.Name == "" || j.Response.Date == "" {
			return nil, errors.New("http get error2")
		}
		if err == nil && j.Response != nil {
			j.Response.UID = raceID
			fmt.Println(j)
			date, err := time.Parse("2006-01-02", j.Response.Date)
			if err != nil {
				fmt.Println(err.Error())
			}
			race := dbModel.Race{RaceID: j.Response.UID, Date: date, Name: j.Response.Name, City: j.Response.City}
			raceError := services.DB.NewRace(race)
			if raceError != nil {
				fmt.Println(raceError.Message, raceError.Classification)
			}
			return &apiModel.OutRaceInfo{RaceID: race.RaceID, Name: race.Name, Date: race.Date, City: race.City}, nil
		} else {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("ERROR NIL")
				fmt.Println(j)
			}
			return nil, errors.New("GET RACE ERROR2")
		}
	}
}

func getRaceInfoRouter(ctx *macaron.Context) {
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(400, "rUID not found")
		return
	}
	raceID, paramExist := params["ruid"]
	if !paramExist {
		fmt.Println("LOG_ERROR: rUID not found")
		log.Println("LOG_ERROR: rUID not found")
		ctx.JSON(400, "rUID not found")
		return
	}
	raceInfo, err := getRaceInfo(raceID[0])
	if err != nil || raceInfo == nil {
		if err != nil {
			fmt.Println(err)
		}
		ctx.JSON(400, "get race info error")
		return
	} else {
		ctx.JSON(200, raceInfo)
		return
	}
}

func getUserAccessRouter(ctx *macaron.Context) {
	if sid := ctx.GetCookie("connect.sid"); sid == "" {
		ctx.JSON(400, "3")
		return
	}
	params, err := url.ParseQuery(ctx.Req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		log.Println("LOG_ERROR: QUERY PARSE ERROR", err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if raceID, exist := params["ruid"]; exist && raceID[0] != "" {
		access, err := getUserAccess(ctx.GetCookie("connect.sid"), raceID[0])
		if err != nil || access == nil {
			fmt.Println(err)
			fmt.Println("1")
			ctx.JSON(200, apiModel.OutRaceAccess{RaceID: raceID[0], IsMedia: false, IsOrg: false,
				IsMediaConfirmed: false, IsCompetitor: false})
			return
		} else {
			fmt.Println("2")
			ctx.JSON(200, apiModel.OutRaceAccess{RaceID: access.RaceUID, UserID: access.PUID, IsMedia: access.IsMedia, IsOrg: access.IsOrg,
				IsMediaConfirmed: access.IsMediaConfirmed, IsCompetitor: access.IsComp})
			return
		}
	} else {
		ctx.JSON(400, "Wrong raceID")
		return
	}
}

func getUserAccess(cookie string, raceID string) (*dbModel.UserAccess, error) {
	fmt.Println("GET USER ACCESS GET USER ACCESS GET USER ACCESS")
	fmt.Println("cookie:", cookie)
	userID, err := getUserID(cookie)
	if err != nil || userID == "" {
		fmt.Println("empty UserID")
		return nil, err
	}
	userAccess, apiErr := services.DB.GetUserAccess(cookie, raceID)
	if userAccess != nil && apiErr == nil {
		return userAccess, nil
	} else {
		raceInfo, err := getRaceInfo(raceID)
		if err != nil || raceInfo == nil {
			return nil, err
		} else {
			fmt.Println(raceInfo.Date.Format("2006-01-02"))
			urlString := fmt.Sprintf("https://marshalone.ru/api/public/race?ruid=%s", raceID)
			fmt.Println(urlString)
			req, err := http.NewRequest("GET", urlString, nil)
			reqCookie := http.Cookie{Name: "connect.sid", Value: cookie}
			req.AddCookie(&reqCookie)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("LOG_ERROR: HTTP LOAD RACES ERROR", err.Error())
				fmt.Println("Errored when sending request to the server2")
				return nil, err
			}
			fmt.Println(resp.StatusCode)
			defer resp.Body.Close()
			var j apiModel.GetRace
			err = json.NewDecoder(resp.Body).Decode(&j)
			fmt.Println("RACE UID", j.Response.Name)
			fmt.Println("RACE", j.Response)
			fmt.Println("USERID", userID)
			if err == nil && j.Response != nil && j.Error == nil {
				userAccess := dbModel.UserAccess{RaceUID: raceID, PUID: userID, Cookie: cookie, IsComp: j.Response.IsCompetitor,
					IsMedia: j.Response.IsMedia, IsOrg: j.Response.IsOrg, IsMediaConfirmed: j.Response.IsMediaConfirmed}
				services.DB.CreateUserAccess(userAccess)
				return &userAccess, nil
			} else {
				return nil, nil
			}
		}
	}
}
