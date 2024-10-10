package main

// Imports
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "io/ioutil"
)

// Article struct to match MongoDB documents
type Article struct {
    ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name    string             `json:"name" bson:"name"`
    Content string             `json:"content" bson:"content"`
}

var client *mongo.Client

// Connect to MongoDB
func connectDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatalf("Missing MongoDB URI in environment variables")
    }

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
}

func getArticleCollection() *mongo.Collection {
    return client.Database(os.Getenv("MONGO_DB_NAME")).Collection("articles")
}

// Return all articles
func returnAllArticles(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllArticles")

    collection := getArticleCollection()
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    var articles []Article
    for cursor.Next(context.TODO()) {
        var article Article
        err := cursor.Decode(&article)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        articles = append(articles, article)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(articles)
}

// Return single article
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    collection := getArticleCollection()
    var article Article
    err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&article)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Article not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(article)
}

// Create new article
func createNewArticle(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article
    json.Unmarshal(reqBody, &article)

    article.ID = primitive.NewObjectID() // Generate new ObjectID
    collection := getArticleCollection()
    _, err := collection.InsertOne(context.TODO(), article)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(article)
}

// Update article
func updateArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    reqBody, _ := ioutil.ReadAll(r.Body)
    var updatedArticle Article
    json.Unmarshal(reqBody, &updatedArticle)

    collection := getArticleCollection()
    filter := bson.M{"_id": id}
    update := bson.M{"$set": bson.M{
        "name":    updatedArticle.Name,
        "content": updatedArticle.Content,
    }}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(updatedArticle)
}

// Delete article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    collection := getArticleCollection()
    filter := bson.M{"_id": id}

    _, err = collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Article with ID %s deleted", vars["id"])
}

// HomePage function
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

// Handle requests (routing)
func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
    myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
    myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

// Main function
func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")

    connectDB()  // Connect to MongoDB
    handleRequests()  // Start the server
}
/* OLD CODE. MONGO DB IS NOT IMPLEMENTED BELOW
type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Articles array (fake database)
//var Articles []Article                         


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

// Delete article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

// Update article
func updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updatedArticle Article
	json.Unmarshal(reqBody, &updatedArticle)

	for i, article := range Articles {
		if article.Id == id {
			Articles[i] = updatedArticle
		}
	}
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
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
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

*/