
//window.onloadのタイミングでは遅いのでDOMを読み終わったらすぐに実行する
document.addEventListener("DOMContentLoaded", init);

//初期化処理
function init(){
    updateAllEvents();
    updateAllTickets();
    updateWeather();
    updateToDO();
    updateSchedule();
    updateLunch();
}

//ゼロ埋め
function ZeroPadding(num, digit) {
    return ('00000000' + num).slice(-digit);
}

// 指定ミリ秒間だけループさせる
function sleep(waitMsec) {
    var startMsec = new Date();
    while (new Date() - startMsec < waitMsec);
  }

////////////////////////////////////////////////////////
// 弁当
////////////////////////////////////////////////////////

function updateLunch(){
    //削除
    document.getElementById('lunch').textContent = null;
    //スピナー
    const spinner_id = "lunch-spinner";
    var spinner = document.getElementById(spinner_id);
    spinner.style.visibility = "visible";

    fetch("/api/gochikuru")
    .then(response => response.json())
    .then(data => {
        data.products.forEach(product => {
            insertLunchItem(product);
        });
        spinner.style.visibility = "hidden";
    }).catch(function(err){
        console.log(err);
        spinner.style.visibility = "hidden";
    });
}

//リストの要素作成
function insertLunchItem(product){
    //販売元
    var maker_tag = document.createElement("small");
    maker_tag.appendChild(document.createTextNode(product.maker));
    //商品名
    var name_tag = document.createElement("h4");
    name_tag.classList.add("mt-0");
    name_tag.classList.add("mb-1");
    name_tag.appendChild(document.createTextNode(product.name));
    //値段
    var price_tag = document.createElement("h5");
    price_tag.appendChild(document.createTextNode(product.price));
    //ボディ
    var media_body = document.createElement("div");
    media_body.classList.add("media-body");
    media_body.appendChild(maker_tag);
    media_body.appendChild(name_tag);
    media_body.appendChild(price_tag);
    //注文済みラベルの追加
    if(product.amount > 0) {
        var badge_tag = document.createElement("span");
        badge_tag.classList.add("badge");
        badge_tag.classList.add("badge-danger");
        badge_tag.appendChild(document.createTextNode("注文済み"));
        var ordered_tag = document.createElement("h1");
        ordered_tag.appendChild(badge_tag);
        ordered_tag.appendChild(document.createTextNode(product.amount + '個'));
        media_body.appendChild(ordered_tag);
    }

    //画像
    var img_tag = document.createElement("img");
    img_tag.classList.add("d-flex");
    img_tag.classList.add("ml-3");
    img_tag.setAttribute("src", product.image_url);
    img_tag.width = 300;

    //親要素
    var media_tag = document.createElement("div");
    media_tag.classList.add("media");
    media_tag.appendChild(media_body);
    media_tag.appendChild(img_tag);

    //追加
    document.getElementById('lunch').appendChild(media_tag);
}

////////////////////////////////////////////////////////
// スケジュール
////////////////////////////////////////////////////////

function updateSchedule(){

    const spinner_id = "schedule-spinner";
    var spinner = document.getElementById(spinner_id);
    spinner.style.visibility = "visible";
    //削除
    deleteAllScheduleTable()

    fetch("/api/schedule")
    .then(response => response.json())
    .then(data => {
        const table_max = 4;
        const rows_max = Math.floor(data.schedules.length / table_max);
        var count = 1;
        data.schedules.forEach(schedule => {
            var index = Math.ceil(count / rows_max);
            // index = (index < 1) ? 1 : index;
            index = (index > 4) ? 4 : index;
            insertScheduleRow(index, schedule.name, schedule.item, schedule.color);
            count++;
        });
        spinner.style.visibility = "hidden";
    }).catch(function(err){
        console.log(err);
        spinner.style.visibility = "hidden";
    });
}

//スケジュール表の中身を消す（タイトルも）
function deleteAllScheduleTable() {
    const max = 4;
    for (let index = 1; index <= max; index++) {
        var table = document.getElementById('schedule-table-' + index);
        table.textContent = null;
    }
}

//Schedule表に行を追加する
function insertScheduleRow(index, name, item, color) {
    var table = document.getElementById('schedule-table-' + index);
    var row = table.insertRow(-1);

    var col_name = row.insertCell(-1);
    col_name.appendChild(document.createTextNode(name));
    col_name.setAttribute("width","50%");

    var col_item = row.insertCell(-1);
    col_item.appendChild(document.createTextNode(item));
    col_item.style.backgroundColor = color;
    col_name.setAttribute("width","50%");
}


