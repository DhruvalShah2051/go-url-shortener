# Go URL Shortener

A simple URL shortener built in Go with SQLite for persistence.  
Think of it as a mini Bit.ly that you can run locally.

---

## Features
- Shorten long URLs to short codes (e.g., `http://localhost:8080/short/abc123`)
- Redirect short codes back to the original URLs
- Stores data in a local SQLite database
- Simple HTML form for easy use
- URL validation before shortening

---

## Prerequisites
Make sure you have:
- [Go](https://go.dev/) v1.18 or higher
- SQLite3 installed

---

## Installation

1. **Clone the repository**
    ```bash
    git clone https://github.com/dhruvalshah2051/go-url-shortener.git
    cd go-url-shortener
    ```
2. **Install dependencies**
    ```bash
    go mod tidy
    ```

**Note:** The app will create `urlshortener.db` automatically if it doesn’t exist.

## Usage

1. **Run the application**
    ```bash
    go run main.go
    ```
2. **Open in your browser**
    ```
    http://localhost:8080
    ```
3. **Shorten a URL**
    - Enter a long URL in the form and click **Shorten**
    - You’ll get a short link like `http://localhost:8080/short/abc123`

---

## License
MIT License
