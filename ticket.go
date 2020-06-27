package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

//ResponseTicket レスポンス形式
type ResponseTicket struct {
	Tickets []Ticket `json:"tickets"`
}

//Ticket チケット
type Ticket struct {
	//TicketType タイプ
	TicketType string `json:"ticket_type"`
	//ID 識別子
	ID string `json:"id"`
	//Title タイトル
	Title string `json:"title"`
	//TimeLimit 期限
	TimeLimit string `json:"timelimit"`
	//MineStone マイルストン
	MineStone string `json:"milestone"`
	//Status ステータス
	Status string `json:"status"`
	//URL チケットページのリンク
	URL string `json:"url"`
}

const (
	TicketTypeTrac            = "trac"
	TicketTypeRedmineBug      = "bug"
	TicketTypeRedmineShipment = "shipment"
	TicketTypeRedmineECO      = "eco"
	TicketTypeRedmineBacklog  = "backlog"
)

func serveTicketTrac(w http.ResponseWriter, r *http.Request) {
	//パラメータ解析
	// r.ParseForm()
	// form := r.Form
	// startdate := form.Get("start")

	tickets := ResponseTicket{}
	for i := 0; i < 10; i++ {
		ID := strconv.Itoa(999 + i)
		tickets.Tickets = append(tickets.Tickets, Ticket{ID: ID, TicketType: TicketTypeTrac, Status: "未着手", Title: "チケットタイトル" + ID, TimeLimit: "2020/06/27", MineStone: "V3.1"})
	}
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketRedmineBug(w http.ResponseWriter, r *http.Request) {

	tickets := ResponseTicket{}
	for i := 0; i < 10; i++ {
		ID := strconv.Itoa(999 + i)
		tickets.Tickets = append(tickets.Tickets, Ticket{ID: ID, TicketType: TicketTypeRedmineBug, Status: "未着手", Title: "チケットタイトル" + ID, TimeLimit: "2020/06/27", MineStone: "V3.1"})
	}
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketRedmineShipment(w http.ResponseWriter, r *http.Request) {

	tickets := ResponseTicket{}
	for i := 0; i < 10; i++ {
		ID := strconv.Itoa(999 + i)
		tickets.Tickets = append(tickets.Tickets, Ticket{ID: ID, TicketType: TicketTypeRedmineShipment, Status: "未着手", Title: "チケットタイトル" + ID, TimeLimit: "2020/06/27", MineStone: "V3.1"})
	}
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketRedmineECO(w http.ResponseWriter, r *http.Request) {

	tickets := ResponseTicket{}
	for i := 0; i < 10; i++ {
		ID := strconv.Itoa(999 + i)
		tickets.Tickets = append(tickets.Tickets, Ticket{ID: ID, TicketType: TicketTypeRedmineECO, Status: "未着手", Title: "チケットタイトル" + ID, TimeLimit: "2020/06/27", MineStone: "V3.1"})
	}
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketBacklog(w http.ResponseWriter, r *http.Request) {
	tickets := ResponseTicket{}
	for i := 0; i < 10; i++ {
		ID := strconv.Itoa(999 + i)
		tickets.Tickets = append(tickets.Tickets, Ticket{ID: ID, TicketType: TicketTypeRedmineBacklog, Status: "未着手", Title: "チケットタイトル" + ID, TimeLimit: "2020/06/27", MineStone: "V3.1"})
	}
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}
