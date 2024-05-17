package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"

    "github.com/gorilla/mux"
)

var notionToken = os.Getenv("NOTION_TOKEN")
var notionDatabaseID = os.Getenv("NOTION_DATABASE_ID")

func fetchDatabaseItems(w http.ResponseWriter, r *http.Request) {
    url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", notionDatabaseID)

    client := &http.Client{}
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    req.Header.Add("Authorization", "Bearer "+notionToken)
    req.Header.Add("Notion-Version", "2022-06-28")
    req.Header.Add("Content-Type", "application/json")

    res, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}

func main() {
    if notionToken == "" || notionDatabaseID == "" {
        log.Fatal("Environment variables NOTION_TOKEN and NOTION_DATABASE_ID must be set")
    }

    r := mux.NewRouter()
    r.HandleFunc("/fetch-items", fetchDatabaseItems).Methods("GET")

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
