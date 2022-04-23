package lottery

type Shop_Registry_RequestData struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductItem struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Count int    `json:"count"`
}

type Product_ItemList_ResponseData struct {
	Items []ProductItem `json:"items"`
}

type ProductImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Buy_Item_ResponseData struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Item    string `json:"item"`
}
