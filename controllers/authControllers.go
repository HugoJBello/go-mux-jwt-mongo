package controllers

import (
	"encoding/json"
	"fmt"
	"go-mux-jwt-mongo/models"
	u "go-mux-jwt-mongo/utils"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		response := map[string]interface{}{"status": false, "message": "Invalid request"}
		u.Respond(w, response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(*account)
	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{"status": false, "message": "Invalid request"}
		u.Respond(w, response)
		return
	}
	fmt.Println(*account)

	resp, code := models.Login(account.Username, account.Password)
	if code == 401 {
		w.WriteHeader(http.StatusUnauthorized)
	} else if code == 500 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	u.Respond(w, resp)
}
