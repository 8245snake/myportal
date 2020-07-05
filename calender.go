package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	AppScriptURL       = "https://script.google.com/macros/s/AKfycbxavm6qHSZ-0oHqfOBkJDxXWf-IChtMB-bfNmD6YUN4UxqU_JPn/exec"
	AppScriptEvents    = "events"
	AppScriptTasks     = "tasks"
	AppScriptSchedules = "schedules"
)

//client httpクライアント
var client = new(http.Client)

//ResponseCalender カレンダー
type ResponseCalender struct {
	Message string          `json:"message"`
	Events  []CalenderEvent `json:"events"`
}

//SetMessage メッセージをセットする
func (r *ResponseCalender) SetMessage(msg string) {
	r.Message = msg
}

//ResponseToDoList リスト
type ResponseToDoList struct {
	Message string     `json:"message"`
	Tasks   []ToDoTask `json:"tasks"`
}

//SetMessage メッセージをセットする
func (r *ResponseToDoList) SetMessage(msg string) {
	r.Message = msg
}

//ResponseSchedule みんなのスケジュール
type ResponseSchedule struct {
	Message   string     `json:"message"`
	Schedules []Schedule `json:"schedules"`
}

//SetMessage メッセージをセットする
func (r *ResponseSchedule) SetMessage(msg string) {
	r.Message = msg
}

//ResponseFromGAS APIへのレスポンス共通化
type ResponseFromGAS interface {
	SetMessage(string)
}

//Schedule スケジュール
type Schedule struct {
	Name  string `json:"name"`
	Item  string `json:"item"`
	Day   string `json:"day"`
	Color string `json:"color"`
}

//ToDoTask タスク
type ToDoTask struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Timelimit   string `json:"timelimit"`
	Completed   bool   `json:"completed"`
}

//CalenderEvent イベント
type CalenderEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Location    string `json:"location"`
}

//getOAuthToken アクセストークンを取得する
func getOAuthToken() string {
	fileName := "C:/Users/USER/MyDrive/API_KEY.txt"
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}
	return string(bytes)
}

//GetFromGAS GASにリクエストを送信する
func GetFromGAS(endpoint string, date *time.Time) (data ResponseFromGAS) {

	req, err := http.NewRequest("GET", AppScriptURL, nil)
	if err != nil {
		data.SetMessage(err.Error())
		return data
	}
	token := getOAuthToken()
	req.Header.Set("Authorization", "Bearer "+token)

	params := req.URL.Query()
	params.Add("type", endpoint)
	if date != nil {
		params.Add("day", date.Format("2006-01-02"))
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		data.SetMessage(err.Error())
		return data
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		data.SetMessage(err.Error())
		return data
	}

	// fmt.Printf("%s\n", b)

	switch endpoint {
	case AppScriptEvents:
		var events ResponseCalender
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&events)
		data = &events
	case AppScriptTasks:
		var tasks ResponseToDoList
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&tasks)
		data = &tasks
	case AppScriptSchedules:
		var schedule ResponseSchedule
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&schedule)
		data = &schedule
	}

	if err != nil {
		data.SetMessage(err.Error())
		return data
	}
	data.SetMessage("OK")
	return data
}

func serveEvents(w http.ResponseWriter, r *http.Request) {
	events := GetFromGAS(AppScriptEvents, nil)
	jsondata, _ := json.Marshal(events)
	w.Write(jsondata)
}

func serveToDoList(w http.ResponseWriter, r *http.Request) {
	todolist := GetFromGAS(AppScriptTasks, nil)
	jsondata, _ := json.Marshal(todolist)
	w.Write(jsondata)
}

func serveSchedule(w http.ResponseWriter, r *http.Request) {
	schedules := GetFromGAS(AppScriptSchedules, nil)
	jsondata, _ := json.Marshal(schedules)
	w.Write(jsondata)
}
