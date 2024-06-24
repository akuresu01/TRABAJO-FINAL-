package main

import (
    "html/template"
    "log"
    "net/http"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "fmt"
)

type RequestData struct {
    Days          int    `json:"days"`
    DaysRegression int    `json:"daysRegression"`
    Date          string `json:"date"`
    MilesPerDay   int    `json:"milesPerDay"`
}

type ResponseData struct {
    Message string `json:"message"`
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/submit", submitHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    days, err := strconv.Atoi(r.FormValue("days"))
    if err != nil {
        http.Error(w, "Invalid value for days", http.StatusBadRequest)
        return
    }

    daysRegression, err := strconv.Atoi(r.FormValue("daysRegression"))
    if err != nil {
        http.Error(w, "Invalid value for daysRegression", http.StatusBadRequest)
        return
    }

    milesPerDay, err := strconv.Atoi(r.FormValue("milesPerDay"))
    if err != nil {
        http.Error(w, "Invalid value for milesPerDay", http.StatusBadRequest)
        return
    }

    requestData := RequestData{
        Days:          days,
        DaysRegression: daysRegression,
        Date:          r.FormValue("date"),
        MilesPerDay:   milesPerDay,
    }

    jsonData, err := json.Marshal(requestData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // URL de la máquina virtual
    vmURL := "http://vm:8081/process"

    resp, err := http.Post(vmURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        errorMessage := fmt.Sprintf("Error connecting to VM: %v", err)
        http.Error(w, errorMessage, http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var responseData ResponseData
    err = json.Unmarshal(body, &responseData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}
