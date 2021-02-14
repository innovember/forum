package main

import (
	config "github.com/innovember/forum/api/config"
	db "github.com/innovember/forum/api/db"
	"github.com/innovember/forum/api/middleware"

	userHandler "github.com/innovember/forum/api/user/delivery"
	userRepo "github.com/innovember/forum/api/user/repository"
	userUsecase "github.com/innovember/forum/api/user/usecases"

	postHandler "github.com/innovember/forum/api/post/delivery"
	postRepo "github.com/innovember/forum/api/post/repository"
	postUsecase "github.com/innovember/forum/api/post/usecases"

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
	postRepository := postRepo.NewPostDBRepository(dbConn)
	postRateRepository := postRepo.NewRateDBRepository(dbConn)
	categoryRepository := postRepo.NewCategoryDBRepository(dbConn)
	commentRepository := postRepo.NewCommentDBRepository(dbConn)
	notificationRepository := postRepo.NewNotificationDBRepository(dbConn)
	//Usecases
	userUcase := userUsecase.NewUserUsecase(userRepository)
	postUcase := postUsecase.NewPostUsecase(postRepository)
	postRateUcase := postUsecase.NewRateUsecase(postRateRepository)
	categoryUcase := postUsecase.NewCategoryUsecase(categoryRepository)
	commentUcase := postUsecase.NewCommentUsecase(commentRepository)
	notificationUcase := postUsecase.NewNotificationUsecase(notificationRepository)
	//Middleware
	mux := http.NewServeMux()
	mw := middleware.NewMiddlewareManager()
	//Delivery
	userHandler := userHandler.NewUserHandler(userUcase)
	userHandler.Configure(mux, mw)

	postHandler := postHandler.NewPostHandler(postUcase, userUcase,
		postRateUcase, categoryUcase,
		commentUcase, notificationUcase)
	postHandler.Configure(mux, mw)
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
