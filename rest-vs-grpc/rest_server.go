package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type HelloRequest struct {
    Name string `json:"name"`
}

type HelloResponse struct {
    Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    var req HelloRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    resp := HelloResponse{Message: "Hello, " + req.Name}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    http.HandleFunc("/hello", helloHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
