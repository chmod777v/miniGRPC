package httpServer

import (
	"encoding/json"
	"fmt"
	"grpc/pkg/database"
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

func getRequest(body *io.ReadCloser, person interface{}, w http.ResponseWriter) error {
	if err := json.NewDecoder(*body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return err
	}
	return nil
}
func (h *httpServer) handler(w http.ResponseWriter, r *http.Request) {
	var personID personID
	var person database.Person
	switch r.Method {
	case http.MethodGet:
		if err := getRequest(&r.Body, &personID, w); err != nil {
			return
		}
		slog.Debug("RequestGet", "Data", personID)
		resp, err := h.grpcServer.GetRequestGrpc(r.Context(), &g_serv.GetRequest{Id: int64(personID.Id)})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("Error when GetRequestGrpc", "ERROR", err.Error())
			return
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
			return
		}
		w.Write(reqbyte)
	case http.MethodPost:
		if err := getRequest(&r.Body, &person, w); err != nil {
			return
		}
		slog.Debug("RequestPost", "Data", person)
		req := &g_serv.PostRequest{
			Info: &g_serv.UserInfo{
				UserId: int64(person.User_id),
				Name:   person.Name,
				Admin:  person.Admin,
			},
		}
		resp, _ := h.grpcServer.PostRequestGrpc(r.Context(), req)

		slog.Debug("ResponsePost", "Data", resp)
		respbyte, err := json.MarshalIndent(resp, "", " ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Error("Error when Marshal", "ERROR", err.Error())
			return
		}
		w.Write(respbyte)
	}
}

func Run(serv *grpcconect.Server, host string, port int) {

	httpServer := &httpServer{grpcServer: serv}
	http.HandleFunc("/", httpServer.handler)

	link := fmt.Sprintf("%s:%v", host, port)
	slog.Info("Server listening", "Host", link)
	if err := http.ListenAndServe(link, nil); err != nil {
		slog.Error("Error starting server", "ERROR", err.Error())
	}
}
