package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"crypto/md5"
	"io"
)

var Cookies map[string]string

var tokenEncodeString string = "somethingsdfsh/sdfdso;fds'SD/d"

func CreateToken(user string, password string) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set some claims
	token.Claims = jwt.MapClaims{
		"username": user,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	//Sign and get the complete encoded token as string
	return (token.SignedString([]byte(tokenEncodeString)))
}

func ParseToken(myToken string) (jwt.Claims) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenEncodeString), nil
	})
	if err != nil {
		//handling errors
		fmt.Println(err)
	}
	return token.Claims
}

func HashPassword(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

func init() {
	Cookies = make(map[string]string)
}
