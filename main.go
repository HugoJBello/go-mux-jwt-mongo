package main

import (
	"go-mux-jwt-mongo/controllers"
	"go-mux-jwt-mongo/middlewares"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const logPath = "logs.log"

var Logger *log.Logger

func main() {

	InitLogger()

	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.Use(middlewares.MiddlewareLogger)  //attach JWT auth middleware
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Access-Control-Allow-Origin"})
	//originsOk := handlers.AllowedOrigins([]string{"http://localhost:8081"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// start server listen
	// with error handling
	log.Println("Starting server on port " + port)
	log.Println(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func InitLogger() {

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Failed to open log file")
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

}
