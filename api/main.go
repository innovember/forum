package main

import (
	"crypto/tls"
	config "github.com/innovember/forum/api/config"
	db "github.com/innovember/forum/api/db"
	"github.com/innovember/forum/api/middleware"
	"github.com/innovember/forum/api/services/loadEnv"
	session "github.com/innovember/forum/api/services/session"
	"time"

	userHandler "github.com/innovember/forum/api/user/delivery"
	userRepo "github.com/innovember/forum/api/user/repository"
	userUsecase "github.com/innovember/forum/api/user/usecases"

	postHandler "github.com/innovember/forum/api/post/delivery"
	postRepo "github.com/innovember/forum/api/post/repository"
	postUsecase "github.com/innovember/forum/api/post/usecases"
	"log"
	"net/http"
	"os"
)

func Run() {
	log.Println("Server is starting...")
	if errEnv := loadEnv.Load(); errEnv != nil {
		log.Fatal(errEnv)
	}
	dbConn, err := db.GetDBInstance()
	if err != nil {
		log.Fatal("DB conn", err)
	}
	if err = db.CheckDB(dbConn, config.DBPath+"/"+config.DBSchema); err != nil {
		log.Fatal("DB schema", err)
	}
	if err = session.ResetAll(dbConn); err != nil {
		log.Fatal("Session reset", err)
	}
	session.Init(dbConn)
	// User repositories
	userRepository := userRepo.NewUserDBRepository(dbConn)
	adminRepository := userRepo.NewAdminDBRepository(dbConn)
	moderatorRepository := userRepo.NewModeratorDBRepository(dbConn)

	// Post repositories
	postRepository := postRepo.NewPostDBRepository(dbConn)
	postRateRepository := postRepo.NewRateDBRepository(dbConn)
	categoryRepository := postRepo.NewCategoryDBRepository(dbConn)
	commentRepository := postRepo.NewCommentDBRepository(dbConn)
	notificationRepository := postRepo.NewNotificationDBRepository(dbConn)
	commentRateRepository := postRepo.NewRateCommentDBRepository(dbConn)

	// User usecases
	userUcase := userUsecase.NewUserUsecase(userRepository)
	adminUcase := userUsecase.NewAdminUsecase(adminRepository)
	moderatorUcase := userUsecase.NewModeratorUsecase(moderatorRepository)
	// Post usecases
	postUcase := postUsecase.NewPostUsecase(postRepository)
	postRateUcase := postUsecase.NewRateUsecase(postRateRepository)
	categoryUcase := postUsecase.NewCategoryUsecase(categoryRepository)
	commentUcase := postUsecase.NewCommentUsecase(commentRepository)
	notificationUcase := postUsecase.NewNotificationUsecase(notificationRepository)
	commentRateUcase := postUsecase.NewRateCommentUsecase(commentRateRepository)

	//Middleware
	mux := http.NewServeMux()
	mw := middleware.NewMiddlewareManager()
	// User delivery
	userHandler := userHandler.NewUserHandler(userUcase, adminUcase,
		moderatorUcase,
		postUcase, postRateUcase,
		categoryUcase, commentUcase,
		notificationUcase, commentRateUcase)
	userHandler.Configure(mux, mw)

	// Post delivery
	postHandler := postHandler.NewPostHandler(postUcase, userUcase,
		postRateUcase, categoryUcase,
		commentUcase, notificationUcase,
		commentRateUcase)
	postHandler.Configure(mux, mw)

	port := config.APIPortDev
	if port == "" {
		port = getPort()
	}
	cer, err := tls.LoadX509KeyPair("ssl/localhost.pem", "ssl/localhost-key.pem")
	if err != nil {
		log.Fatal("SSL", err)
		return
	}
	server := &http.Server{
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    30 * time.Second,
		Addr:           config.APIURLDev,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{cer},
			InsecureSkipVerify: true,
		},
		Handler:      middleware.Limit(mux),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// log.Println("Server is listening", "http://"+config.APIURLDev)
	// err = http.ListenAndServe(":"+port, mux)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
	log.Println("Server is listening", "https://"+config.APIURLDev)
	err = server.ListenAndServeTLS("", "")
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
