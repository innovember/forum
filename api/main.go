package main

import (
	config "github.com/innovember/forum/api/config"
	db "github.com/innovember/forum/api/db"

	userHandler "github.com/innovember/forum/api/user/delivery"
	userRepo "github.com/innovember/forum/api/user/repository"
	userUsecase "github.com/innovember/forum/api/user/usecases"
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
	if err = db.CheckDB(dbConn, config.DBPath+"/"+config.DBSchema); err != nil {
		log.Fatal("DB schema", err)
	}
	//Repository
	userRepository := userRepo.NewUserDBRepository(dbConn)

	//Usecases
	userUcase := userUsecase.NewUserUsecase(userRepository)

	//Middleware
	mux := http.NewServeMux()

	//Delivery
	userHandler := userHandler.NewUserHandler(userUcase)
	userHandler.Configure(mux)
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
