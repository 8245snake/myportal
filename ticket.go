package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//RedmineAPIKey APIキー
const RedmineAPIKey = "6e353c242541c24d7879ca11e60f6a55541a227d"

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

//RedmineAPI APIレスポンス構造体
type RedmineAPI struct {
	Issues     []RedmineIssue `json:"issues"`
	TotalCount int            `json:"total_count"`
	Offset     int            `json:"offset"`
	Limit      int            `json:"limit"`
}

//RedmineIssue 課題
type RedmineIssue struct {
	ID           int                  `json:"id"`
	Project      RedmineField         `json:"project"`
	Tracker      RedmineField         `json:"tracker"`
	Status       RedmineField         `json:"status"`
	Priority     RedmineField         `json:"priority"`
	Author       RedmineField         `json:"author"`
	AssignedTo   RedmineField         `json:"assigned_to"`
	Category     RedmineField         `json:"category"`
	Subject      string               `json:"subject"`
	Description  string               `json:"description"`
	StartDate    string               `json:"start_date"`
	DueDate      string               `json:"due_date"`
	DoneRatio    int                  `json:"done_ratio"`
	CustomFields []RedmineCustomField `json:"custom_fields"`
	CreatedOn    time.Time            `json:"created_on"`
	UpdatedOn    time.Time            `json:"updated_on"`
}

//RedmineField 標準項目
type RedmineField struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//RedmineCustomField カスタム項目
type RedmineCustomField struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Multiple bool   `json:"multiple,omitempty"`
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
	params.Add("col", "id")
	params.Add("col", "summary")
	params.Add("col", "status")
	params.Add("col", "milestone")
	params.Add("col", "due_close")
	params.Add("order", "due_close")

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

		ticket.TimeLimit = strings.TrimSpace(row.Find("td.due_close").Text())

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

func getRedmineTickets(endpoint string) RedmineAPI {
	var data RedmineAPI
	req, err := http.NewRequest("GET", "https://10.212.252.83/redmine/projects/"+endpoint+"/issues.json", nil)
	if err != nil {
		return data
	}

	params := req.URL.Query()
	params.Add("key", RedmineAPIKey)
	params.Add("assigned_to_id", "me")
	params.Add("limit", "100")

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return data
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data
	}

	err = json.NewDecoder(bytes.NewReader(b)).Decode(&data)
	if err != nil {
		return data
	}

	return data
}

func getCustonFieldValue(issue RedmineIssue, name string) string {
	for _, cf := range issue.CustomFields {
		if name == cf.Name {
			return cf.Value
		}
	}
	return ""
}

func serveTicketRedmineBug(w http.ResponseWriter, r *http.Request) {
	apiData := getRedmineTickets("support-request")
	var tickets ResponseTicket

	for _, issue := range apiData.Issues {
		var ticket Ticket
		ticket.TicketType = TicketTypeRedmineBug
		ticket.ID = strconv.Itoa(issue.ID)
		ticket.Title = issue.Subject
		ticket.URL = "https://10.212.252.83/redmine/issues/" + ticket.ID
		ticket.MineStone = getCustonFieldValue(issue, "施設名")
		ticket.TimeLimit = issue.DueDate
		ticket.Status = issue.Status.Name
		tickets.Tickets = append(tickets.Tickets, ticket)
	}

	tickets.Message = "OK"
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketRedmineShipment(w http.ResponseWriter, r *http.Request) {
	apiData := getRedmineTickets("shipping")
	var tickets ResponseTicket

	for _, issue := range apiData.Issues {
		var ticket Ticket
		ticket.TicketType = TicketTypeRedmineShipment
		ticket.ID = strconv.Itoa(issue.ID)
		ticket.Title = issue.Subject
		ticket.URL = "https://10.212.252.83/redmine/issues/" + ticket.ID
		ticket.MineStone = getCustonFieldValue(issue, "施設リスト")
		ticket.Status = issue.Status.Name
		tickets.Tickets = append(tickets.Tickets, ticket)
	}

	tickets.Message = "OK"
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketRedmineECO(w http.ResponseWriter, r *http.Request) {
	apiData := getRedmineTickets("ecr-eco_asakai")
	var tickets ResponseTicket

	for _, issue := range apiData.Issues {
		var ticket Ticket
		ticket.TicketType = TicketTypeRedmineECO
		ticket.ID = getCustonFieldValue(issue, "文書管理番号")
		ticket.Title = issue.Subject
		ticket.URL = "https://10.212.252.83/redmine/issues/" + strconv.Itoa(issue.ID)
		ticket.MineStone = ""
		ticket.TimeLimit = issue.DueDate
		ticket.Status = issue.Status.Name
		tickets.Tickets = append(tickets.Tickets, ticket)
	}

	tickets.Message = "OK"
	jsindata, _ := json.Marshal(tickets)
	w.Write(jsindata)
}

func serveTicketBacklog(w http.ResponseWriter, r *http.Request) {
}
