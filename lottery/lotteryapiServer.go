package lottery

import (
	"LotteryServer/module/db"
	"log"
	"net/http"
	"sync"
	"syscall"
	"time"
)

const api_name = "/api/lottery/"
const version = "v1"

var AppVersion = "1.0.0"

const (
	URI_VERSION      = api_name + version + "/version"
	URI_USER_INFO    = api_name + version + "/user/info"  // POST : 가입 GET : 조회 (ID 찾기 또는 비번) PUT : 수정 DELETE : 탈퇴
	URI_USER_LOGIN   = api_name + version + "/user/login" // PUT : Login
	URI_USER_COINS   = api_name + version + "/user/coins" // GET, PUT
	URI_TODAY_NUMBER = api_name + version + "/today"

	//로또 결과
	URI_RESULT       = api_name + version + "/result/"
	URI_LOTTERY_DATA = api_name + version + "/lotterydata"

	//상품
	URI_PRODUCT_REGISTRY = api_name + version + "/shop/product"
	URI_PRODUCT_IMAGE    = api_name + version + "/shop/product/image/"
	URI_PRICE_IMAGE      = api_name + version + "/shop/price/image/"
	URI_BUY_ITEM         = api_name + version + "/shop/product/buy/"
	//URI_LOTTERY_DATA = api_name + version + "/lotterydata"
	//URI_LOTTERY_DATA = api_name + version + "/lotterydata"
)

func RunningServer(wg *sync.WaitGroup, done chan int) {
	if wg != nil {
		defer wg.Done()
	}
	server := &http.Server{
		Addr:    ":44001",
		Handler: nil,
	}
	db.InitDB()
	var server_wg sync.WaitGroup
	server_wg.Add(1)
	go OpenServer(&server_wg, server)

	for len(done) == 1 {
		time.Sleep(1 * time.Second)
	}
	server.Close()
	server_wg.Wait()
}

func OpenServer(wg *sync.WaitGroup, server *http.Server) {
	if wg != nil {
		defer wg.Done()
	}

	http.HandleFunc(URI_VERSION, Handler_getVersion())
	// 유저 관련
	http.HandleFunc(URI_USER_INFO, Handler_userInfo())
	http.HandleFunc(URI_USER_LOGIN, Handler_userLogin())
	http.HandleFunc(URI_USER_COINS, Handler_userCoins())
	http.HandleFunc(URI_TODAY_NUMBER, Handler_today())

	//결과 관련
	http.HandleFunc(URI_RESULT, Handler_result())

	http.HandleFunc(URI_LOTTERY_DATA, Handler_lotterData())

	http.HandleFunc(URI_PRODUCT_REGISTRY, Handler_product())
	http.HandleFunc(URI_PRODUCT_IMAGE, Handler_product_image())
	http.HandleFunc(URI_PRICE_IMAGE, Handler_price_image())
	http.HandleFunc(URI_BUY_ITEM, Handler_buy_item())

	//db 관련
	log.Println("Server is running...")
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Failed to listen : ", err)
		return
	}
	log.Println("Server is ended.")
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

}
