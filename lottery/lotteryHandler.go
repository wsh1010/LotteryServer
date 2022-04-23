package lottery

import (
	"LotteryServer/module/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func Handler_userInfo() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(GetUser())
		case http.MethodPost:
			var data User_Sign_RequestData
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNotImplemented)
			} else {
				json.Unmarshal(body, &data)
				result := SaveUser(data)
				if result {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusNotImplemented)
				}
			}
		case http.MethodPut:
			var result int
			accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
			accept_id, exist_id := r.Header["Accept-Id"]
			accept_type, exist_type := r.Header["Id-Type"]
			if !(exist_uuid && exist_id && exist_type) {
				result = http.StatusNotAcceptable
			} else {
				if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
					result = http.StatusNotAcceptable
				} else {
					var data Modify_UserInfo_RequestData
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusNotImplemented)
					} else {
						json.Unmarshal(body, &data)
						modifyNick(data.Nick, accept_id[0], accept_type[0])
						result = http.StatusOK
					}
				}
			}
			w.WriteHeader(result)
		case http.MethodDelete:
			var result int
			accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
			accept_id, exist_id := r.Header["Accept-Id"]
			accept_type, exist_type := r.Header["Id-Type"]
			if !(exist_uuid && exist_id && exist_type) {
				result = http.StatusNotAcceptable
			} else {
				if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
					result = http.StatusNotAcceptable
				} else {
					signOutUser(accept_id[0], accept_type[0])
					result = http.StatusOK
				}
			}
			w.WriteHeader(result)
		}
	}
}

func Handler_userLogin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var data User_Login_RequestData
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNotImplemented)
			} else {
				json.Unmarshal(body, &data)
				id_valid, _ := uuid.NewV4()
				result, nick := LoginUser(data, id_valid.String())
				responsebody_data := User_Login_ResponseData{Nick: nick, Valid: id_valid.String()}
				responseBody, _ := json.Marshal(responsebody_data)
				w.WriteHeader(result)
				w.Write(responseBody)
			}
		case http.MethodDelete:
			var result int
			accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
			accept_id, exist_id := r.Header["Accept-Id"]
			accept_type, exist_type := r.Header["Id-Type"]
			if !(exist_uuid && exist_id && exist_type) {
				result = http.StatusNotAcceptable
			} else {
				if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
					result = http.StatusNotAcceptable
				} else {
					logoutUser(accept_id[0], accept_type[0])
					result = http.StatusOK
				}
			}
			w.WriteHeader(result)
		}
	}
}

func Handler_lotterData() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			responseBody, _ := json.Marshal(lotteryinfo)
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)

		}
	}
}

func Handler_result() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		URISplit := strings.Split(r.RequestURI, "/")
		switch r.Method {
		case http.MethodGet:
			var result int
			if len(URISplit) == 6 {
				if URISplit[5] == "ranks" {
					var data Rank_Result_ResponseData
					result, data = getResultRank(0)
					responseBody, _ := json.Marshal(data)
					w.WriteHeader(result)
					w.Write(responseBody)
				} else {
					var data Rank_Result_ResponseData
					round, err := strconv.Atoi(URISplit[5])
					if err != nil {
						result = http.StatusNotImplemented
					} else {
						result, data = getResultRank(round)
					}
					responseBody, _ := json.Marshal(data)
					w.WriteHeader(result)
					w.Write(responseBody)
				}
			} else if len(URISplit) == 7 {
				if URISplit[5] == "complete" {
					accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
					accept_id, exist_id := r.Header["Accept-Id"]
					accept_type, exist_type := r.Header["Id-Type"]
					if !(exist_uuid && exist_id && exist_type) {
						result = http.StatusNotAcceptable
					} else {
						if len(URISplit) == 7 {
							if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
								result = http.StatusNotAcceptable
							} else {
								round, err := strconv.Atoi(URISplit[6])
								if err != nil {
									result = http.StatusNotFound
								} else {
									result = isComplete(accept_id[0], accept_type[0], round)
								}
							}
						} else {
							result = http.StatusNotFound
						}
					}
					w.WriteHeader(result)
				}
			}

		case http.MethodPut:
			var data Rank_Result_RequestData
			var result int
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNotImplemented)
			} else {
				json.Unmarshal(body, &data)
				accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
				accept_id, exist_id := r.Header["Accept-Id"]
				accept_type, exist_type := r.Header["Id-Type"]
				if !(exist_uuid && exist_id && exist_type) {
					result = http.StatusNotAcceptable
				} else {
					if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
						result = http.StatusNotAcceptable
					} else {
						if len(URISplit) == 6 {
							round, err := strconv.Atoi(URISplit[5])
							if err != nil {
								result = http.StatusNotImplemented
							} else {
								result = setResult(round, accept_id[0], accept_type[0], data)
							}
						}
					}

				}

				w.WriteHeader(result)
			}
		}
	}
}

