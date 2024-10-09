package main

// Imports
import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
)

// Article struct
type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Articles array (fake database)
var Articles []Article

/*********** Functions*************/

// Return all articles
func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

// Return single article
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    for _, article := range Articles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
        }
    }
}

// Create new article
func createNewArticle(w http.ResponseWriter, r *http.Request) { 
    reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	Articles = append(Articles, article)
    json.NewEncoder(w).Encode(article)
}

// Home page function
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}


// Handle requests (routing)
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
    log.Fatal(http.ListenAndServe(":8081", myRouter))
}


// Main function
func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}