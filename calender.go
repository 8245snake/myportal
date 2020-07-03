package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//AppScriptURL エンドポイント
const AppScriptURL = "https://script.google.com/macros/s/AKfycbxavm6qHSZ-0oHqfOBkJDxXWf-IChtMB-bfNmD6YUN4UxqU_JPn/exec"

//client httpクライアント
var client = new(http.Client)

//ResponseCalender カレンダー
type ResponseCalender struct {
	Message string          `json:"message"`
	Events  []CalenderEvent `json:"events"`
}

//ResponseToDoList リスト
type ResponseToDoList struct {
	Message string     `json:"message"`
	Tasks   []ToDoTask `json:"tasks"`
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

//getEvents GASにリクエストを送信する
func getEvents(date *time.Time) (calender ResponseCalender) {

	req, err := http.NewRequest("GET", AppScriptURL, nil)
	if err != nil {
		calender.Message = err.Error()
		return calender
	}
	token := getOAuthToken()
	req.Header.Set("Authorization", "Bearer "+token)

	params := req.URL.Query()
	params.Add("type", "events")
	if date != nil {
		params.Add("day", date.Format("2006-01-02"))
	}
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
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
	calender.Message = "OK"
	return calender
}

//getToDoList GASにリクエストを送信する
func getToDoList() ResponseToDoList {
	var todo ResponseToDoList
	req, err := http.NewRequest("GET", AppScriptURL, nil)
	if err != nil {
		todo.Message = err.Error()
		return todo
	}
	token := getOAuthToken()
	req.Header.Set("Authorization", "Bearer "+token)

	params := req.URL.Query()
	params.Add("type", "tasks")
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		todo.Message = err.Error()
		return todo
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		todo.Message = err.Error()
		return todo
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&todo)
	if err != nil {
		todo.Message = err.Error()
		return todo
	}
	todo.Message = "OK"
	return todo
}

func serveEvents(w http.ResponseWriter, r *http.Request) {
	events := getEvents(nil)
	jsondata, _ := json.Marshal(events)
	w.Write(jsondata)
}

func serveToDoList(w http.ResponseWriter, r *http.Request) {
	todolist := getToDoList()
	jsondata, _ := json.Marshal(todolist)
	w.Write(jsondata)
}