////////////////////////////////////////////////////////
// ToDoリスト
////////////////////////////////////////////////////////

function updateToDO(){
    var list = document.getElementById("todo-list");
    //子要素を全て削除
    list.textContent = null;
    var spinner_id = 'todo-spinner';
    document.getElementById('todo-spinner').style.visibility = "visible";

    fetch("/api/todo")
    .then(response => response.json())
    .then(data => {
        var count = 1;
        data.tasks.forEach(task => {
            var item = createTiDoItem(task.id, task.title, task.timelimit, task.description);
            item.id = "todo-item-" + count;
            list.appendChild(item);
            count++;
        });
        document.getElementById(spinner_id).style.visibility = "hidden";
    }).catch(function(){
        var item = createTiDoItem("", "エラーが発生しました", "", "");
        item.id = "todo-item-error"
        list.appendChild(item);
        document.getElementById(spinner_id).style.visibility = "hidden";
    });
}

//ノードを作成
function createTiDoItem(ID, title, timelimit, description) {

    var h5_title = document.createElement('h5');
    h5_title.classList.add("mb-1");
    h5_title.innerText = title;

    var small_timelimit = document.createElement('small');
    var datetime = new Date(timelimit);
    if (datetime != 'Invalid Date'){
        datetime = new Date(datetime.toLocaleString({ timeZone: 'Asia/Tokyo' }));
        small_timelimit.innerText = '期限：' + datetime.getFullYear() + '/' + (datetime.getMonth() + 1) + '/' + datetime.getDate();
    }else{
        small_timelimit.innerText = '期限：' + ((timelimit != "")? timelimit:"なし");
    }
    
    var div = document.createElement('div');
    div.classList.add("d-flex");
    div.classList.add("w-100");
    div.classList.add("justify-content-between");
    div.appendChild(h5_title);
    div.appendChild(small_timelimit);

    var p_description  = document.createElement('p');
    p_description.classList.add("mb-1");
    p_description.innerText = description;

    var hidden_ID = document.createElement('input');
    hidden_ID.value = ID;
    hidden_ID.hidden = true;

    var element = document.createElement('a');
    element.classList.add("list-group-item");
    element.classList.add("list-group-item-action");
    element.classList.add("flex-column");
    element.classList.add("align-items-start");

    element.appendChild(div);
    element.appendChild(p_description);
    element.appendChild(hidden_ID);
    element.ondblclick = function(e){
        console.log(e.path);
        $('#todo-modal').modal('show');
    };
    return element;
}

//登録
function regTask(){
    var title = document.getElementById('todo-modal-title');
    var limit = document.getElementById('todo-modal-limit');
    var description = document.getElementById('todo-modal-description');
    //スピナー
    const spinner_id = "todo-modal-spinner";
    var spinner = document.getElementById(spinner_id);
    spinner.removeAttribute("hidden");

    var data = new Object();
    data.type = "tasks";
    var tasks = [];
    var task = new Object();
    task.title = title.value;
    task.description = description.value;
    task.timelimit = limit.value;
    tasks.push(task);
    data.tasks = tasks;

    const method = "POST";
    const body = JSON.stringify(data);
    const headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };
    fetch("/api/todo", {
        method,
        headers,
        body
    })
    .then(response => response.json())
    .then(data => {
        spinner.setAttribute("hidden","true");
        //modalを消す
        title.value = "";
        limit.value = "";
        description.value = "";
        $('#todo-modal').modal('hide');
        updateToDO();
    }).catch(function(err){
        console.log(err);
        spinner.setAttribute("hidden","true");
        //modalを消す
        title.value = "";
        limit.value = "";
        description.value = "";
        $('#todo-modal').modal('hide');
    });
}


////////////////////////////////////////////////////////
// 天気予報
////////////////////////////////////////////////////////

