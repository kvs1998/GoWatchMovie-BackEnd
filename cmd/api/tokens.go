package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var validUser = models.User {
	ID: 12,
	Email: "Keya@something.com",
	Password: "$2a$12$wS5tsxKCp4k4cuUVyGiVtehzlk8ZurWPJh0UG0W208fZ4Ts47aSRS",
}

type Credentials struct {
	Username string `json:"Email"`
	Password string `json:"Password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request){
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.logger.Println(errors.New("decoder error"))
		app.errorJSON(w, err)
		return
	}
	hashedPassword := validUser.Password
	// to check if valid user
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer ="something"
	claims.Audiences = []string{"something.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	app.writeJSON(w,http.StatusOK, string(jwtBytes), "response")
}
