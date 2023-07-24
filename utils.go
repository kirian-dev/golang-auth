package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

var storage = make(map[string]User)

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error generate password: %w", err)
	}
	return hashedPassword, nil
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

	token, err := createToken(username)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie := http.Cookie{
		Name:    "jwt_token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24), // Set the cookie expiration time
	}
	http.SetCookie(w, &cookie)
	log.Printf("User is logging in with password")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func comparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
