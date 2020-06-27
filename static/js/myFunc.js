window.onload = function (event) {
    updateAllEvents();
    updateAllTickets();
}

function ZeroPadding(num, digit) {
    return ('00000000' + num).slice(-digit);
}

////////////////////////////////////////////////////////
// カレンダー
////////////////////////////////////////////////////////

//カレンダーイベント取得
function updateAllEvents() {
    //リストを全て消す
    deleteAllEvents()
    //APIから取得
    getTodaysEvents();
}

function getTodaysEvents() {
    var list = document.getElementById("event-list");
    var spinner_id = 'event-spinner';
    document.getElementById('event-spinner').style.visibility = "visible";
    fetch("/api/events")
        .then(response => response.json())
        .then(data => {
            // console.log(data.ticket_type);
            data.events.forEach(event => {
                console.log(event.title);
                var start = new Date(event.start);
                var end = new Date(event.end);
                var time = ZeroPadding(start.getHours(), 2) + ':' + ZeroPadding(start.getMinutes(), 2) +
                    '～' + ZeroPadding(end.getHours(), 2) + ':' + ZeroPadding(end.getMinutes(), 2);
                var item = createListItemNode(time, event.title, event.description, event.location);
                list.appendChild(item);
            });
            document.getElementById(spinner_id).style.visibility = "hidden";
        });
}

function createListItemNode(datetime, title, description, place) {
    //要素の作成
    var element = document.createElement('a');
    element.classList.add("list-group-item");
    element.classList.add("list-group-item-action");
    element.classList.add("flex-column");
    element.classList.add("align-items-start");
    element.setAttribute("data-toggle", "tooltip");
    element.setAttribute("data-placement", "top");
    element.setAttribute("title", description);
    //時刻
    var h5 = document.createElement('h5');
    h5.classList.add("mb-1");
    h5.appendChild(document.createTextNode(datetime));
    //タイトル
    var main = document.createElement('p');
    main.classList.add("mb-1");
    main.appendChild(document.createTextNode(title));
    //場所
    var small = document.createElement('small');
    small.classList.add("text-muted");
    small.appendChild(document.createTextNode(place));
    //子ノードをセット
    element.appendChild(h5);
    element.appendChild(main);
    element.appendChild(small);
    return element;
}

//イベントを全て削除する
function deleteAllEvents() {
    //子要素を全て削除
    var list = document.getElementById("event-list");
    list.textContent = null;
}

////////////////////////////////////////////////////////
// チケット
////////////////////////////////////////////////////////

const TICKET_TYPE_TRAC = "trac";
const TICKET_TYPE_BUG = "bug";
const TICKET_TYPE_SHIPMENT = "shipment";
const TICKET_TYPE_ECO = "eco";
const TICKET_TYPE_BACKLOG = "backlog";



//全てのチケット取得
function updateAllTickets() {
    updateTicketTable(TICKET_TYPE_TRAC, '/api/ticket/trac');
    updateTicketTable(TICKET_TYPE_BUG, '/api/ticket/redmine/bug');
    updateTicketTable(TICKET_TYPE_SHIPMENT, '/api/ticket/redmine/shipment');
    updateTicketTable(TICKET_TYPE_ECO, '/api/ticket/redmine/eco');
    updateTicketTable(TICKET_TYPE_BACKLOG, '/api/ticket/backlog');
}

//Tracのチケットを取得して表に出力する
function updateTicketTable(ticket_type, url) {
    deleteAllRows(ticket_type);
    var spinner_id = ticket_type + '-spinner';
    document.getElementById(spinner_id).style.visibility = "visible";
    fetch(url)
        .then(response => response.json())
        .then(data => {
            // console.log(data.ticket_type);
            data.tickets.forEach(ticket => {
                insertNewRow(ticket);
            });
            document.getElementById(spinner_id).style.visibility = "hidden";
        });
}

//チケット表に行を追加する
function insertNewRow(ticket) {
    var table = document.getElementById(ticket.ticket_type + '-table');
    var row = table.insertRow(-1);

    //IDにはリンクあり
    var link = document.createElement("a");
    link.href = ticket.url;
    link.appendChild(document.createTextNode(ticket.id));

    var col_ID = row.insertCell(-1);
    col_ID.appendChild(link);

    var col_title = row.insertCell(-1);
    col_title.appendChild(document.createTextNode(ticket.title));

    var col_timelimit = row.insertCell(-1);
    col_timelimit.appendChild(document.createTextNode(ticket.timelimit));

    var col_milestone = row.insertCell(-1);
    col_milestone.appendChild(document.createTextNode(ticket.milestone));

    var col_status = row.insertCell(-1);
    col_status.appendChild(document.createTextNode(ticket.status));
}

//指定した表の行を全て削除する（ヘッダ行は除く）
function deleteAllRows(ticket_type) {
    var table = document.getElementById(ticket_type + '-table');
    var row_num = table.rows.length;
    //行数がどんどん変化するため最初に行数を覚えておいてループ内で常に1行目を消すようにする
    for (let index = 1; index <= row_num - 1; index++) {
        table.deleteRow(1);
    }
}