package main

import (
	config "github.com/innovember/forum/api/config"
	db "github.com/innovember/forum/api/db"
	//_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func Run() {
	dbConn, err := db.GetDBInstance()
	if err != nil {
		log.Fatal("DB conn", err)
	}
	if err = db.CheckDB(dbConn, config.DBSchema); err != nil {
		log.Fatal("DB schema", err)
	}
	mux := http.NewServeMux()
	port := config.APIPortDev
	if port == "" {
		port = getPort()
	}
	log.Println("Server is listening", config.APIURLDev)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

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
