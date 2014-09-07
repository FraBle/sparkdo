package main

import (
	// standard library
	"log"
	"net/http"

	// external packages
	"github.com/gorilla/mux"
)

func main() {

	// initial configuration of the application
	if err := changeDirectoryToExecutable(); err != nil {
		log.Printf("Error changing directory: %v", err)
	}
	if err := initializeLogger(); err != nil {
		log.Printf("Error initializing the logger: %v", err)
	}
	if err := readConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	// HTTP Server
	router := mux.NewRouter().StrictSlash(true)
	// router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler).Methods("GET")
	router.HandleFunc("/spark", sparkHandler).Methods("POST")
	router.HandleFunc("/monitor", monitorHandler).Methods("POST")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../src/github.com/FraBle/sparkdo/assets/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../src/github.com/FraBle/sparkdo/assets/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("../src/github.com/FraBle/sparkdo/assets/img"))))
	http.Handle("/", router)
	log.Printf("Error starting http server: %v", http.ListenAndServe(":"+CONFIG.HttpPort, nil))
}
