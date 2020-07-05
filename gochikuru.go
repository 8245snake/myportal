package main

import (
	"encoding/json"
	"net/http"

	"github.com/8245snake/gochikurunow"
)

//ResponseGochikuruMenue メニュー情報
type ResponseGochikuruMenue struct {
	Message  string                 `json:"message"`
	Date     string                 `json:"date"`
	Products []gochikurunow.Product `json:"products"`
}

func getGochikuruMenue() ResponseGochikuruMenue {
	var menueInfo ResponseGochikuruMenue
	menueInfo.Products = []gochikurunow.Product{}
	//クライアントを初期化
	api, err := gochikurunow.NewGochiClient(config.Gochikuru.Mail, config.Gochikuru.Password)
	if err != nil {
		menueInfo.Message = err.Error()
		return menueInfo
	}
	//メニューを取得
	menue, err := api.GetMenu()
	if err != nil {
		menueInfo.Message = err.Error()
		return menueInfo
	}
	//情報を取得
	menueInfo.Message = "OK"
	menueInfo.Date = menue.Date
	menueInfo.Products = menue.Products
	return menueInfo
}

func serveGochikuruMenue(w http.ResponseWriter, r *http.Request) {
	data := getGochikuruMenue()
	jsondata, _ := json.Marshal(data)
	w.Write(jsondata)
}
