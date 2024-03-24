package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/actors/handler"
	"github.com/tkhrk1010/go-samples/scatter-gather/fruits-trading-demo/api"
)

func main() {

	// register purchasing API
	go runPurchasingAPI()
	time.Sleep(3 * time.Second) // Give server time to start
	post()

	// getTradeSupportInformation
	tradeSupportInformationHandler := handler.NewTradeSupportInformationHandler()
	jsonResponse, err := tradeSupportInformationHandler.GetTradeSupportInformation()
	if err != nil {
		fmt.Println("Error while getting trade information:", err)
		return
	}

	fmt.Println("Aggregate Results (JSON):", string(jsonResponse))

}

func runPurchasingAPI() {
	r := gin.Default()
	r.POST("/add/purchasing/apple", api.HandleApplePurchasingPost)
	r.Run(":8080")
}

func post() {
	url := "http://localhost:8080/add/purchasing/apple" // Specify the URL to send the request to

	// Create the JSON data to send
	jsonData := []byte(`{"id":"1","count": 5,"price":100.0,"purchasedAt":"2019-01-01T00:00:00Z"}`)

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create the HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status code
	fmt.Println("Response Status Code:", resp.Status)
}

func retryPost(attempts int, delay time.Duration) {
	for i := 0; i < attempts; i++ {
		post()
		time.Sleep(delay)
	}
}
