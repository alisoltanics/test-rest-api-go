package main

import (
  "fmt"
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Book struct {
  gorm.Model
  Isbn   string  json:"isbn"
  Title  string  json:"title"
}


var db *gorm.DB
var err error

func initialMigration() {
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()

  db.AutoMigrate(&Book{})
}

func getBooks(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()

  var mybooks []Book
  db.Find(&mybooks)
  json.NewEncoder(w).Encode(mybooks)
}

func getBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()
  params := mux.Vars(r)
  var book Book
  db.Where("id = ?", params["id"]).Find(&book)
  json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()

    var book Book
  _ = json.NewDecoder(r.Body).Decode(&book)
  db.Create(&book)
  json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()
  params := mux.Vars(r)
  var book Book
  var bodyData Book
  _ = json.NewDecoder(r.Body).Decode(&bodyData)
  db.Where("id = ?", params["id"]).Find(&book)
  book.Isbn = bodyData.Isbn
  book.Title = bodyData.Title
  db.Save(&book)
  json.NewEncoder(w).Encode(updateBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
  db, err = gorm.Open("sqlite3", "test_database.db")
  if err != nil {
    fmt.Println(err.Error())
    panic("Failed to connect to database")
  }
  defer db.Close()
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  var book Book
  db.Where("id = ?", params["id"]).Find(&book)
  db.Delete(&book)
  fmt.Fprintf(w, "User deleted")
}

func handleRequests() {
  r := mux.NewRouter().StrictSlash(true)

  r.HandleFunc("/api/books", getBooks).Methods("GET")
  r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
  r.HandleFunc("/api/books", createBook).Methods("POST")
  r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
  r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

  log.Fatal(http.ListenAndServe(":8000", r))
}

func main() {
  initialMigration()
  handleRequests()
}