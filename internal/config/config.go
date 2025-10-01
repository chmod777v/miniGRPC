package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string
	Database    Database
	Grpc_server Grpc_server
	Http_server Http_server
}
type Database struct {
	Username string
	Password string
	Host     string
	Port     int
	Db_name  string
}
type Grpc_server struct {
	Port int
	Host string
}
type Http_server struct {
	Port        int
	Host        string
	Grpc_client Grpc_client
}
type Grpc_client struct {
	Port int
	Host string
}

func LoadConfig() *Config {
	var path string
	flag.StringVar(&path, "config", "", "path") //"config" - имя флага (--config) "path" - описание для справки
	flag.Parse()
	if path == "" {
		panic("CONFIG_PATH is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) { //os.Stat- - Проверка существует ли файл, os.IsNotExist-если нет то
		panic("Config file does not exist: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failet to read config" + err.Error())
	}
	return &cfg
}
