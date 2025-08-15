// package main

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"sync"
// 	"time"
// 	"encoding/json"
// 	"os"
// 	"database/sql"
// 	_ "github.com/mattn/go-sqlite3"
// 	"log"
// )

// const storeFile = "urlstore.json"

// var (
// 	db *sql.DB
// 	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// )

// func initDB() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "./urlshortener.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	createTable := `
// 	CREATE TABLE IF NOT EXISTS urls (
// 		code TEXT PRIMARY KEY,
// 		long_url TEXT NOT NULL UNIQUE
// 	);
// 	`
// 	if _, err := db.Exec(createTable); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
// 	initDB()  // initialize SQLite DB and table

// 	rand.Seed(time.Now().UnixNano())

// 	http.HandleFunc("/", handleHome)
// 	http.HandleFunc("/shorten", handleShorten)
// 	http.HandleFunc("/short/", handleRedirect)

// 	fmt.Println("Starting server on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handleHome(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8") // ✅ This tells browser to render HTML
// 	fmt.Fprintln(w, `
// 		<h2>Simple Go URL Shortener</h2>
// 		<form action="/shorten" method="POST">
// 			Long URL: <input name="url" type="text" style="width:300px">
// 			<input type="submit" value="Shorten">
// 		</form>
// 	`)
// }

// func handleShorten(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	longURL := r.FormValue("url")
// 	if longURL == "" {
// 		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if URL already exists
// 	var code string
// 	err := db.QueryRow("SELECT code FROM urls WHERE long_url = ?", longURL).Scan(&code)
// 	if err == sql.ErrNoRows {
// 		// Insert new short code
// 		for {
// 			code = generateCode(6)
// 			_, err := db.Exec("INSERT INTO urls (code, long_url) VALUES (?, ?)", code, longURL)
// 			if err == nil {
// 				break // success
// 			}
// 			// If code collision, generate again
// 		}
// 	} else if err != nil {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprintf(w, "Short URL: <a href=\"/short/%s\">http://localhost:8080/short/%s</a>", code, code)
// }

// func handleRedirect(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Path[len("/short/"):]
// 	var longURL string
// 	err := db.QueryRow("SELECT long_url FROM urls WHERE code = ?", code).Scan(&longURL)
// 	if err == sql.ErrNoRows {
// 		http.NotFound(w, r)
// 		return
// 	} else if err != nil {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, longURL, http.StatusFound)
// }

// func generateCode(length int) string {
// 	b := make([]rune, length)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

// // try: https://www.golang.org/doc/effective_go





package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"net/url"
)

var (
	db      *sql.DB
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./urlshortener.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		code TEXT PRIMARY KEY,
		long_url TEXT NOT NULL UNIQUE
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/short/", handleRedirect)

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `
		<h2>Simple Go URL Shortener</h2>
		<form action="/shorten" method="POST">
			Long URL: <input name="url" type="text" style="width:300px">
			<input type="submit" value="Shorten">
		</form>
	`)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate URL
	if _, err := url.ParseRequestURI(longURL); err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	var code string
	err := db.QueryRow("SELECT code FROM urls WHERE long_url = ?", longURL).Scan(&code)
	if err == sql.ErrNoRows {
		// Insert new short code
		for {
			code = generateCode(6)
			_, err := db.Exec("INSERT INTO urls (code, long_url) VALUES (?, ?)", code, longURL)
			if err == nil {
				break // success
			}
			// If code collision, generate again
		}
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Short URL: <a href=\"/short/%s\">http://localhost:8080/short/%s</a>", code, code)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/short/"):]
	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE code = ?", code).Scan(&longURL)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}

func generateCode(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
