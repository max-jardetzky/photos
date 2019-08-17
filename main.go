package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// A User is a user.
type User struct {
	Username string
	Password string
}

var collection *mongo.Collection

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/index.html")
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)

	/* Connecting to MongoDB */
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database("test").Collection("users")
	/*
		max := User{"max", "bruh"}

		insertResult, err := collection.InsertOne(context.TODO(), max)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)

		filter := bson.D{primitive.E{Key: "username", Value: "max"}}
		var result User
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found a single document: %+v\n", result)
	*/
	log.Fatal(http.ListenAndServe(":80", nil))
}

// Signin handles signins.
func Signin(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `User` instance
	user, err := parseToUser(w, r)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if a user with the given username exists
	filter := bson.D{primitive.E{Key: "username", Value: user.Username}}
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// If the user doesn't exist, return a 404 status
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Compare the stored hashed password and the given password
	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		// If the password is wrong, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// Signup handles signups.
func Signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `User` instance
	user, err := parseToUser(w, r)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	// Checking that the username is unique
	filter := bson.D{primitive.E{Key: "username", Value: user.Username}}
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Writing to the database
	insertResult, err := collection.InsertOne(context.TODO(), &User{user.Username, string(hashedPassword)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New user registered: ", insertResult.InsertedID)
}

func parseToUser(w http.ResponseWriter, r *http.Request) (*User, error) {
	user := &User{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
