package main

import (
	"log"
	"net/http"
)

//home トップ画面
func home(w http.ResponseWriter, r *http.Request) {

}

func main() {
	port := "4000"

	//CSSとjsにアクセスするために必要
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	//ページ
	http.Handle("/", http.FileServer(http.Dir("page")))
	// http.HandleFunc("/", home)

	// API
	//チケット
	http.HandleFunc("/api/ticket/trac", serveTicketTrac)
	http.HandleFunc("/api/ticket/redmine/bug", serveTicketRedmineBug)
	http.HandleFunc("/api/ticket/redmine/shipment", serveTicketRedmineShipment)
	http.HandleFunc("/api/ticket/redmine/eco", serveTicketRedmineECO)
	http.HandleFunc("/api/ticket/backlog", serveTicketBacklog)
	//カレンダーイベント
	http.HandleFunc("/api/events", serveEvents)
	//TODOリスト
	http.HandleFunc("/api/todo", serveToDoList)
	//天気予報
	http.HandleFunc("/api/weather", serveWeatherReport)
	//みんなのスケジュール
	http.HandleFunc("/api/schedule", serveSchedule)

	//開始
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
