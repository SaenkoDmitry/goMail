package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"crypto/md5"
	"io"
	"unicode/utf8"
	"strconv"
	"go.uber.org/zap"
	"errors"
	"sync"
	"net/http"
	"encoding/json"
)

var (
	Logger, _ = zap.NewProduction()
	Mu = &sync.Mutex{}
)

var Cookies map[string]string

func WriteResponse(w http.ResponseWriter, err error, resp interface{}, code int) {
	result := generateResult(err, resp)
	w.WriteHeader(code)
	js, _ := json.MarshalIndent(result, "", "  ")
	w.Write(js)
}

func generateResult(err error, resp interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	if err != nil {
		res["error"] = err.Error()
	} else {
		res["error"] = nil
	}
	res["response"] = resp
	return
}

func ConvertToNameInDB(name string, user_id uint64) string {
	return name + fmt.Sprintf("%v", user_id)
}

func ConvertToRealName(name string, user_id uint64) string {
	return name[:utf8.RuneCountInString(name) - utf8.RuneCountInString(strconv.FormatUint(user_id, 10))]
}

var tokenEncodeString string = "somethingsdfsh/sdfdso;fds'SD/d"

func HandleToken(s string) (string, error) {
	if utf8.RuneCountInString(s) > 8 {
		s = s[7:]
	} else {
		return "", errors.New("too short authorization string")
	}
	return s, nil
}

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
	Mu.Lock()
	Cookies = make(map[string]string)
	Mu.Unlock()
}
