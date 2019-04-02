package controllers

import (
	"fmt"
	"net/http"
	u "go-mux-jwt-mongo/utils"
	"go-mux-jwt-mongo/models"
	"encoding/json"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	fmt.Println("----")
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println("invalid request")
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

