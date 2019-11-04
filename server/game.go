package main

import (
    "fmt"
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

type GameUpdate struct {
    SecondsLeft int
    RecordInfo string
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
                e.Data = "Your final rank was #" + strconv.Itoa(rank+1)
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

func (g *Game) word_valid(word string) *event {
    if ok := is_word(word); !ok {
        e := &event{
            Name: "wrong_because",
            Data: fmt.Sprintf("%v not an english word", word),
        }
        return e
    }

    if ok := is_substring(word, g.Phrase); !ok {
        e := &event{
            Name: "wrong_because",
            Data: fmt.Sprintf("%v is not a subword", word),
        }
        return e
    }

    return nil
}

func (g *Game) word_pts(word string) int {
    result := 100
    for i := 0; i < len(word) - 1; i++ {
        result *= 2
    }
    return result
}

func (g *Game) is_prefix(word string) (string, string) {
    for used := range g.IsUsed {
        if strings.HasPrefix(word, used) {
            return word, used
        }
        if strings.HasPrefix(used, word) {
            return used, word
        }
    }
    return "",""
}

func (g *Game) submit_word(word string, p *Player) {
    word = strings.ToLower(word)

    if e := g.word_valid(word); e != nil {
        p.Send(e)
        return
    }

    if wrd, prx := g.is_prefix(word); wrd != "" {
        e := &event{
            Name: "wrong_because",
            Data: fmt.Sprintf("%v is a prefix of %v", prx, wrd),
        }
        p.Send(e)
        return
    }

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
            if len(word) > 0 {
                g.IsUsed[word] = true
                sb.Reset()
            }
        }
    }

    if sb.Len() > 0 {
        word := sb.String()
        g.IsUsed[word] = true
    }
}

func (g *Game) remove_player(p *Player) {
    // find & delete player
    to_del := -1
    arr := g.Records
    for i,rec := range arr {
        if rec.Player == p {
            to_del = i
            break
        }
    }
    if to_del != -1 {
        arr[to_del] = arr[len(arr)-1]
        arr = arr[:len(arr)-1]
        g.Records = arr
    }
}
