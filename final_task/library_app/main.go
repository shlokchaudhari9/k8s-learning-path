package main

import (
	"os"
	"log"
	"net/http"
	"gorm.io/gorm"
	"encoding/json"
	"gorm.io/driver/mysql"
)

var MYSQL_SERVER_IP = os.Getenv("MYSQL_SERVER_IP")
var MYSQL_SERVER_PORT = os.Getenv("MYSQL_SERVER_PORT")
var MYSQL_SERVER_USER = os.Getenv("MYSQL_SERVER_USER")
var MYSQL_SERVER_PASSWORD = os.Getenv("MYSQL_SERVER_PASSWORD")

type PersistenceManager struct {
	db_connection *gorm.DB
}

func (db *PersistenceManager) handle_books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetBook(db, w, r)
	case "POST":
		SetBook(db, w, r)
	case "DELETE":
		DeleteBook(db, w, r)
	case "PUT":
		UpdateBook(db, w, r)
	default:
		log.Println("Unhandled request method found")
	}
}

func (db *PersistenceManager) handle_list_books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetBooks(db, w, r)
	default:
		log.Println("Unhandled request method found")
	}
}

type Books struct {
	BookID uint `gorm:"primaryKey" json:"book_id"`
	BookName string `json:"book_name"`
	BookCost uint `json:"book_cost"`
}

func initialize_db() *gorm.DB {
	dsn := MYSQL_SERVER_USER + ":" + MYSQL_SERVER_PASSWORD + "@tcp(" + MYSQL_SERVER_IP + ":" + MYSQL_SERVER_PORT + ")/"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_ = db.Exec("CREATE DATABASE IF NOT EXISTS library;")
	_ = db.Exec("USE library;")
  	if err != nil {
    	panic("failed to connect database")
	}
	db.AutoMigrate(&Books{})
	return db
}

func RouteMux() *http.ServeMux {
	mux := http.NewServeMux()
	
	db := &PersistenceManager{db_connection: initialize_db()}
	// Library books handler
	LibraryHandler := http.HandlerFunc(db.handle_books)
	mux.Handle("/book", LibraryHandler)
	
	ListBooksHandler := http.HandlerFunc(db.handle_list_books)
	mux.Handle("/books", ListBooksHandler)

	return mux
}

func read_request(r *http.Request) Books {
	books := Books{}
	err := json.NewDecoder(r.Body).Decode(&books)
	if err != nil {
		panic(err)
	}
	return books
}

func GetBooks(db *PersistenceManager, w http.ResponseWriter, r *http.Request) {
	var books []Books
	db.db_connection.Find(&books)
	bookJson, err := json.Marshal(books)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bookJson)
}

func GetBook(db *PersistenceManager, w http.ResponseWriter, r *http.Request) {
	books := read_request(r)
	var books_struct Books
	db.db_connection.Limit(1).Find(&books_struct, "book_name = ?", books.BookName)
	bookJson, err := json.Marshal(books_struct)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bookJson)
}

func SetBook(db *PersistenceManager, w http.ResponseWriter, r *http.Request) {
	books := read_request(r)
	var books_struct Books
	db.db_connection.Limit(1).Find(&books_struct, "book_name = ?", books.BookName)
	if books_struct.BookName == "" {
		db.db_connection.Create(&Books{BookName: books.BookName, BookCost: books.BookCost})
	}
	db.db_connection.Limit(1).Find(&books_struct, "book_name = ?", books.BookName)
	bookJson, err := json.Marshal(books_struct)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bookJson)
}

func DeleteBook(db *PersistenceManager, w http.ResponseWriter, r *http.Request) {
	books := read_request(r)
	var books_struct Books
	db.db_connection.Limit(1).Find(&books_struct, "book_name = ?", books.BookName)
	if books_struct.BookName != "" {
		db.db_connection.Delete(&books_struct, "book_name = ?", books.BookName)
	}
	bookJson, err := json.Marshal(books_struct)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bookJson)
}

func UpdateBook(db *PersistenceManager, w http.ResponseWriter, r *http.Request) {
	books := read_request(r)
	var books_struct Books
	db.db_connection.Limit(1).Find(&books_struct, "book_name = ?", books.BookName)
	if books_struct.BookName != "" {
		db.db_connection.Model(&books_struct).Update("BookCost", books.BookCost)
	}
	bookJson, err := json.Marshal(books_struct)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bookJson)
}

func main() {

	mux := RouteMux()
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}