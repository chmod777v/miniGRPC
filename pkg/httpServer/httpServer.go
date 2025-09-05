package main

import (
	"encoding/json"
	"fmt"
	my_grpc "grpc/internal/grpc"
	"io"
	"log/slog"
	"net/http"
)

type personID struct {
	Id int
}

func GetRequest(body *io.ReadCloser, person interface{}) error {
	if err := json.NewDecoder(*body).Decode(&person); err != nil {
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return err
	}
	return nil
}
func handler(w http.ResponseWriter, r *http.Request) {
	var personID personID
	var person my_grpc.Person
	if r.Method == http.MethodGet {
		if err := GetRequest(&r.Body, &personID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("Request Get", "Data", personID)

	} else if r.Method == http.MethodPost {
		if err := GetRequest(&r.Body, &person); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info("Request Post", "Data", person)

	}
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера")
	}
}
