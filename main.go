package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	)
	"github.com/dgrijalva/jwt-go/request"
	//RSA KEYS AND INITIALISATION
	const (
	privKeyPath = "app.rsa"
	pubKeyPath = "app.rsa.pub"
	)
	var VerifyKey, SignKey []byte
	
	func initKeys() {
	var err error
	
	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	
	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
	}
	
	//STRUCT DEFINITIONS
	
	type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	}
	
	type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	}
	
	type Response struct {
	Data string `json:"data"`
	}
	
	type Token struct {
	Token string `json:"token"`
	}
	
	//SERVER ENTRY POINT
	
	func StartServer() {
	
	//PUBLIC ENDPOINTS
	http.HandleFunc("/login", LoginHandler)
	
	//PROTECTED ENDPOINTS
	http.Handle("/resource/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))
	
	log.Println("Now listening...")
	http.ListenAndServe(":8000", nil)
	
	}
	
	func main() {
	
	initKeys()
	StartServer()
	
	}
	
	//////////////////////////////////////////
	
	/////////////ENDPOINT HANDLERS////////////
	
	/////////////////////////////////////////
	
	func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
	
	}
	
	func LoginHandler(w http.ResponseWriter, r *http.Request) {
	
	var user UserCredentials
	
	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}
	
	fmt.Println(user.Username, user.Password)
	
	//validate user credentials
	if strings.ToLower(user.Username) != "alexcons" {
		if user.Password != "kappa123" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}
	
	//create a rsa 256 signer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Minute * 20).Unix(),
		"CustomUserInfo": struct {
			Name string
			Role string
		}{user.Username, "Member"},
	})

tokenString, err := token.SignedString(SignKey)
	fmt.Println(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
		
	}
	
	//create a token instance using the token string
	response := Token{tokenString}
	fmt.Println(tokenString)
	JsonResponse(response, w)
	
	}
	
	//AUTH TOKEN VALIDATION
	
	func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	
	//validate token
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})
	
	if err == nil {
	
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}
	
	}
	
	//HELPER FUNCTIONS
	
	func JsonResponse(response interface{}, w http.ResponseWriter) {
	
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	
	}