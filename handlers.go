package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func createHtmlForm(w http.ResponseWriter, r *http.Request) {
	// Check if the request has a cookie named "jwt_token"
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		// Cookie not found, proceed with the regular HTML form
		renderHtmlForm(w, "", "")
		return
	}

	// Parse the token from the cookie value
	tokenString := cookie.Value
	claims, err := parseToken(tokenString)
	if err != nil {
		// Error parsing the token, redirect to the login page
		http.Redirect(w, r, "/?errorMsg="+url.QueryEscape("Invalid token"), http.StatusSeeOther)
		return
	}

	// The token is valid, extract the username from the claims
	username := claims.Username

	// Proceed with the logged-in user interface
	// Here, you can show a welcome message or any other content for the logged-in user
	renderHtmlForm(w, fmt.Sprintf("Welcome, %s!", username), "")
}

func renderHtmlForm(w http.ResponseWriter, registrationMsg, loginMsg string) {
	html, err := fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Register and Login</title>
		</head>

		<body style="width: 100%;">
			<h1>Register form</h1>
			<p>%s</p>
			<form method="post" action="/register">
				<input type="text" name="username" placeholder="Please enter your username" required />
				<input type="password" name="password" placeholder="Please enter your password" required />
				<button type="submit" name="button">Register</button>	
			</form>

			<h1>Login form</h1>
			<p>%s</p>
			<form method="post" action="/login">
				<input type="text" name="username" placeholder="Please enter your username" required />
				<input type="password" name="password" placeholder="Please enter your password" required />
				<button type="submit" name="button">Login</button>	
			</form>
		</body>
	</html>
	`, registrationMsg, loginMsg)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, string(html))
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorMsg := url.QueryEscape("your method must be POST")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		errorMsg := url.QueryEscape("email is required")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}
	password := r.FormValue("password")
	if password == "" {
		errorMsg := url.QueryEscape("password is required")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}
	hashPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}
	storage[username] = User{
		Username: username,
		Password: string(hashPassword),
	}

	token, err := createToken(username)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie := http.Cookie{
		Name:    "jwt_token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie)

	log.Printf("User registered")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
