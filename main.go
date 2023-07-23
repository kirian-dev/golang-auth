package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var storage = make(map[string]string)

func main() {
	http.HandleFunc("/", createHtmlForm)
	http.HandleFunc("/register", register)

	http.ListenAndServe(":8080", nil)
}

func createHtmlForm(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en>
		<head>
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Register form</title>
		</head>

		<body style="width: 100%;">
			<h1>Register form</h1>
			<form method="post" action="/register">
				<input type="text" name="username" placeholder="Please enter your username" required />
				<input type="password" name="password" placeholder="Please enter your password" required />
				<button type="submit" name="button">Submit</button>	
			</form>
		</body>
	</html>
	`
	io.WriteString(w, html)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusAccepted)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	hashPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}

	storage["username"] = username
	storage["password"] = string(hashPassword)
	log.Printf("storage: %v", storage)
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error generate password: %w", err)
	}
	return hashedPassword, nil
}
