package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Book struct {
	id           int64  `json: "id" gorm:"primary key"`
	name         string `json: "name"`
	author       string `json: "author"`
	published_at string `json: "published_at"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("CREATE DATABASE orders_db")
	db.Exec("USE orders_db")

	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&Book{})
}

func main() {
	router := mux.NewRouter()
	// Create
	router.HandleFunc("/book", createBook).Methods("POST")
	// Read
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	// Read-all
	router.HandleFunc("/book", getBook).Methods("GET")
	// Update
	router.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	// Delete
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")
	// Initialize db connection
	initDB()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book
	json.NewDecoder(r.Body).Decode(&newBook)
	// Creates new order by inserting records in the `orders` and `items` table
	db.Create(&newBook)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBooks Book
	db.Preload("Items").Find(&newBooks)
	json.NewEncoder(w).Encode(newBooks)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderID := params["id"]

	var newBook Book
	db.Preload("Items").First(&newBook, inputOrderID)
	json.NewEncoder(w).Encode(newBook)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var updatedBook Book
	json.NewDecoder(r.Body).Decode(&updatedBook)
	db.Save(&updatedBook)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inputOrderID := params["id"]
	// Convert `orderId` string param to uint64
	id64, _ := strconv.ParseUint(inputOrderID, 10, 64)
	// Convert uint64 to uint
	idToDelete := uint(id64)

	db.Where("id = ?", idToDelete).Delete(&Book{})
	w.WriteHeader(http.StatusNoContent)
}
