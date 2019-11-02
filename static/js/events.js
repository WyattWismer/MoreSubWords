// helper methods
function set_body(message) {
    var body = document.getElementById("body");
    body.innerHTML = message;
}

// send events
function send_code() {
    var code_input = document.getElementById("code_input");
    send("submit_room_code", code_input.value);
}

function update_name() {
    var inp = document.getElementById('name_input');
    send('set_name', inp.value);
}

function submit_word() {
    var inp = document.getElementById('word_input');
    send('submit_word', inp.value);
    inp.value = "";
}

// recieve events
function show_page(json) {
    var body_text = json["data"]
    set_body(body_text);

    // if this is game screen
    var inp = document.getElementById('word_input');
    if (inp != undefined) {
        inp.addEventListener('keyup', (e) => {
            if (e.key === "Enter") {
                submit_word()
            }
        });
    }
}

function wrong_code(json) {
    var info_box = document.getElementById("info_box");
    info_box.innerHTML = "Can not find a room with this code";
}

function format_seconds(time) {
    var min = parseInt(time/60);
    var sec = time%60;
    var res = '';
    if (min > 0) {
        res += min;
        res += ':';
    }
    res += sec;
    return res; 
}

var clock_time;
var repeater;
function decrement_update() {
    clock_time -= 1
    update_time();
}

function update_time(seconds=clock_time) {
    if (seconds < 0) {
        clearInterval(repeater);
        clock_time = undefined;
        return;
    }
    if (clock_time == undefined) {
        clock_time = seconds;
        repeater = setInterval(decrement_update, 1000);
    }
    var eps = 3;
    if (Math.abs(clock_time - seconds) >= eps) {
        clock_time = seconds;
    }

    var game_time = document.getElementById("game_time");
    if (game_time == undefined) {
        clearInterval(repeater);
        clock_time = undefined;
        return;
    }
    game_time.innerHTML = format_seconds(clock_time);
}

function update_game(json) {
    var new_records = json["data"]["RecordInfo"];
    var records = document.getElementById("game_records");
    records.innerHTML = new_records;

    var time = json["data"]["SecondsLeft"];
    update_time(time);
}

function end_game(json) {
    var header = document.getElementById("announcement");
    header.innerHTML = json["data"];

    var button = document.getElementById("back");
    button.innerHTML = "Return to room";
    button.onclick = function () {
        send('Return_room');
    }
}

