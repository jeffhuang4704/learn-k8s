package main

import (
    "crypto/tls"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // Import SQLite3 driver
)

func logToSQLite(db *sql.DB, timestamp time.Time, url string, filename string) error {
    // Insert log entry into SQLite database
    query := `INSERT INTO logs (timestamp, url, filename) VALUES (?, ?, ?)`
    _, err := db.Exec(query, timestamp, url, filename)
    return err
}

func genericHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    log.Printf("Received %s request at %s", r.Method, r.URL.Path)
    
    if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read request body", http.StatusInternalServerError)
            return
        }
        
        // Generate a unique filename based on timestamp and random number
        timestamp := time.Now()
        filename := fmt.Sprintf("/output/request_%d_%d", timestamp.UnixNano(), rand.Intn(10000))
        
        contentType := r.Header.Get("Content-Type")
        if strings.HasPrefix(contentType, "text/") || strings.Contains(contentType, "json") || strings.Contains(contentType, "xml") {
            filename += ".txt"
            err = os.WriteFile(filename, body, 0644)
        } else {
            filename += ".bin"
            encodedData := base64.StdEncoding.EncodeToString(body)
            err = os.WriteFile(filename, []byte(encodedData), 0644)
        }
        
        if err != nil {
            log.Printf("Failed to write request body to file: %v", err)
            http.Error(w, "Failed to write request body to file", http.StatusInternalServerError)
            return
        }
        
        log.Printf("Request body written to %s", filename)
        
        // Log the information into SQLite
        err = logToSQLite(db, timestamp, r.URL.Path, filename)
        if err != nil {
            log.Printf("Failed to log to SQLite: %v", err)
        }
    }
    
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Received %s request at %s", r.Method, r.URL.Path)
}

func setupDatabase() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "/output/logs.db")
    if err != nil {
        return nil, err
    }
    
    // Create the table if it doesn't exist
    query := `
    CREATE TABLE IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME,
        url TEXT,
        filename TEXT
    );`
    _, err = db.Exec(query)
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

func main() {
    certFile := "cert.pem"
    keyFile := "key.pem"
    port := 8443
    
    // Set up SQLite database
    db, err := setupDatabase()
    if err != nil {
        log.Fatalf("Failed to set up database: %v", err)
    }
    defer db.Close()
    
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS10,
        InsecureSkipVerify: true,
        ClientAuth: tls.NoClientCert,
    }
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        genericHandler(w, r, db)
    })
    
    server := &http.Server{
        Addr:      fmt.Sprintf(":%d", port),
        Handler:   http.DefaultServeMux,
        TLSConfig: tlsConfig,
    }
    
    log.Printf("HTTPS server is listening on port %d...", port)
    err = server.ListenAndServeTLS(certFile, keyFile)
    if err != nil {
        log.Fatalf("Failed to start HTTPS server: %v", err)
    }
}

