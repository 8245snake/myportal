package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const RedmineAPIKey = ""

//ResponseTicket レスポンス形式
type ResponseTicket struct {
	Message string   `json:"message"`
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

//GetTicketTrac tracをスクレイピングする
func GetTicketTrac() (tickets ResponseTicket) {
	req, err := http.NewRequest("GET", "http://orangeright/trac/OR/query", nil)
	if err != nil {
		tickets.Message = err.Error()
		return tickets
	}

	params := req.URL.Query()
	params.Add("owner", "~shingo.hanyu")
	params.Add("status", "assigned")
	params.Add("status", "new")
	params.Add("status", "reopened")
	params.Add("keywords", "~出荷")
	params.Add("milestone", "!33.リリース済")
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		tickets.Message = err.Error()
		return tickets
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		tickets.Message = err.Error()
		return tickets
	}

	doc.Find("table.tickets").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		var ticket Ticket
		ticket.TicketType = TicketTypeTrac
		aTag := row.Find("td.id").Find("a")
		ID := aTag.Text()
		if ID == "" {
			return
		}
		ticket.ID = strings.TrimSpace(strings.Replace(ID, "#", "", -1))
		ticket.URL = "http://orangeright" + aTag.AttrOr("href", "")

		title := row.Find("td.summary").Find("a").Text()
		ticket.Title = strings.TrimSpace(title)

		milestoine := row.Find("td.milestone").Find("a").Text()
		ticket.MineStone = strings.TrimSpace(milestoine)

		ticket.Status = strings.TrimSpace(row.Find("td.status").Text())

		tickets.Tickets = append(tickets.Tickets, ticket)
	})

	return tickets
}

func serveTicketTrac(w http.ResponseWriter, r *http.Request) {
	//パラメータ解析
	// r.ParseForm()
	// form := r.Form
	// startdate := form.Get("start")
	tickets := GetTicketTrac()
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
