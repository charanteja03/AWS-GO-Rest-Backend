package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sfr-backend/models"
	"sfr-backend/user"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateToken(t *testing.T) {
	username := "userName"
	rval, _ := generateToken(user.User{
		Username: username,
	})

	// access token decoding
	tokenString := rval.AccessToken
	accessTokenClaims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, accessTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return signinKey, nil
	})

	tokenString = rval.RefreshToken
	refreshTokenCalims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, refreshTokenCalims, func(token *jwt.Token) (interface{}, error) {
		return signinKey, nil
	})

	assert.Equal(t, username, accessTokenClaims["user"])
	assert.Equal(t, true, accessTokenClaims["authorized"])
	assert.Equal(t, username, refreshTokenCalims["user"])
	assert.Equal(t, float64(1), refreshTokenCalims["sub"])
}

func TestComparePasswords(t *testing.T) {
	userPassword := "password"
	userBadPassword := "anything"
	userPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.MinCost)
	testTable := []struct {
		userSavedPasswordHash string
		userPassedPassword    string
		expectedResult        bool
	}{
		{string(userPasswordHash), userPassword, true},
		{string(userPasswordHash), userBadPassword, false},
	}

	for _, testCase := range testTable {
		rval := comparePasswords(
			user.User{
				Password: testCase.userPassedPassword,
			}, user.UserDetails{
				Password: testCase.userSavedPasswordHash,
			})
		assert.Equal(t, testCase.expectedResult, rval)
	}
}

func TestLoginHandler(t *testing.T) {
	userPassword := "password"
	username := "username"
	userPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.MinCost)
	preparedUserDetails := user.UserDetails{
		Password: string(userPasswordHash),
		Username: username,
	}
	emptyPayload := strings.NewReader("")
	payloadCorrectPassword := strings.NewReader(
		fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			username, userPassword))
	payloadBadPassword := strings.NewReader(
		fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			username, username))
	getUserDetailsFunction = func(string) user.UserDetails {
		return preparedUserDetails
	}
	testTable := []struct {
		requestPayload       *strings.Reader
		expectedResponseCode int
		checkDetails         bool
	}{
		{emptyPayload, http.StatusBadRequest, false},
		{payloadBadPassword, http.StatusUnauthorized, false},
		{payloadCorrectPassword, http.StatusOK, true},
	}
	for _, testCase := range testTable {
		req, _ := http.NewRequest("POST", "/", testCase.requestPayload)
		rr := httptest.NewRecorder()
		LoginHandler(rr, req)
		assert.Equal(t, testCase.expectedResponseCode, rr.Code)
		if testCase.checkDetails {
			var rval models.UserDetailsDto
			json.Unmarshal(rr.Body.Bytes(), &rval)
			assert.Equal(t, username, rval.Username)
		}
	}
}

func TestCheckAuthentication(t *testing.T) {

}

func TestCreateUserSuccess(t *testing.T) {
	userPassword := "password"
	username := "username"
	emptyPayload := strings.NewReader("")
	userPayload := strings.NewReader(
		fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			username, userPassword))
	correctUserCreationFunction := func(userDetails user.UserDetails) string {
		return "success"
	}
	createUserFunction = correctUserCreationFunction
	testTable := []struct {
		requestPayload       *strings.Reader
		expectedResponseCode int
		checkDetails         bool
	}{
		{emptyPayload, http.StatusBadRequest, false},
		{userPayload, http.StatusOK, true},
	}
	for _, testCase := range testTable {
		req, _ := http.NewRequest("POST", "/", testCase.requestPayload)
		rr := httptest.NewRecorder()
		CreateUser(rr, req)
		assert.Equal(t, testCase.expectedResponseCode, rr.Code)
		if testCase.checkDetails {
			var rval user.UserDetails
			json.Unmarshal(rr.Body.Bytes(), &rval)
			assert.Equal(t, username, rval.Username)
		}
	}
}

func TestCreateUserFail(t *testing.T) {
	userPassword := "password"
	username := "username"
	userPayload := strings.NewReader(
		fmt.Sprintf(`{"username": "%s", "password": "%s"}`,
			username, userPassword))
	errorUserCreationFunction := func(userDetails user.UserDetails) string {
		return "error"
	}
	createUserFunction = errorUserCreationFunction
	testTable := []struct {
		requestPayload       *strings.Reader
		expectedResponseCode int
	}{
		{userPayload, http.StatusInternalServerError},
	}
	for _, testCase := range testTable {
		req, _ := http.NewRequest("POST", "/", testCase.requestPayload)
		rr := httptest.NewRecorder()
		CreateUser(rr, req)
		assert.Equal(t, testCase.expectedResponseCode, rr.Code)
	}
}

func TestRefreshTokenCheck(t *testing.T) {
	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["user"] = "username"
	rtClaims["exp"] = expirationTime
	autorizedRefreshToken, _ := refreshToken.SignedString(signinKey)
	rtClaims["sub"] = 0
	notAutorizedRefreshToken, _ := refreshToken.SignedString(signinKey)
	autorizedPayload := strings.NewReader(
		fmt.Sprintf(`{"RefreshToken": "%s"}`, autorizedRefreshToken))
	notAutorizedPayload := strings.NewReader(
		fmt.Sprintf(`{"RefreshToken": "%s"}`, notAutorizedRefreshToken))
	emptyPayload := strings.NewReader("")
	badTokenPayload := strings.NewReader(fmt.Sprintf(`{"RefreshToken": "%s"}`, "asd"))
	testTable := []struct {
		requestPayload       *strings.Reader
		expectedResponseCode int
	}{
		{emptyPayload, http.StatusBadRequest},
		{badTokenPayload, http.StatusUnauthorized},
		{notAutorizedPayload, http.StatusUnauthorized},
		{autorizedPayload, http.StatusOK},
	}
	for _, testCase := range testTable {
		req, _ := http.NewRequest("POST", "/", testCase.requestPayload)
		rr := httptest.NewRecorder()
		RefreshTokenCheck(rr, req)
		assert.Equal(t, testCase.expectedResponseCode, rr.Code)
	}
}

func TestLogout(t *testing.T) {
	username := "username"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	var expiration = time.Now().Add(time.Minute * 2).Unix()
	claims["authorized"] = true
	claims["user"] = username
	claims["exp"] = expiration

	tokenString, _ := token.SignedString(signinKey)
	badToken := "token"
	testTable := []struct {
		autorizationHeaderContent string
		expectedResponseCode      int
	}{
		{"", http.StatusUnauthorized},
		{badToken, http.StatusUnauthorized},
		{tokenString, http.StatusOK},
	}
	for _, testCase := range testTable {
		req, _ := http.NewRequest("POST", "/", nil)
		req.Header.Add("Authorization", testCase.autorizationHeaderContent)
		rr := httptest.NewRecorder()
		Logout(rr, req)
		assert.Equal(t, testCase.expectedResponseCode, rr.Code)
	}
}
