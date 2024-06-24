package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type RequestData struct {
    Data string `json:"data"`
}

type ResponseData struct {
    Message string `json:"message"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
    var requestData RequestData
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    responseData := ResponseData{Message: "Processed data: " + requestData.Data}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(responseData)
}

func main() {
    http.HandleFunc("/process", processHandler)

    log.Println("VM Server started at :8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
