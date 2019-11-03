package main

import (
    "strings"
    "strconv"
    "time"
    "sort"
    "unicode"
)

type Record struct {
    Player *Player
    Words []string
    Score int
}

type Game struct {
    Records []*Record
    Phrase string
    IsUsed map[string]bool
    SecondsLeft int
    Room *Room
}

func create_game(r *Room) *Game {
    g := &Game{
        Phrase: get_phrase(),
        Records: make([]*Record,0),
        IsUsed: make(map[string]bool),
        SecondsLeft: 90,
    }
    g.add_players(r.Players)
    g.add_phrase_words()
    r.Game = g
    g.Room = r
    return g
}

func (g *Game) add_players(players []*Player) {
    for _,p := range players {
        rec := &Record{}
        rec.Player = p
        rec.Words = make([]string,0)
        rec.Score = 0
        g.Records = append(g.Records, rec)
    }
}

type GameUpdate struct {
    SecondsLeft int
    RecordInfo string
}

func (g *Game) create_game_update() *GameUpdate {
    gu := &GameUpdate{}
    gu.SecondsLeft = g.SecondsLeft
    gu.RecordInfo = applyTemplate("static/pages/game_records.html", g)
    return gu
}


func (g *Game) update_game() {
    e := &event{}
    e.Name = "game_update"
    e.Data = g.create_game_update()
    g.Room.broadcast(e)
}

func (g *Game) end_game() {
    for _,p := range g.Room.Players {
        for rank,rec := range g.Records {
            if rec.Player == p {
                p.Score += rec.Score
                e := &event{}
                e.Name = "end_game"
                e.Data = "Your final rank was " + strconv.Itoa(rank+1)
                p.Conn.WriteJSON(e)
                break
            }
        }
    }
    g.Room.Game = nil
}

func (g *Game) time_game() {
    for ; g.SecondsLeft > 0; g.SecondsLeft -= 1 {
        time.Sleep(time.Second)
        if g.SecondsLeft % 10 == 0 {
            g.update_game()
        }
    }
    g.end_game()
}

func (g *Game) word_valid(word string) bool {
    return is_word(word) && is_substring(word, g.Phrase)
}

func (g *Game) word_pts(word string) int {
    result := 100
    for i := 0; i < len(word) - 1; i++ {
        result *= 2
    }
    return result
}

func (g *Game) is_used(word string) bool {
    _, ok := g.IsUsed[word]
    if ok {
        return true
    }
    return false
}

func (g *Game) submit_word(word string, p *Player) {
    word = strings.ToLower(word)
    if (g.word_valid(word) && !g.is_used(word)) {
        // mark word as used
        g.IsUsed[word] = true
        // find record
        ind := -1
        for i,rec := range g.Records {
            if rec.Player == p {
                ind = i
                break
            }
        }
        g.Records[ind].Score += g.word_pts(word)
        g.Records[ind].Words = append(g.Records[ind].Words, word)
        sort.Slice(g.Records, func(i,j int) bool {
            return g.Records[i].Score > g.Records[j].Score
        })
        g.update_game()
    }
}

func (g *Game) to_page_event() *event {
    e := &event{}
    e.Name = "show_page"
    e.Data = applyTemplate("static/pages/game.html", g)
    return e
}

func (g *Game) add_phrase_words() {
    phrase := strings.ToLower(g.Phrase)
    var sb strings.Builder
    for _,r := range phrase {
        if unicode.IsLetter(r) {
            sb.WriteRune(r)
        } else {
            word := sb.String()
            g.IsUsed[word] = true
            sb.Reset()
        }
    }

    if sb.Len() > 0 {
        word := sb.String()
        g.IsUsed[word] = true
    }
}
