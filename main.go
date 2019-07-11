package main
import(
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"JSMPJ_jwt/models"
	"strings"
	"github.com/dgrijalva/jwt-go"
)
const (
	privKeyPath="app.rsa"
	pubKeyPath="app.rsa.pub"
)
type UserCredentials struct{
	UserName string `json:"username"`
	Password string `json:"password"`
}
var VerifyKey, SignKey []byte
func LoginHandler(w http.ResponseWriter, r *http.Request){
	var user UserCredentials
	err:= json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w,"Error in request")
		return
	}
	fmt.Println(user.UserName, user.Password)
	if strings.ToLower(user.UserName) != "alexcons"{
		if user.Password != "kappa123"{
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error in logging")
			fmt.Fprint(w,"Invalid Credentials")
		}
	}
	signer := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := make(jwt.MapClaims)
	claims["iss"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	claims["CustomUserInfo"] = struct {
		Name string
		Role string
	}{user.UserName, "Member"}
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"Error while signin token")
		log.Printf("Error signing token : %v\n",err)
	}
	response := Token{tokenString}
	JsonResponse(response,w)

}
func initKeys(){
	var err error
	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil{
		log.Fatal("Error in reading private key")
	}
	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil{
		log.Fatal("Error in reading public key")
	}
}
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

	json, err :=  json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func handler(){
  origins:=handlers.AllowedOrigins([]string{"*"})
  methods:=handlers.AllowedMethods([]string{"GET","PUT","POST","DELETE"})
  newRouter:=mux.NewRouter()
  newRouter.HandleFunc("/",models.HomePage).Methods("GET")
  log.Fatal(http.ListenAndServe(":8000",handlers.CORS(origins,methods)(newRouter)))
}
func main(){
fmt.Println("JSMPJ Corporation Server is active now.....")
initKeys()
handler()
}