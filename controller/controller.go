package main

import (
	"D/Works/GO/31/models"
	"D/Works/GO/31/mongodb"
	"D/Works/GO/31/repository"
	"D/Works/GO/31/usecase"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

var Storage usecase.Storage

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	host := "localhost"
	port := ":27017"
	database := "clients"
	mongoDBClient, err := mongodb.NewClien(context.Background(), host, port, database)
	if err != nil {
		panic(err)
	}
	Storage = repository.NewStorage(mongoDBClient, "users")
	router := chi.NewRouter()
	router.Post("/users", CreateNewUser)
	router.Get("/users/{userId}", GetUser)
	router.Get("/users", GetAllUsers)
	router.Get("/users/{userId}/friends", GetUserFriends)
	router.Put("/users/{userId}/friends", MakeFriends)
	router.Patch("/users/{userId}", UpdateUser)
	router.Delete("/users/{userId}", DeleteUser)
	serverPort := flag.String("port", ":8080", "port number")
	flag.Parse()
	server := &http.Server{
		Addr:    *serverPort,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server don't started %s\n", err)
			os.Exit(1)
		}
	}()
	fmt.Printf("Server started on port:%s\n", server.Addr)

	<-done
	fmt.Println("Stop signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error %s\n", err)
	}
	fmt.Println("Server stoped")
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var modelUser models.User

	err := json.NewDecoder(r.Body).Decode(&modelUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if modelUser.Name == "" {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}
	if modelUser.Age == 0 {
		http.Error(w, "age must be greater than 0", http.StatusBadRequest)
		return
	}

	err = Storage.CreateNewUser(context.Background(), modelUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User was created")
}

func MakeFriends(w http.ResponseWriter, r *http.Request) {
	type FriendRequest struct {
		User1 string `json:"userId1"`
		User2 string `json:"userId2"`
	}

	var friendRequest FriendRequest
	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err1, err2 := Storage.MakeFriends(context.Background(), friendRequest.User1, friendRequest.User2)
	if err1 != nil {
		http.Error(w, "User1 not found", http.StatusNotFound)
		return
	} else if err2 != nil {
		http.Error(w, "User2 not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s and %s are friends now", friendRequest.User1, friendRequest.User2)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	user, err := Storage.GetUser(context.Background(), userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := Storage.GetAllUsers(context.Background())
	if err != nil {
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}
	for _, v := range users {
		json.NewEncoder(w).Encode(v)
	}
}

func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	userFriends, err := Storage.GetUserFriends(context.Background(), userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(userFriends)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	var updateUser models.User

	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = Storage.UpdateUser(context.Background(), userId, updateUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s updated", userId)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	err := Storage.DeleteUser(context.Background(), userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s deleted", userId)
}
