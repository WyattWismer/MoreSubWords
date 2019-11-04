var conn;

function send(event_name, data = "") {
    console.assert(conn);
    console.assert(conn.readyState === WebSocket.OPEN);
    conn.send(JSON.stringify({"event":event_name,
                             "data":data}));
}

function recieve(socket_event) {
    console.log("Recieved data: ", socket_event.data);
    var json = JSON.parse(socket_event.data);

    switch(json["event"]) {
        case "show_page":
            show_page(json);
            break;
        case "room_not_found":
            wrong_code(json);
            break;
        case "game_update":
            update_game(json);
            break;
        case "end_game":
            end_game(json);
            break;
        case "wrong_because":
            wrong_because(json);
            break;
        default:
            console.log("Unsupported event type");
    }
};

window.onload = function () {
    // check for websocket support
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onopen = function (evt) {
            console.log("connection opened");
            send("GET_home");
        };
        conn.onclose = function (evt) {
            console.log("connection closed");
        };

        conn.onmessage = recieve; 
    } else {
        set_body("Your browser does not support WebSockets.");
    }
};

