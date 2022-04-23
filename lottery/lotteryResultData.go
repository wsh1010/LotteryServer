package lottery

import "time"

var lotteryinfo lotteryInfo

type lotteryInfo struct {
	CurrRound    int              `json:"currRound"`
	IsClosed     bool             `json:"isClosed"`
	ExpectedDate time.Time        `json:"expectedDate"`
	LastData     lotteryRoundData `json:"lastData"`
}

type lotteryRoundData struct {
	ReturnValue string `json:"returnValue"`
	DrwNoDate   string `json:"drwNoDate"`
	DrwNo       int    `json:"drwNo"`
	DrwtNo1     int    `json:"drwtNo1"`
	DrwtNo2     int    `json:"drwtNo2"`
	DrwtNo3     int    `json:"drwtNo3"`
	DrwtNo4     int    `json:"drwtNo4"`
	DrwtNo5     int    `json:"drwtNo5"`
	DrwtNo6     int    `json:"drwtNo6"`
	BnusNo      int    `json:"bnusNo"`
}
