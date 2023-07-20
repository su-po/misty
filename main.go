package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Raindrop struct {
	ID     int32  `json:"_id"`
	Format string `json:"type"`
	Title  string `json:"title"`
	URL    string `json:"link"`
}

type RaindropCollection struct {
	Result bool       `json:"result"`
	Items  []Raindrop `json:"items"`
}

func main() {

	cursor := 0
	var items []Raindrop
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	for {
		url := "https://api.raindrop.io/rest/v1/raindrops/0?search=type%3Aarticle&perpage=50&page=" + fmt.Sprintf("%d", cursor)
		req, _ := http.NewRequest("GET", url, nil)
		headerValue := "Bearer " + string(os.Getenv("APITOKEN"))
		req.Header.Add("Authorization", headerValue)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Something went wrong\n." + err.Error())

		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Failed to read the response body")
		}

		var raindropCollection RaindropCollection

		err = json.Unmarshal(body, &raindropCollection)
		if err != nil {
			log.Fatal("Program failed.", err)
		}

		if len(raindropCollection.Items) < 1 {
			break
		}
		items = append(items, raindropCollection.Items...)
		cursor += 1
	}
	// j, _ := json.MarshalIndent(items, "", " ")
	// fmt.Println(string(j))
}
