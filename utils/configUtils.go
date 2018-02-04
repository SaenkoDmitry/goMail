package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"log"
	"go.uber.org/zap"
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
	b, err := ioutil.ReadFile("config/mysql.yml")
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
	b, err := ioutil.ReadFile("config/tarantool.yml")
	if err != nil {
		Logger.Info("Cannot open file tarantool.yml",
			zap.Error(err),
		)
	}
	data := string(b)
	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
		Logger.Info("unmarshal yaml error",
			zap.Error(err),
		)
	}
	return t
}