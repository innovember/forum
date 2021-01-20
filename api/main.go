package main

import (
	config "github.com/innovember/forum/api/config"
	//_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func Run() {
	mux := http.NewServeMux()
	port := config.apiPortDev
	if port == "" {
		port = getPort()
	}
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("Server is listening %s", config.apiURLDev)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}
	return port
}

func main() {
	Run()
}
