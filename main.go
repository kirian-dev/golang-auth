package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

var storage = make(map[string]User)

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

func createHtmlForm(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.FormValue("errormsg")
	html, err := fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en>
		<head>
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Register form</title>
		</head>

		<body style="width: 100%;">
			<h1>Register form</h1>
			<p>%s</p>
			<form method="post" action="/register">
				<input type="text" name="username" placeholder="Please enter your username" required />
				<input type="password" name="password" placeholder="Please enter your password" required />
				<button type="submit" name="button">Submit</button>	
			</form>

			<h1>Login form</h1>
			<p>%s</p>
			<form method="post" action="/login">
				<input type="text" name="username" placeholder="Please enter your username" required />
				<input type="password" name="password" placeholder="Please enter your password" required />
				<button type="submit" name="button">Submit</button>	
			</form>
		</body>
	</html>
	`, errorMsg)
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

	log.Printf("User registered")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error generate password: %w", err)
	}
	return hashedPassword, nil
}

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorMsg := url.QueryEscape("your method must be POST")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}

	password := r.FormValue("password")
	username := r.FormValue("username")

	user, ok := storage[username]
	if !ok {
		errorMsg := url.QueryEscape("username not found in storage")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}

	if !comparePassword(user.Password, password) {
		errorMsg := url.QueryEscape("password not correct")
		http.Redirect(w, r, "/?errorMsg="+errorMsg, http.StatusSeeOther)
		return
	}

	log.Printf("User is logging in with password")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createToken(token string, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.At((time.Now())),
		},
		Username: username,
	})

	return token
}
func parseToken(accessToken string, key []byte) (string, error) {

}
