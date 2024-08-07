package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type User struct {
	// Fields have to be capital
	Name     string `json:"name"`
	Language string `json:"language"`
	Id       string `json:"id"`
	Bio      string `json:"bio"`
	Version  string `json:"version"`
}

func main() {
	// Open JSON File
	file, err := os.Open("sample1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Store JSON to array of objects
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	var Users []User
	err = json.Unmarshal(byteValue, &Users)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "wassap")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			result, err := json.Marshal(Users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(result)
			return
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			query := r.URL.Query()
			id := query.Get("id")
			if id == "" {
				http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
				return
			}
			for _, user := range Users {
				if user.Id == id {
					result, err := json.Marshal(user)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Write(result)
					return
				}
			}

			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if r.Method == "POST" {
			// use decode for reading data from io.reader stream (ex. body of a http req), unmarshal for data already in memory.
			var newUser User
			if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			Users = append(Users, newUser)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "User added successfully:\n")
			fmt.Fprintf(w, "Name: %s\n", newUser.Name)
			fmt.Fprintf(w, "Language: %s\n", newUser.Language)
			fmt.Fprintf(w, "Id: %s\n", newUser.Id)
			fmt.Fprintf(w, "Bio: %s\n", newUser.Bio)
			fmt.Fprintf(w, "Version: %s\n", newUser.Version)
		} else if r.Method == "DELETE" {
			query := r.URL.Query()
			id := query.Get("id")
			if id == "" {
				http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
				return
			}
			for index, user := range Users {
				if user.Id == id {
					Users = append(Users[:index-1], Users[index+1:]...)
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintf(w, "User deleted successfully:\n")
					fmt.Fprintf(w, "Name: %s\n", user.Name)
					fmt.Fprintf(w, "Language: %s\n", user.Language)
					fmt.Fprintf(w, "Id: %s\n", user.Id)
					fmt.Fprintf(w, "Bio: %s\n", user.Bio)
					fmt.Fprintf(w, "Version: %s\n", user.Version)
					return
				}
			}
		} else if r.Method == "PUT" {
			var editedUser User
			if err := json.NewDecoder(r.Body).Decode(&editedUser); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			query := r.URL.Query()
			id := query.Get("id")
			if id == "" {
				http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
				return
			}
			for index := range Users {
				if Users[index].Id == id {
					Users[index] = editedUser
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintf(w, "User edited successfully:\n")
					fmt.Fprintf(w, "Name: %s\n", Users[index].Name)
					fmt.Fprintf(w, "Language: %s\n", Users[index].Language)
					fmt.Fprintf(w, "Id: %s\n", Users[index].Id)
					fmt.Fprintf(w, "Bio: %s\n", Users[index].Bio)
					fmt.Fprintf(w, "Version: %s\n", Users[index].Version)
					return
				}
			}
		}
	})

	fmt.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
