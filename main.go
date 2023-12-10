package main

import (
	"encoding/json"
	"fmt"
	"golang/offers"
	"net/http"
	"os"
	"time"
)

func getOffers(URL string) offers.Offers {
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	body := resp.Body
	defer body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected http GET status : %s", resp.Status)
		panic(err)
	}

	var offers offers.Offers
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&offers)
	if err != nil {
		panic(err)
	}
	return offers
}

func writeJsonToFile(offers any, fileName string) {
	filePtr, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	err = json.NewEncoder(filePtr).Encode(offers)
	if err != nil {
		panic(err)
	}
}

func main() {
	startTime := time.Now()
	URL := "https://61c3deadf1af4a0017d990e7.mockapi.io/offers/near_by?lat=1.313492&lon=103.860359&rad=20"
	offers := getOffers(URL)
	//checkInDate := "2019-12-25"
	fmt.Println("Input checkin date (YYYY-MM-DD): ")
	var checkInDate string
	_, err := fmt.Scanln(&checkInDate)
	if err != nil {
		panic(err)
	}

	filteredOffer := offers.FilterOffers(checkInDate, 5)
	writeJsonToFile(filteredOffer, "output.json")
	fmt.Println("Performance time: ", time.Since(startTime))
}
