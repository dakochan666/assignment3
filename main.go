package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const PORT = ":8080"

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type DataJSON struct {
	Status `json:"status"`
}

type Weather struct {
	WaterValue  int    `json:"water_value"`
	WaterStatus string `json:"water_status"`
	WindValue   int    `json:"wind_value"`
	WindStatus  string `json:"wind_status"`
}

func main() {
	http.HandleFunc("/", randomize)

	log.Println("Server started at port", PORT)
	http.ListenAndServe(PORT, nil)
}

func randomize(rw http.ResponseWriter, r *http.Request) {
	var randomValue DataJSON
	max := 100
	min := 1

	rand.Seed(time.Now().UnixNano())
	n := min + rand.Intn(max-min+1)
	randomValue.Status.Water = n

	n = min + rand.Intn(max-min+1)
	randomValue.Status.Wind = n

	payload, err := json.Marshal(randomValue)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("data.json", payload, 0644)
	if err != nil {
		log.Fatal(err)
	}

	byteValue, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	var dataJSON DataJSON
	err = json.Unmarshal(byteValue, &dataJSON)
	if err != nil {
		log.Fatal(err)
	}

	var waterStatus string
	if dataJSON.Status.Water <= 5 {
		waterStatus = "Aman"
	} else if dataJSON.Status.Water >= 6 && dataJSON.Status.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	var windStatus string
	if dataJSON.Status.Wind <= 6 {
		windStatus = "Aman"
	} else if dataJSON.Status.Wind >= 7 && dataJSON.Status.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	weather := Weather{
		WaterValue:  dataJSON.Status.Water,
		WaterStatus: waterStatus,
		WindValue:   dataJSON.Status.Wind,
		WindStatus:  windStatus,
	}

	htmlTemplate, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Println("[ERROR]", r.Method, r.URL.Path, "error :", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlTemplate.Execute(rw, weather)
}
