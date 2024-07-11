package main

import (
	"fmt"
	"net/http"
	"time"

	// Package Management
	"github.com/dgrijalva/jwt-go"
)

// Structs
type Person struct {
	Name string
	Age  int
}

// Interfaces
type Greeter interface {
	Greet() string
}

type EnglishPerson struct {
	Name string
}

func (e EnglishPerson) Greet() string {
	return "Hello, " + e.Name
}

// Authentication with JWT
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

func getToken(w http.ResponseWriter, r *http.Request) {
	token, err := generateJWT("username")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token))
}

func complexFunction(a, b int, ch chan int) {
	go func() {
		time.Sleep(2 * time.Second)
		ch <- a + b
	}()
}

func main() {
	p := Person{Name: "Fahmi", Age: 26}
	fmt.Println(p)

	e := EnglishPerson{Name: "Prasanda"}
	fmt.Println(e.Greet())

	http.HandleFunc("/welcome", authenticate(welcome))
	http.HandleFunc("/token", getToken)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