//天気予報取得
function updateWeather(){
    const src_weather = "../static/etc/iframe_weather.html";
    const src_loading = "../static/etc/iframe_loading.html";
    var frame = document.getElementById("weather-frame");
    frame.contentDocument.location.replace(src_loading);

    fetch("/api/weather")
        .then(response => response.json())
        .then(data => {
            frame.contentDocument.location.replace(src_weather);
        }).catch(function(err){
            console.log(err);
            frame.contentDocument.location.replace(src_weather);
        });
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
    document.getElementById(spinner_id).style.visibility = "visible";

    fetch("/api/events")
        .then(response => response.json())
        .then(data => {
            var count = 1;
            data.events.forEach(event => {
                var start = new Date(event.start);
                var end = new Date(event.end);
                var time = ZeroPadding(start.getHours(), 2) + ':' + ZeroPadding(start.getMinutes(), 2) +
                    '～' + ZeroPadding(end.getHours(), 2) + ':' + ZeroPadding(end.getMinutes(), 2);
                var item = createListItemNode(time, event.title, event.description, event.location);
                item.id = "event-item-" + count;
                list.appendChild(item);
                count++;
            });
            document.getElementById(spinner_id).style.visibility = "hidden";
        }).catch(function(){
            var item = createListItemNode("エラーが発生しました", "", "", "");
            item.id = "event-item-error"
            list.appendChild(item);
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

    element.ondblclick = function(e){
        console.log(e.path);
    }
    return element;
}

//イベントを全て削除する
function deleteAllEvents() {
    //子要素を全て削除
    var list = document.getElementById("event-list");
    list.textContent = null;
}

//開始時間に加算する
function addMinute(minute){
    var today = new Date();
    var startTimeTag = document.getElementById("event-modal-starttime");
    var startTime = startTimeTag.value;
    if (startTime == ""){
        startTime = ZeroPadding(today.getHours(), 2) + ":" + ZeroPadding(today.getMinutes(), 2);
        startTimeTag.value = startTime;
    }
    
    var startDate = new Date(today.getFullYear() + "-" + (today.getMonth()+1) + "-" + today.getDate() + " " + startTimeTag.value);
    startDate.setMinutes(startDate.getMinutes() + minute);
    document.getElementById("event-modal-endtime").value = ZeroPadding(startDate.getHours(), 2) + ":" + ZeroPadding(startDate.getMinutes(), 2);
}

//新規登録
function regEvent(){
    //入力欄
    var startTime = document.getElementById("event-modal-starttime");
    var endTime = document.getElementById("event-modal-endtime");
    var title = document.getElementById('event-modal-title');
    var description = document.getElementById('event-modal-description');
    //スピナー
    const spinner_id = "event-modal-spinner";
    var spinner = document.getElementById(spinner_id);
    spinner.removeAttribute("hidden");

    var data = new Object();
    data.type = "events";
    var events = [];
    var event = new Object();
    event.title = title.value;
    event.description = description.value;
    var today = new Date();
    event.start = today.getFullYear() + "-" + (today.getMonth()+1) + "-" + today.getDate() + " " + startTime.value;
    event.end = today.getFullYear() + "-" + (today.getMonth()+1) + "-" + today.getDate() + " " + endTime.value;
    events.push(event);
    data.events = events;

    const method = "POST";
    const body = JSON.stringify(data);
    const headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };
    fetch("/api/todo", {
        method,
        headers,
        body
    })
    .then(response => response.json())
    .then(data => {
        spinner.setAttribute("hidden","true");
        //modalを消す
        title.value = "";
        description.value = "";
        startTime.value = "";
        endTime.value = "";
        $('#event-modal').modal('hide');
        updateAllEvents();
    }).catch(function(err){
        console.log(err);
        spinner.setAttribute("hidden","true");
        //modalを消す
        title.value = "";
        limit.value = "";
        description.value = "";
        $('#event-modal').modal('hide');
    });
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
        }).catch(function(){
            var ticket = new Object
            ticket.ticket_type = ticket_type;
            ticket.url = "";
            ticket.title = "エラーが発生しました";
            ticket.timelimit = "";
            ticket.milestone = "";
            ticket.status = "";
            insertNewRow(ticket);
            document.getElementById(spinner_id).style.visibility = "hidden";
        });
}

//チケット表に行を追加する
function insertNewRow(ticket) {
    var table = document.getElementById(ticket.ticket_type + '-table');
    var row = table.insertRow(-1);

    //IDにはリンクあり
    var id_node = new Object;
    if (ticket.url != "") {
        var link = document.createElement("a");
        link.href = ticket.url;
        link.appendChild(document.createTextNode(ticket.id));
        id_node = link;
    }else{
        id_node = document.createTextNode("");
    }

    var col_ID = row.insertCell(-1);
    col_ID.appendChild(id_node);

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