func Handler_userCoins() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var result int
			var data Coins_ResponseData
			accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
			accept_id, exist_id := r.Header["Accept-Id"]
			accept_type, exist_type := r.Header["Id-Type"]

			if !(exist_uuid && exist_id && exist_type) {
				result = http.StatusNotAcceptable
			} else {
				if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
					result = http.StatusNotAcceptable
				} else {
					result, data = getUsersCoins(accept_id[0], accept_type[0])
				}
			}
			responseBody, _ := json.Marshal(data)
			w.WriteHeader(result)
			w.Write(responseBody)
		}
	}
}

func isComplete(user_id string, id_type string, round int) int {
	var query string
	if id_type == "kakao" {
		query = fmt.Sprintf("SELECT `id` FROM T_USER_INFO WHERE `kakao_id` = '%s';", user_id)
	} else {
		query = fmt.Sprintf("SELECT `id` FROM T_USER_INFO WHERE `app_id` = '%s'; ", user_id)
	}
	var num_id int
	row := db.SelectQueryRow(query)
	row.Scan(&num_id)

	query = fmt.Sprintf("SELECT EXISTS (SELECT * FROM T_USER_COINS WHERE `id` = '%d' AND `round` = '%d');", num_id, round)
	var result int
	row = db.SelectQueryRow(query)
	row.Scan(&result)
	if result == 0 {
		return http.StatusOK
	} else {
		return http.StatusNotImplemented
	}

}

func getUsersCoins(id string, id_type string) (int, Coins_ResponseData) {
	var result int
	var data Coins_ResponseData
	var query string
	if id_type == "kakao" {
		query = fmt.Sprintf("SELECT id, nick, coins FROM T_USER_INFO WHERE kakao_id = '%s';", id)
	} else {
		query = fmt.Sprintf("SELECT id, nick, coins FROM T_USER_INFO WHERE app_id = '%s';", id)
	}

	row := db.SelectQueryRow(query)

	var iden_id int
	row.Scan(&iden_id, &data.Nick, &data.Coins)

	result = http.StatusOK
	return result, data

}

func SaveUser(data User_Sign_RequestData) bool {
	if data.ID_Type == "kakao" {
		var user_id int
		var query string
		if data.App_ID != "" {
			query = fmt.Sprintf("SELECT id FROM T_USER_INFO WHERE app_id = '%s';", data.App_ID)
			row := db.SelectQueryRow(query)
			row.Scan(&user_id)
			now := time.Now().Format("2006-01-02")
			query = fmt.Sprintf("UPDATE T_USER_INFO SET kakao_id='%s', last_login_date='%s' WHERE id=%d;", data.Kakao_ID, now, user_id)
			db.ExcuteQuery(query)
		} else {
			now := time.Now().Format("2006-01-02")
			query = fmt.Sprintf("INSERT INTO T_USER_INFO (nick, kakao_id, last_login_date) VALUES ('%s', '%s', '%s');", data.Nick, data.Kakao_ID, now)
			db.ExcuteQuery(query)
		}
	} else if data.ID_Type == "app" {
		now := time.Now().Format("2006-01-02")
		query := fmt.Sprintf("INSERT INTO T_USER_INFO (nick, app_id, last_login_date) VALUES ('%s', '%s', '%s');", data.Nick, data.App_ID, now)
		db.ExcuteQuery(query)
	} else {
		return false
	}
	return true
}

func GetUser() int {
	return http.StatusOK
}

func LoginUser(data User_Login_RequestData, valid string) (int, string) {
	var nick string
	if data.ID_Type == "kakao" {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM T_USER_INFO WHERE kakao_id = '%s';", data.Kakao_ID)
		row := db.SelectQueryRow(query)
		row.Scan(&count)
		if count == 0 {
			return http.StatusNotImplemented, ""
		}
		query = fmt.Sprintf("SELECT nick FROM T_USER_INFO WHERE kakao_id = '%s';", data.Kakao_ID)
		row = db.SelectQueryRow(query)
		row.Scan(&nick)

		if nick != "" {
			query = fmt.Sprintf("UPDATE T_USER_INFO SET `last_login_date`= now(), login_uuid = '%s' WHERE kakao_id = '%s';", valid, data.Kakao_ID)
			db.ExcuteQuery(query)
		} else {
			return http.StatusNotImplemented, ""
		}
	} else if data.ID_Type == "app" {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM T_USER_INFO WHERE app_id = '%s';", data.App_ID)
		row := db.SelectQueryRow(query)
		row.Scan(&count)
		if count == 0 {
			return http.StatusNotImplemented, ""
		}

		query = fmt.Sprintf("SELECT nick FROM T_USER_INFO WHERE app_id = '%s';", data.App_ID)
		row = db.SelectQueryRow(query)
		row.Scan(&nick)
		if nick != "" {
			query = fmt.Sprintf("UPDATE T_USER_INFO SET `last_login_date`= now(), login_uuid = '%s' WHERE app_id = '%s';", valid, data.App_ID)
			db.ExcuteQuery(query)
		} else {
			return http.StatusNotImplemented, ""
		}
	} else {
		return http.StatusBadRequest, ""
	}
	return http.StatusOK, nick
}

