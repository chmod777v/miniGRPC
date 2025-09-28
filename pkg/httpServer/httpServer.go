package httpServer

import (
	"encoding/json"
	"fmt"
	"grpc/internal/metrics"
	"grpc/pkg/database"
	grpcconect "grpc/pkg/httpServer/grpcConect"
	g_serv "grpc/pkg/proto"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type personID struct {
	Id int
}
type HttpServer struct {
	requestGrpc grpcconect.RequestGrpc
}

func getRequest(body *io.ReadCloser, person interface{}, w http.ResponseWriter) error {
	if err := json.NewDecoder(*body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return err
	}
	return nil
}
func (h *HttpServer) handleGet(w http.ResponseWriter, r *http.Request) int {
	var personID personID
	if err := getRequest(&r.Body, &personID, w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return http.StatusBadRequest
	}
	slog.Debug("RequestGet", "Data", personID)

	resp, err := h.requestGrpc.GetRequestGrpc(r.Context(), &g_serv.GetRequest{Id: int64(personID.Id)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error when GetRequestGrpc", "ERROR", err.Error())
		return http.StatusInternalServerError
	}
	slog.Debug("ResponseGet", "Data", resp)

	req := database.Person{
		Id:      personID.Id,
		User_id: int(resp.Info.UserId),
		Name:    resp.Info.Name,
		Admin:   resp.Info.Admin,
	}
	reqbyte, err := json.MarshalIndent(req, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error when Marshal", "ERROR", err.Error())
		return http.StatusInternalServerError
	}
	w.Write(reqbyte)
	return http.StatusOK
}
func (h *HttpServer) handlePost(w http.ResponseWriter, r *http.Request) int {
	var person database.Person
	if err := getRequest(&r.Body, &person, w); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return http.StatusBadRequest
	}
	slog.Debug("RequestPost", "Data", person)
	req := &g_serv.PostRequest{
		Info: &g_serv.UserInfo{
			UserId: int64(person.User_id),
			Name:   person.Name,
			Admin:  person.Admin,
		},
	}
	resp, err := h.requestGrpc.PostRequestGrpc(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error when GetRequestGrpc", "ERROR", err.Error())
		return http.StatusInternalServerError
	}

	slog.Debug("ResponsePost", "Data", resp)
	respbyte, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error when Marshal", "ERROR", err.Error())
		return http.StatusInternalServerError
	}
	w.Write(respbyte)
	return http.StatusOK
}
func (h *HttpServer) handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var status int
	defer func() {
		metrics.ObserveRequest(r.Method, time.Since(start), status)
	}()

	switch r.Method {
	case http.MethodGet:
		status = h.handleGet(w, r)
	case http.MethodPost:
		status = h.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Debug("ERROR Method not allowed", "method", r.Method)
		status = http.StatusMethodNotAllowed
	}
}

func Run(requestGrpc grpcconect.RequestGrpc, host string, port int) {

	httpServer := &HttpServer{requestGrpc: requestGrpc}
	http.HandleFunc("/", httpServer.handler)

	go func() {
		_ = metrics.Listen(host + ":8081")
	}()

	link := fmt.Sprintf("%s:%v", host, port)
	slog.Info("Server listening", "Host", link)
	if err := http.ListenAndServe(link, nil); err != nil {
		slog.Error("Error starting server", "ERROR", err.Error())
	}
}
