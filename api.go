package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID    int    `json:"ID"`
	Name  string `json:"Name"`
	Type  string `json:"Type"`
	Age   int    `json:"Age"`
	Email string `json:"Email"`
	Phone int    `json:"Phone"`
	City  string `json:"City"`
}

type Users struct {
	Users []User `json:"users"`
}

var users []User

// Function to read the users from a JSON file
func readUsersFromFile(filename string) ([]User, error) {
	var users Users
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty slice
			return users.Users, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return nil, err
	}

	return users.Users, nil
}

// Function to write the users to a JSON file
func writeUsersToFile(filename string, users []User) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create the Users struct to wrap the slice of users in the "users" key
	usersStruct := Users{Users: users}

	// Use a JSON encoder to write pretty-printed JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(usersStruct)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract the ID as a string from the URL
	params := mux.Vars(r)
	idStr := params["ID"]

	// Convert the ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Search for the user with the given ID
	for _, item := range users {
		if item.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// If no user was found with that ID, return a 404
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&User{})
}

// In the CreateUser function

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Generate a unique integer ID for the new user
	user.ID = len(users) + 1 // Directly assign the integer value

	// Append the new user to the users slice
	users = append(users, user)

	// Write the updated users list to the JSON file
	if err := writeUsersToFile("users.json", users); err != nil {
		http.Error(w, fmt.Sprintf("Error writing to file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["ID"]

	// Convert the ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Search for the user with the given ID
	for index, item := range users {
		if item.ID == id {
			// Delete the user from the slice
			users = append(users[:index], users[index+1:]...)

			// Write the updated users list to the JSON file
			if err := writeUsersToFile("users.json", users); err != nil {
				http.Error(w, fmt.Sprintf("Error writing to file: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(users)
			return
		}
	}

	// If no user was found with that ID, return a 404
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("User not found")
}

func main() {
	// Load existing users from file (if any)
	var err error
	users, err = readUsersFromFile("users.json")
	if err != nil {
		log.Fatalf("Error reading users from file: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Apend users
	users = append(users, User{ID: 8, Name: "Romario", Type: "Autor", Age: 18, Phone: 32458798, Email: "romario@gmail.com", City: "Rio de Janeiro"})
	users = append(users, User{ID: 9, Name: "Ana", Type: "Autor", Age: 68, Phone: 38956474, Email: "ana@gmail.com", City: "Berlin"})
	users = append(users, User{ID: 11, Name: "Gwen", Type: "Reader", Age: 43, Phone: 55663214, Email: "gwen@gmail.com", City: "Dublin"})
	// remove users
	// simulating removal of user id 3
	for index, user := range users {
		if user.ID == 3 {
			// Remove the user from the slice
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	// Define routes
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{ID}", GetUser).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{ID}", DeleteUser).Methods("DELETE")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))

}