func getResultRank(round int) (int, Rank_Result_ResponseData) {
	var data Rank_Result_ResponseData
	var query string
	if round == 0 {
		query = "SELECT SUM(`rank1`), SUM(`rank2`), SUM(`rank3`) FROM T_RESULT_LOTTO;"
	} else {
		query = fmt.Sprintf("SELECT `rank1`, `rank2`, `rank3` FROM T_RESULT_LOTTO WHERE `round` = '%d'", round)
	}

	row := db.SelectQueryRow(query)
	data.Round = round
	row.Scan(&data.Rank1, &data.Rank2, &data.Rank3)

	return http.StatusOK, data
}
func CheckIdAvailable(id string, user_uuid string, id_type string) bool {
	var query string

	if id_type == "kakao" {
		query = fmt.Sprintf("SELECT login_uuid FROM T_USER_INFO WHERE kakao_id = '%s';", id)
	} else {
		query = fmt.Sprintf("SELECT login_uuid FROM T_USER_INFO WHERE app_id = '%s';", id)
	}
	row := db.SelectQueryRow(query)
	var uuid sql.NullString
	row.Scan(&uuid)
	if uuid.Valid {
		if uuid.String == user_uuid {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func setResult(round int, user_id string, id_type string, data Rank_Result_RequestData) int {
	query := fmt.Sprintf("INSERT INTO T_RESULT_LOTTO (`round`, `rank1`, `rank2`, `rank3`) VALUES ('%d', '%d', '%d', '%d') ON  DUPLICATE KEY UPDATE `rank1` = `rank1` + '%d' , `rank2` = `rank2` + '%d', `rank3` = `rank3` +'%d';",
		round, data.Rank1, data.Rank2, data.Rank3, data.Rank1, data.Rank2, data.Rank3)

	db.ExcuteQuery(query)

	if id_type == "kakao" {
		query = fmt.Sprintf("SELECT id FROM T_USER_INFO WHERE kakao_id = '%s';", user_id)
	} else {
		query = fmt.Sprintf("SELECT id FROM T_USER_INFO WHERE app_id = '%s';", user_id)
	}
	row := db.SelectQueryRow(query)

	var id int
	row.Scan(&id)
	if id == 0 {
		return http.StatusNotImplemented
	}
	query = fmt.Sprintf("INSERT INTO T_USER_COINS VALUES ('%d', '%d', '%d'); ", id, round, data.Coins)
	_, err := db.ExcuteQuery(query)

	query = fmt.Sprintf("UPDATE `T_USER_INFO` SET `coins` = `coins` + %d WHERE `id` = '%d'; ", data.Coins, id)
	db.ExcuteQuery(query)

	if err != nil {
		return http.StatusNotImplemented
	}

	return http.StatusOK
}

func modifyNick(nick string, id string, idType string) {
	var query string
	if idType == "kakao" {
		query = fmt.Sprintf("UPDATE `T_USER_INFO` SET `nick` = '%s' WHERE `kakao_id` = '%s'", nick, id)
	} else {
		query = fmt.Sprintf("UPDATE `T_USER_INFO` SET `nick` = '%s' WHERE `app_id` = '%s'", nick, id)
	}
	db.ExcuteQuery(query)

}

func signOutUser(id string, idType string) {
	var query string

	if idType == "kakao" {
		query = fmt.Sprintf("SELECT `id` FROM `T_USER_INFO` WHERE `kakao_id` = '%s';", id)
		row := db.SelectQueryRow(query)
		var idKey int
		row.Scan(&idKey)
		query = fmt.Sprintf("DELETE FROM `T_USER_COINS` WHERE `id` = '%d'", idKey)
		db.ExcuteQuery(query)
		query = fmt.Sprintf("DELETE FROM `T_USER_INFO` WHERE `kakao_id` = '%s'", id)
	} else {
		query = fmt.Sprintf("SELECT `id` FROM `T_USER_INFO` WHERE `app_id` = '%s';", id)
		row := db.SelectQueryRow(query)
		var idKey int
		row.Scan(&idKey)
		query = fmt.Sprintf("DELETE FROM `T_USER_COINS` WHERE `id` = '%d'", idKey)
		db.ExcuteQuery(query)
		query = fmt.Sprintf("DELETE FROM `T_USER_INFO` WHERE `app_id` = '%s'", id)

	}
	db.ExcuteQuery(query)
}

func logoutUser(id string, idType string) {
	var query string

	if idType == "kakao" {
		query = fmt.Sprintf("UPDATE `T_USER_INFO` set `login_uuid` = null WHERE `kakao_id` = '%s';", id)

	} else {
		query = fmt.Sprintf("UPDATE `T_USER_INFO` set `login_uuid` = null WHERE `app_id` = '%s';", id)
	}
	db.ExcuteQuery(query)
}
