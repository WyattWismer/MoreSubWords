package main

import (
    "io/ioutil"
    "log"
    "strings"
    "math/rand"
    "encoding/json"
)

var word_lookup map[string]int
var phrase_lookup []string

func setup_phrase_lookup() {
    bytes, err := ioutil.ReadFile("text/phrases.txt")
    if err != nil {
        log.Fatal(err)
    }
    raw := string(bytes)
    phrase_lookup = strings.Split(raw,"\n")
}

func get_phrase() string {
    ind := rand.Intn(len(phrase_lookup))
    return phrase_lookup[ind]
}

func setup_word_lookup() {
    bytes, err := ioutil.ReadFile("text/words.json")
    if err != nil {
        log.Fatal(err)
    }
    json.Unmarshal(bytes, &word_lookup)
}

func is_word(word string) bool {
    _, ok := word_lookup[strings.ToLower(word)]
    return ok
}

func is_substring(word, phrase string) bool {
    word = strings.ToLower(word)
    phrase = strings.ToLower(phrase)

    lcs := longest_common_subsequence(word, phrase)
    return lcs == len(word)
}

func longest_common_subsequence(a, b string) int {
    n, m := len(a), len(b)
    LCS := make([][]int, n)
    for i := range(LCS) {
        LCS[i] = make([]int, m)
    }
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            if (a[i] == b[j]) {
                LCS[i][j] = 1
                if i > 0 && j > 0 {
                    LCS[i][j] += LCS[i-1][j-1]
                }
            } else {
                if i > 0 {
                    LCS[i][j] = LCS[i-1][j]
                }
                if j > 0 && LCS[i][j-1] > LCS[i][j] {
                    LCS[i][j] = LCS[i][j-1]
                }
            }
        }
    }
    return LCS[n-1][m-1];
}
