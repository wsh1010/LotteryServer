package lottery

import (
	"LotteryServer/module/db"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Handler_product() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var data Shop_Registry_RequestData
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			json.Unmarshal(body, &data)
			result := registry_product(data)
			w.WriteHeader(result)
		case http.MethodGet:
			var datas Product_ItemList_ResponseData
			datas.Items = make([]ProductItem, 0)
			result := getProductList(&datas)
			responseBody, _ := json.Marshal(datas)
			w.WriteHeader(result)
			w.Write(responseBody)
		}

	}
}

func Handler_product_image() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			paths := strings.Split(r.URL.Path, "/")
			uploadFile, header, err := r.FormFile("filename")
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			dirname := "./item"
			os.MkdirAll(dirname, 0777)
			filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
			file, err := os.Create(filepath)
			defer file.Close()
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			io.Copy(file, uploadFile)
			registry_product_image(paths[len(paths)-1], filepath)

			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			paths := strings.Split(r.URL.Path, "/")
			result, files := getImage(paths[len(paths)-1])
			w.WriteHeader(result)
			w.Write(files)
		}
	}
}

func Handler_buy_item() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			paths := strings.Split(r.URL.Path, "/")
			accept_uuid, exist_uuid := r.Header["Accept-Uuid"]
			accept_id, exist_id := r.Header["Accept-Id"]
			accept_type, exist_type := r.Header["Id-Type"]
			if !(exist_uuid && exist_id && exist_type) {
				w.WriteHeader(http.StatusNotAcceptable)
				log.Println("request error.")
				return
			} else {
				if !CheckIdAvailable(accept_id[0], accept_uuid[0], accept_type[0]) {
					w.WriteHeader(http.StatusNotAcceptable)
					log.Println("not check available")
					return
				} else {
					log.Println("buy item ,", accept_id[0], accept_type[0], paths[len(paths)-1])
					result, data := buyItem(paths[len(paths)-1], accept_id[0], accept_type[0])
					responseBody, _ := json.Marshal(data)
					w.WriteHeader(result)
					w.Write(responseBody)
				}
			}

		}
	}
}

func Handler_price_image() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			paths := strings.Split(r.URL.Path, "/")
			uploadFile, header, err := r.FormFile("filename")
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			dirname := "./item"
			os.MkdirAll(dirname, 0777)
			filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
			file, err := os.Create(filepath)
			defer file.Close()
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			io.Copy(file, uploadFile)
			registry_price_image(paths[len(paths)-1], filepath)

			w.WriteHeader(http.StatusOK)
		}
	}
}

func registry_product(data Shop_Registry_RequestData) int {
	query := fmt.Sprintf("SELECT count(*) FROM `T_POINT_SHOP` WHERE `item_name` = '%s';", data.Name)
	row := db.SelectQueryRow(query)
	var count int
	row.Scan(&count)
	if count != 0 {
		return http.StatusFound
	}
	query = fmt.Sprintf("INSERT INTO `T_POINT_SHOP` (`item_name`, `item_price`) VALUES ('%s', '%d'); ", data.Name, data.Price)
	db.ExcuteQuery(query)

	return http.StatusOK
}

func registry_product_image(item string, path string) int {
	query := fmt.Sprintf("UPDATE `T_POINT_SHOP` SET `item_image` = '%s' WHERE `item_name` = '%s'; ", path, item)
	db.ExcuteQuery(query)

	return http.StatusOK
}

func registry_price_image(item string, path string) int {
	query := fmt.Sprintf("SELECT `item_id` FROM `T_POINT_SHOP` WHERE `item_name` = '%s';", item)
	row := db.SelectQueryRow(query)
	id := 0
	row.Scan(&id)
	if id == 0 {
		return http.StatusNotImplemented
	}
	query = fmt.Sprintf("INSERT INTO `T_SHOP_PRICE` (`item_id` ,`item_price_image`) VALUES ('%d', '%s'); ", id, path)
	db.ExcuteQuery(query)

	return http.StatusOK
}

func getProductList(datas *Product_ItemList_ResponseData) int {
	query := "SELECT `T_POINT_SHOP`.*, COUNT(`T_SHOP_PRICE`.`item_id`) AS `item_count` FROM `T_POINT_SHOP` LEFT JOIN `T_SHOP_PRICE` ON `T_POINT_SHOP`.`item_id` = `T_SHOP_PRICE`.`item_id` AND `T_SHOP_PRICE`.`buyer_id` IS NULL GROUP BY `T_POINT_SHOP`.`item_id`;"
	rows, err := db.SelectQueryRows(query)
	if err != nil {
		return http.StatusNotImplemented
	}

	for rows.Next() {
		var data ProductItem
		var id int
		var path string

		rows.Scan(&id, &data.Name, &data.Price, &path, &data.Count)
		datas.Items = append(datas.Items, data)
	}

	return http.StatusOK
}

func getImage(name string) (int, []byte) {
	query := fmt.Sprintf("SELECT `item_image` FROM `T_POINT_SHOP` WHERE item_name = '%s';", name)
	row := db.SelectQueryRow(query)

	var path string
	row.Scan(&path)

	filebytes, _ := ioutil.ReadFile(path)

	/*enc := base64.StdEncoding.EncodeToString(filebytes)

	data.Name = name
	data.Image = enc*/

	return http.StatusOK, filebytes
}

func buyItem(itemName string, id string, idType string) (int, Buy_Item_ResponseData) {
	var data Buy_Item_ResponseData
	var query string
	if idType == "kakao" {
		query = fmt.Sprintf("SELECT `id`, `coins` FROM `T_USER_INFO` WHERE `kakao_id` = '%s';", id)
	} else {
		query = fmt.Sprintf("SELECT `id`, `coins` FROM `T_USER_INFO` WHERE `app_id` = '%s';", id)
	}
	row := db.SelectQueryRow(query)
	var user_id, user_coins int
	row.Scan(&user_id, &user_coins)

	query = fmt.Sprintf("SELECT `item_id`, `item_price` FROM `T_POINT_SHOP` WHERE `item_name` = '%s';", itemName)
	row = db.SelectQueryRow(query)
	var item_id, item_price int
	row.Scan(&item_id, &item_price)

	if item_price > user_coins {
		data.Result = 1001
		data.Message = "Lack of Money"
		return http.StatusOK, data
	}

	query = fmt.Sprintf("SELECT `item_price_image` FROM `T_SHOP_PRICE` WHERE `item_id` = '%d' and `buyer_id` IS NULL;", item_id)
	row = db.SelectQueryRow(query)
	var price string
	row.Scan(&price)

	if price == "" {
		data.Result = 1002
		data.Message = "price is empty"
		return http.StatusOK, data
	}

	user_coins -= item_price

	query = fmt.Sprintf("UPDATE `T_USER_INFO` SET `coins` = '%d' WHERE `id` = '%d';", user_coins, user_id)
	db.ExcuteQuery(query)

	now := time.Now().Format("2006-01-02")
	query = fmt.Sprintf("UPDATE `T_SHOP_PRICE` SET `buyer_id`='%d', `buy_date` = '%s' WHERE `item_price_image` = '%s';", user_id, now, price)
	db.ExcuteQuery(query)

	filebytes, _ := ioutil.ReadFile(price)
	enc := base64.StdEncoding.EncodeToString(filebytes)
	data.Result = 1000
	data.Message = "ok"
	data.Item = enc

	return http.StatusOK, data
}
