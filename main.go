package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", createHtmlForm)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	log.Printf("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
