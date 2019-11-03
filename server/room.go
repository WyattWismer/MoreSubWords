package main

import (
    "strings"
    "sort"
    "math/rand"
    "log"
)

const RoomCodeLength = 4

var room_lookup map[string]*Room

type Room struct {
    Players []*Player
    Code string
    Game *Game
}

func setup_room_lookup() {
    room_lookup = make(map[string]*Room)
}

func generate_rune() rune {
    ind := rand.Intn(36)
    if ind < 26 {
        return rune('A' + ind)
    }
    return rune(ind - 26 + '0')
}

func generate_room_code() string {
    var sb strings.Builder
    for i := 0; i < RoomCodeLength; i++ {
        sb.WriteRune(generate_rune())
    }
    return sb.String()
}

func get_unused_code() string {
    for {
        c := generate_room_code()
        if _, ok := room_lookup[c]; !ok {
            return c
        }
    }
}

func create_room() *Room {
    code := get_unused_code()
    r := &Room{
        Players: make([]*Player, 0),
        Code: code,
    }
    room_lookup[r.Code] = r
    return r
}

func (r *Room) sort_players() {
    sort.Slice(r.Players, func(i,j int) bool {
        if r.Players[i].Score == r.Players[j].Score {
            return r.Players[i].Name < r.Players[j].Name
        }
        return r.Players[i].Score > r.Players[j].Score
    })
}

func (r *Room) to_page_event() *event {
    r.sort_players()
    return &event{
        Name: "show_page",
        Data: applyTemplate("static/pages/room.html", r),
    }
}

func (r *Room) broadcast(e *event) {
    for _,p := range r.Players {
        p.Send(e)
    }
}

func (r *Room) update_room() {
    e := r.to_page_event()
    r.broadcast(e)
}

func (r *Room) setup_game() {
    g := create_game(r)
    e := g.to_page_event()
    r.broadcast(e)
    g.update_game()
    go g.time_game()
}

func (r *Room) add_player(p *Player) {
    r.Players = append(r.Players, p)
    p.Room = r
    r.update_room()
}

func (r *Room) remove_player(p *Player) {
    // find & delete player
    to_del := -1
    arr := r.Players
    for i,p2 := range arr {
        if p2 == p {
            to_del = i
        }
    }
    if to_del != -1 {
        arr[to_del] = arr[len(arr)-1]
        arr = arr[:len(arr)-1]
        r.Players = arr
    }

    // remove room if empty
    if len(arr) == 0 {
        delete(room_lookup, r.Code)
        log.Printf("Removed room %v", r.Code)
    }
    p.Room = nil
    r.update_room()
}
