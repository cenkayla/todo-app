package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cenkayla/todo-app/models"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

//Create connection with mongo db
func init() {
	databaseConnection()
}

//mongoDB database configuration.
func databaseConnection() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error when loading .env file")
	}
	//Getting variable values from .env file.
	dbURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	dbCollName := os.Getenv("DB_COLLECTION_NAME")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(dbCollName)
	fmt.Println("Connected to mongoDB!")
}

//Getting all the task route.
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	task := getAllTasks()
	json.NewEncoder(w).Encode(task)
}

//Creates task route
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	var task models.Todo
	_ = json.NewDecoder(r.Body).Decode(&task)
	insertTask(task)
	json.NewEncoder(w).Encode(task)
}

// DeleteTask delete one task route
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	params := mux.Vars(r)
	deleteOne(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

// DeleteAllTask delete all tasks route
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	count := deleteAll()
	json.NewEncoder(w).Encode(count)
}

//Get all tasks from DB and return it.
func getAllTasks() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

//Inserting value to DB
func insertTask(task models.Todo) {
	insertResult, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted single record ", insertResult.InsertedID)
}

//Delete one task from the DB, delete by ID
func deleteOne(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted document", d.DeletedCount)
}

//Delete all the tasks from the DB
func deleteAll() int64 {
	d, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted document", d.DeletedCount)
	return d.DeletedCount
}
