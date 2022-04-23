package lottery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func RunningClient(wg *sync.WaitGroup, done chan int) {
	if wg != nil {
		defer wg.Done()
	}

	lotteryinfo.CurrRound = 1005

	initLotteryInfo(done)

	need_result := false
	checkTime := time.Now()
	for len(done) == 1 {
		now := time.Now()
		if need_result {
			if checkTime.Sub(now).Seconds() < 0 {
				if checkLotteryResult() {
					need_result = false
					log.Println(lotteryinfo.CurrRound)
					log.Println("Show Result.")
				} else {
					checkTime = time.Now().Add(1 * time.Minute)
				}
			}
		} else {
			if checkTime.Sub(now).Seconds() < 0 {
				if now.Format("2006-01-02") == lotteryinfo.LastData.DrwNoDate {
					lotteryinfo.IsClosed = true
				} else {
					lastdate, _ := time.Parse("2006-01-02", lotteryinfo.LastData.DrwNoDate)
					if now.Weekday() == time.Saturday {
						if now.Hour() >= 20 {
							need_result = true
							lotteryinfo.IsClosed = true
							log.Println("wait for result.")
						}
					} else if (now.Sub(lastdate).Hours() / 24) > 7 {
						need_result = true
						log.Println("wait for result.")
					}

					if lotteryinfo.IsClosed {
						if lastdate != time.Now() {
							lotteryinfo.IsClosed = false
							log.Println("Open Lotto.")
						}
					}
					checkTime = time.Now().Add(1 * time.Second)
				}
			}

		}

		time.Sleep(500 * time.Millisecond)
	}
}

func initLotteryInfo(done chan int) {
	for len(done) == 1 {
		url := fmt.Sprintf("https://www.dhlottery.co.kr/common.do?drwNo=%d&method=getLottoNumber", lotteryinfo.CurrRound)
		resp, _ := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		var data lotteryRoundData
		err := json.Unmarshal(body, &data)

		if err != nil {
			log.Println(body[len(body)-1])
			time.Sleep(1 * time.Minute)
			continue
		}

		if data.ReturnValue == "fail" {
			lastdate, _ := time.Parse("2006-01-02", lotteryinfo.LastData.DrwNoDate)
			d, _ := time.ParseDuration("24h")
			lotteryinfo.ExpectedDate = lastdate.Add(d * 7)
			log.Println("expect : ", lotteryinfo.ExpectedDate)
			log.Println("result : ", lastdate)
			break
		}
		lotteryinfo.LastData = data
		lotteryinfo.CurrRound++
	}
}

func checkLotteryResult() bool {
	url := fmt.Sprintf("https://www.dhlottery.co.kr/common.do?drwNo=%d&method=getLottoNumber", lotteryinfo.CurrRound)
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var data lotteryRoundData
	err := json.Unmarshal(body, &data)
	if err != nil {
		return false
	}
	if data.ReturnValue == "success" {
		lotteryinfo.LastData = data
		lotteryinfo.CurrRound++
		lastdate, _ := time.Parse("2006-01-02", lotteryinfo.LastData.DrwNoDate)
		d, _ := time.ParseDuration("24h")
		lotteryinfo.ExpectedDate = lastdate.Add(d * 7)
		log.Println("expect : ", lotteryinfo.ExpectedDate)
		log.Println("result : ", lastdate)
		return true
	}
	return false
}
