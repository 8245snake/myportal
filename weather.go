package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

//ResponseWeather 天気予報
type ResponseWeather struct {
	Message    string `json:"message"`
	RainMapURL string `json:"rain_map_url"`
}

//GetRainMap 雨雲レーダーの画像URLを取得
func GetRainMap() string {
	URL := "https://weather.yahoo.co.jp/weather/raincloud/13/"
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return ""
	}
	return doc.Find("td.mainImg").First().Find("img").First().AttrOr("src", "")
}

//GetWeatherReportFrame 今日と明日の天気のフレーム
func GetWeatherReportFrame() {
	URL := "https://weather.yahoo.co.jp/weather/jp/13/4410.html"

	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return
	}

	html, err := doc.Find("div.forecastCity").Html()
	if err != nil {
		return
	}

	//新規作成
	file, err := os.OpenFile("./static/etc/iframe_weather.html", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	t := template.Must(template.ParseFiles("template/_iframe_weather.html"))
	if err := t.Execute(writer, html); err != nil {
		log.Fatal(err)
	}
	//書き込み
	writer.Flush()
}

//天気予報
func serveWeatherReport(w http.ResponseWriter, r *http.Request) {
	var data ResponseWeather

	data.RainMapURL = GetRainMap()
	GetWeatherReportFrame()

	data.Message = "正常に取得しました"
	jsondata, _ := json.Marshal(data)
	w.Write(jsondata)
}
