package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"sfr-backend/database"
	"sfr-backend/models"
	"sfr-backend/response"
	"sfr-backend/user"
)

var signinKey = []byte("singinKey")
var getUserDetailsFunction = database.GetUserDetails
var createUserFunction = database.CreateUser

// generateToken Function used to generate JWT token for checkAuthentication
func generateToken(userDetails user.User) (models.ResponseObject, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	var responseObject models.ResponseObject

	claims := token.Claims.(jwt.MapClaims)

	var expiration = time.Now().Add(time.Minute * 2).Unix()
	claims["authorized"] = true
	claims["user"] = userDetails.Username
	claims["exp"] = expiration

	tokenString, err := token.SignedString(signinKey)

	if err != nil {
		log.Println("Error Occurred")
		return responseObject, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["user"] = userDetails.Username
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString(signinKey)
	if err != nil {
		return responseObject, err
	}

	responseObject.AccessToken = tokenString
	responseObject.RefreshToken = rt
	responseObject.TokenExpiryTime = fmt.Sprint(expiration)

	return responseObject, nil
}

func comparePasswords(usr user.User, userCreds user.UserDetails) bool {
	log.Println("Comparing Passwords")

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(userCreds.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(usr.Password))
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Compared Passwords")
	return true
}

func hashAndSaltPassword(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func getTokenFromTokenString(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return signinKey, nil
	})

	if err != nil {
		log.Println("Unable to retrieve token from tokenString")
		log.Println(err.Error())
		return nil
	}
	return token
}

//LoginHandler - Handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var usr user.User
	log.Println(r.Body)
	error := json.NewDecoder(r.Body).Decode(&usr)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	userDetails := getUserDetailsFunction(usr.Username)
	if !comparePasswords(usr, userDetails) {
		http.Error(w, "error", http.StatusUnauthorized)
		return
	}
	token, err := generateToken(usr)
	if err != nil {
		log.Println("Error Occurred")
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	userDto := models.UserDetailsDto{}
	userDto.Username = userDetails.Username
	userDto.Email = userDetails.Email
	userDto.Firstname = userDetails.Firstname
	userDto.Lastname = userDetails.Lastname
	userDto.AccessToken = token.AccessToken
	userDto.RefreshToken = token.RefreshToken
	userDto.TokenExpiresAt = token.TokenExpiryTime

	response.WriteResponse(w, userDto)
}

//CheckAuthentication - Checks authentication for user
func CheckAuthentication(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("DISABLE_AUTH") == "true" {
			endpoint(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Header Not Found", http.StatusUnauthorized)
			fmt.Printf("Cant find cookie :/\r\n")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return signinKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			fmt.Println("Invalid token. Unable to retrieve token claims")
			fmt.Println(err.Error())
			return
		}

		if token.Valid {
			endpoint(w, r)
		}
	})
}

//CreateUser - function to create a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user user.UserDetails

	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	user.Password = hashAndSaltPassword(user.Password)
	status := createUserFunction(user)

	if status == "error" || status == "" {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response.WriteResponse(w, user)
}

// RefreshTokenCheck - Refreshes the access token with the Refresh token once the token is expired
func RefreshTokenCheck(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Body)
	var requestObj models.RefreshToken

	err := json.NewDecoder(r.Body).Decode(&requestObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(requestObj.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signinKey, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		fmt.Println("Invalid token. Unable to retrieve token claims")
		fmt.Println(err.Error())
		return
	}

	fmt.Println(token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userDetails := claims["user"]
		var usr user.User
		usr.Username = userDetails.(string)

		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := generateToken(usr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				fmt.Println("Error while retrieving new Token")
				fmt.Println(err.Error())
				return
			}

			response.WriteResponse(w, newTokenPair)
			return

		}
		http.Error(w, "Unauthorized User Access", http.StatusUnauthorized)
		fmt.Println("Unauthorized User Access")
		return
	}
}

//Logout - Function to logout
func Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Header Not Found", http.StatusUnauthorized)
		fmt.Printf("Cant find cookie :/\r\n")
		return
	}

	token := getTokenFromTokenString(tokenString)
	if token == nil {
		http.Error(w, "Authorization token not found on the request", http.StatusUnauthorized)
		fmt.Printf("Cant find Auth token :/\r\n")
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Unix()

	if claimsCheck, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println(claimsCheck["exp"])
		var claimsExp = int64(claimsCheck["exp"].(int64))
		unixTimeUTC := time.Unix(claimsExp, 0) //gives unix time stamp in utc
		log.Println("checking token exxpiry time")
		log.Println(unixTimeUTC.Sub(time.Now()))
		log.Println("rewritten token exxpiry time")
	} else {
		http.Error(w, "Invalid JWT Token Found", http.StatusUnauthorized)
		log.Printf("Invalid JWT Token")
		return
	}

	response.WriteResponse(w, "Successfully logged out user")
}
