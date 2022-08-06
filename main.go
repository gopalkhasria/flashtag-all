package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	StartFirestoreClient()
	defer Client.Close()
	r := mux.NewRouter()
	log.Println("Connected to firestore")
	log.Println("Start listening")
	fileServer := http.FileServer(http.Dir("./assets"))
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fileServer))
	r.HandleFunc("/{name}", Serve)
	http.ListenAndServe(":8080", r)
}
