package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
)

func main() {
    reqBody, _ := json.Marshal(map[string]string{
        "name": "world",
    })

    // RESTのリクエスト
		// URLでresourceを指定し、HTTPメソッドで操作を指定するのはRESTならでは
    resp, err := http.Post("http://localhost:8080/hello", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    defer resp.Body.Close()

    var respBody map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
        log.Fatalf("failed to decode response: %v", err)
    }
    log.Printf("Greeting: %s", respBody["message"])
}
