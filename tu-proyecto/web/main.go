package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type CO2Data struct {
	Days           int    `json:"days"`
	DaysRegression int    `json:"daysRegression"`
	Date           string `json:"date"`
	MilesPerDay    int    `json:"milesPerDay"`
}

type CO2Response struct {
	TotalCO2 float64 `json:"totalCO2"`
}

func calculateCO2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var request CO2Data
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	response, err := sendToVM(request)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func sendToVM(request CO2Data) (CO2Response, error) {
	conn, err := net.Dial("tcp", "IP_MAQUINA_VIRTUAL:PUERTO")
	if err != nil {
		log.Println("Error connecting to VM:", err)
		return CO2Response{}, err
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&request); err != nil {
		log.Println("Error sending request to VM:", err)
		return CO2Response{}, err
	}

	var response CO2Response
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		log.Println("Error receiving response from VM:", err)
		return CO2Response{}, err
	}

	return response, nil
}

func main() {
	http.HandleFunc("/calculate", calculateCO2Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server is running at port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
