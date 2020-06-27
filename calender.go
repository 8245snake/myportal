package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//AppScriptURL エンドポイント
const AppScriptURL = "https://script.google.com/macros/s/youtkey/exec?"

//ResponseCalender カレンダー
type ResponseCalender struct {
	Message string          `json:"message"`
	Events  []CalenderEvent `json:"events"`
}

//CalenderEvent イベント
type CalenderEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Location    string `json:"location"`
}

//getEvents GASにリクエストを送信する
func getEvents(date *time.Time) ResponseCalender {
	var calender ResponseCalender
	resp, err := http.Get(AppScriptURL + "type=events")
	if err != nil {
		calender.Message = err.Error()
		return calender
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		calender.Message = err.Error()
		return calender
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&calender)
	if err != nil {
		calender.Message = err.Error()
		return calender
	}
	calender.Message = "正常に取得しました"
	return calender
}

//getToDoList GASにリクエストを送信する
func getToDoList(date *time.Time) (ResponseCalender, error) {
	var calender ResponseCalender
	resp, err := http.Get(AppScriptURL + "type=tasks")
	if err != nil {
		return calender, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return calender, err
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&calender)
	if err != nil {
		return calender, err
	}

	return calender, nil
}

func serveEvents(w http.ResponseWriter, r *http.Request) {
	events := getEvents(nil)
	jsondata, _ := json.Marshal(events)
	w.Write(jsondata)
}

func serveToDoList(w http.ResponseWriter, r *http.Request) {

}
