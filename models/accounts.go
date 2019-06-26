package models

import (
	"context"
	"fmt"
	u "go-mux-jwt-mongo/utils"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

/*
JWT claims struct
*/
type Token struct {
	UserId string
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	Username string `bson:"username" json:"username,omitempty"`
	Password string `bson:"password" json:"password,omitempty"`
	Token    string `bson:"token" json:"token,omitempty"`
	ID       string `bson:"_id,omitempty" json:"_id,omitempty"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if account.Username == "" {
		return u.Message(false, "Username address is required"), false
	}

	if len(account.Password) < 2 {
		return u.Message(false, "Password is required"), false
	}

	//check for errors and duplicate Usernames
	db := GetDB()
	collection := db.Collection("users")
	foundAccount := &Account{}
	err := collection.FindOne(context.Background(), bson.M{"username": account.Username}).Decode(foundAccount)

	if err != nil {
		fmt.Println(err)
	}

	if foundAccount.Username != "" {
		return u.Message(false, "Username address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}
	fmt.Println(account)
	account.ID = bson.NewObjectId().Hex()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	db := GetDB()
	collection := db.Collection("users")
	res, err := collection.InsertOne(context.Background(), account)
	fmt.Println(res)

	if err != nil {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	fmt.Println(account.ID)

	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(username, password string) (resp map[string]interface{}, code int) {
	account := &Account{}
	db := GetDB()
	collection := db.Collection("users")
	foundAccount := &Account{}
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(foundAccount)
	if err != nil {
		fmt.Println(err)
		return u.Message(false, "Invalid login credentials. Please try again"), 401
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundAccount.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again"), 401
	}
	//Worked! Logged In
	account.Password = ""
	account.Username = foundAccount.Username
	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resps := u.Message(true, "Logged In")
	resps["account"] = account
	return resps, 200
}

func GetUser(u string) *Account {

	acc := &Account{}
	db := GetDB()
	collection := db.Collection("users")
	foundAccount := Account{}
	err := collection.FindOne(context.Background(), bson.M{"_d": foundAccount.ID}).Decode(foundAccount)
	if err != nil { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
