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
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://flashtag.site/", http.StatusTemporaryRedirect)
	})
	http.ListenAndServe(":8080", r)
}
