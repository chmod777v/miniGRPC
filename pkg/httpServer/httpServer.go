package main

import (
	"encoding/json"
	my_grpc "grpc/internal/grpc"
	grpcconect "grpc/pkg/httpServer/grpcConect"
	g_serv "grpc/pkg/proto"
	"io"
	"log/slog"
	"net/http"
)

type personID struct {
	Id int
}
type httpServer struct {
	grpcServer *grpcconect.Server
}

func GetRequest(body *io.ReadCloser, person interface{}, w http.ResponseWriter) error {
	if err := json.NewDecoder(*body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return err
	}
	return nil
}
func (h *httpServer) handler(w http.ResponseWriter, r *http.Request) {
	var personID personID
	var person my_grpc.Person
	if r.Method == http.MethodGet {
		if err := GetRequest(&r.Body, &personID, w); err != nil {
			return
		}
		slog.Info("Request Get", "Data", personID)
		resp, _ := h.grpcServer.PostRequestGrpc(r.Context(), &g_serv.GetRequest{Id: int64(personID.Id)})
		slog.Info("Response Get", "Data", resp)
	} else if r.Method == http.MethodPost {
		if err := GetRequest(&r.Body, &person, w); err != nil {
			return
		}
		slog.Info("Request Post", "Data", person)

	}
}

func main() {
	serv := grpcconect.NewServer()
	defer serv.Close()

	httpServer := &httpServer{grpcServer: serv}
	http.HandleFunc("/", httpServer.handler)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		slog.Error("Error starting server", "ERROR", err.Error())
	}
	slog.Info("Server started successfully")
}
