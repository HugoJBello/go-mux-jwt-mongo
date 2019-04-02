package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go-mux-jwt-mongo/app"
	"go-mux-jwt-mongo/controllers"
	u "go-mux-jwt-mongo/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.HandleFunc("/", u.LogHandler(u.MessageHandler))

	router.Use(app.JwtAuthentication) //attach JWT auth middleware



	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// start server listen
	// with error handling
	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
