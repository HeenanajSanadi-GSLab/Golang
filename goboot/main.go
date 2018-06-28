package main

import (
	//"flag"
	"os"
	"log"
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/mux"
	"./config"
	"./controllers"
	"github.com/google/logger"
)

//const logPath = "./golang.log"

//var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
	//flag.Parse()

  	

	fmt.Printf("hiiiiiiiii")
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	lf, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()

  	//defer logger.Init("LoggerExample", *verbose, true, lf).Close()

  	loggerOne := logger.Init("LoggerFile", false, true, lf)
  	//log.Print(typeof(loggerOne))
	defer loggerOne.Close()
	loggerOne.Info("This will log to the log file and the system log")
	logger.Info("This is the same as using loggerOne")


	r := httprouter.New()
	//r := mux.NewRouter()
	uc := controllers.NewMyUserController()
	r.GET("/user/:id", uc.GetUser)
	r.GET("/users", uc.GetUsers)
	r.POST("/user", uc.CreateUser)
	r.PUT("/user", uc.UpdateUser)
	r.DELETE("/user/:id", uc.RemoveUser)
	/*r.HandleFunc("/user/:id", uc.GetUser).Methods("GET")
	r.HandleFunc("/users", uc.GetUsers).Methods("GET")
	r.HandleFunc("/user", uc.CreateUser).Methods("POST")
	r.HandleFunc("/user", uc.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/:id", uc.RemoveUser).Methods("DELETE")*/

	server := config.Server_port
	log.Printf("Started server %s .....", server)
	http.ListenAndServe(server, r)
}
