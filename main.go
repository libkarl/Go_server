package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/urfave/cli/v2"

)


var app = &cli.App{
	Name:        "pokerai",
	Usage:       "pokerai play",
	Description: "PokerAI bot build, raw power, superhuman!",

	Commands: []*cli.Command{
		BuildCMD,
	},
}

var BuildCMD = &cli.Command{
	Name:  "build",
	Usage: "will try to train the agent",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "cfr.db", Value: ".cfr.db"},
	},
	Action: func(ctx *cli.Context) error {
		
		return nil
	},
}

type Person struct {
	gorm.Model
	Name string
	Email string `gorm:"typevarchar(100);unique_index"` // nastavení gormu, aby byl jen jeden email pro každého uživatele 
	Books []Book
}

type Book struct {
	gorm.Model 

	Title string
	Author string
	CallNumber int `gorm:"unique_index"`
	PersonID int
}

var (
	person = &Person{
		Name: "Jack", Email: "jack@gmail.com",
	}
	books = []Book{
		{Title: "The rules of Thinking", Author:"Richard Templar", CallNumber: 1234, PersonID: 1},
		{Title: "Happy world champion", Author:"Deko Montera", CallNumber: 12345, PersonID: 1},
	}
)

var db *gorm.DB
var err error

func main () {
	
	// Loading env variables it needs fist call command (source .env) in terminal
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	//DAtabase connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort )

	// Open connection to the database
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	// Close connection to database when the main function finishes
	defer db.Close()

	// Make migration to the database it passes our struct ot the database
	// so if my database will get request to create Person it will need to folow the defined struct if they have not already been created.
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	
	// API router
	router := mux.NewRouter()
	// returns all people
	router.HandleFunc("/people", getPeople).Methods("GET")
	// returns person by id
	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	// return book by id
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	// create person 
	router.HandleFunc("/create/person", createPerson).Methods("POST")
	// create book
	router.HandleFunc("/create/book", createBook).Methods("POST")
	// get all books
	router.HandleFunc("/books", getBooks).Methods("GET")
	// delete person by id
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	// delete book by id
	router.HandleFunc("/delete/book/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// API controllers 
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person 
	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}


func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var books []Book
	var person Person 
	db.First(&person, params["id"])
	db.Model(&person).Related(&books) // It will find all books for person with sended id 

	person.Books = books 
	json.NewEncoder(w).Encode(&person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person 
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	var person Person
	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}

// Books controllers

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book 

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func getBook (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(&book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book 
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	var book Book
	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}