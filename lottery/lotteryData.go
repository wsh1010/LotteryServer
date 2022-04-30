package lottery

type AppVersion_ResponseData struct {
	Version string `json:"version"`
}

// User data type
type User_Sign_RequestData struct {
	ID_Type  string `json:"id_type"` // kakao, app
	App_ID   string `json:"app_id"`
	Kakao_ID string `json:"kakao_id"`
	Nick     string `json:"nick"`
}

type User_Login_RequestData struct {
	ID_Type  string `json:"id_type"` // kakao, app
	App_ID   string `json:"app_id"`
	Kakao_ID string `json:"kakao_id"`
}

type User_Login_ResponseData struct {
	Nick  string `json:"nick"`
	Valid string `json:"valid"`
}

type Rank_Result_ResponseData struct {
	Round int `json:"round"`
	Rank1 int `json:"rank1"`
	Rank2 int `json:"rank2"`
	Rank3 int `json:"rank3"`
}

type Rank_Result_RequestData struct {
	Rank1 int `json:"rank1"`
	Rank2 int `json:"rank2"`
	Rank3 int `json:"rank3"`
	Ball  int `json:"ball"`
	Coins int `json:"coins"`
}

type Coins_ResponseData struct {
	Coins int    `json:"coins"`
	Nick  string `json:"nick"`
}

type Modify_UserInfo_RequestData struct {
	Nick string `json:"nick"`
}

type Today_Number_ResponseData struct {
	Number int    `json:"number"`
	Time   string `json:"time"`
}

type AddCoins_RequestData struct {
	Coins int `json:"coins"`
}

// Result data type

// Store data type
