package utils

import (
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
	"time"
)

type mysql struct {
	Username string
	Password string
	Protocol string
	Host     string
	Port     int
	Database string
	Encoding string
}

type tarantool struct {
	Username      string
	Password      string
	Host          string
	Port          int
	Timeout       time.Duration
	Reconnect     time.Duration
	Maxreconnects uint
}

func GetMysql() mysql {
	t := mysql{}
	b, err := ioutil.ReadFile("src/main/config/mysql.yml")
	if err != nil {
		log.Fatalf("Cannot open file mysql.yml : ", err)
	}
	data := string(b)
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return t
}

func GetTarantool() tarantool {
	t := tarantool{}
	b, err := ioutil.ReadFile("src/main/config/tarantool.yml")
	if err != nil {
		log.Fatalf("Cannot open file tarantool.yml : ", err)
	}
	data := string(b)
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return t
}