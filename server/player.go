package main

import (
    "log"
    "github.com/gorilla/websocket"
)

type Player struct {
    Name string
    Conn *websocket.Conn
    Room *Room
    Score int
}

const DefaultPlayerName = "user"

func CreatePlayer(conn *websocket.Conn) *Player {
    return &Player{
        Name: DefaultPlayerName,
        Conn: conn,
        Score: 0,
    }
}

func (p *Player) Read() (*event, error) {
    e := &event{}
    err := p.Conn.ReadJSON(e)
    if err != nil {
        p.Reset()
        log.Printf("Player Error: %v", err)
    }
    log.Printf("Recieved %v event", e.Name)
    return e, err
}

func (p *Player) Send(e *event) {
    p.Conn.WriteJSON(e)
}

func (p *Player) Reset() {
    if p.InRoom() {
        p.Room.remove_player(p)
    }
    p.Score = 0
}

func (p *Player) InRoom() bool {
    return p.Room != nil
}

func (p *Player) InGame() bool {
    return p.InRoom() && p.Room.Game != nil
}
