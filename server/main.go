package main

import (
    "log"
    "time"
    "strings"
	"net/http"
    "math/rand"
    "bytes"
    "html/template"
    "github.com/gorilla/websocket"
)

type event struct {
    Name string `json:"event"`
    Data interface{} `json:"data"`
}

func applyTemplate(path string, data interface{}) string {
    var out bytes.Buffer
    tmpl := template.Must(template.ParseFiles(path))
    err := tmpl.Execute(&out, data)
    if err != nil {
        log.Printf("Error templating %v when templating %v", err, path)
    }
    return out.String()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serve_socket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
    log.Println("Opened a websocket to user.")
    p := CreatePlayer(conn)

    for {
        e, err := p.Read()
		if err != nil {
			break
		}

        switch e.Name {
        case "GET_home":
            p.Reset()

            e2 := &event{
                Name: "show_page",
                Data: applyTemplate("static/pages/home.html",p.Name),
            }
            p.Send(e2)
        case "start_game":
            if p.Room == nil {
                log.Println("Illegal action, booting player %v", p.Name)
                break
            }
            p.Room.setup_game()
        case "submit_word":
            if !p.InGame() {
                log.Println("No game running, word submission ignored");
            } else {
                p.Room.Game.submit_word(e.Data.(string), p)
            }
        case "GET_room":
            r := create_room()
            r.add_player(p)
        case "Return_room":
            if (!p.InRoom()) {
                log.Println("Illegal action, booting player");
                break
            } else {
                e2 := p.Room.to_page_event()
                p.Send(e2)
            }
        case "GET_join_room":
            e2 := &event{
                Name: "show_page",
                Data: applyTemplate("static/pages/join_room.html",nil),
            }
            p.Send(e2)
        case "submit_room_code":
            code := e.Data.(string)
            code = strings.ToUpper(code)
            r, ok := room_lookup[code]
            if ok && r.Game == nil {
                r.add_player(p)
                e2 := r.to_page_event()
                p.Send(e2)
            } else {
                e2 := &event{Name:"room_not_found"}
                p.Send(e2)
            }
        case "set_name":
            p.Name = e.Data.(string)
        default:
            log.Println("Unsupported event type")
        }
        log.Printf("Sent response for %v event ", e.Name)
    }
}

func main() {
    // seed rng
    rand.Seed(time.Now().UnixNano())

    // setup globals
    setup_room_lookup()
    setup_word_lookup()
    setup_phrase_lookup()

    // serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    // accept websocket connections
	http.HandleFunc("/ws", serve_socket)

    err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
