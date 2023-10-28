package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	http.HandleFunc("/", s.handleGetUser)
	return http.ListenAndServe(listenAddr, nil)
}

func (s *ApiServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.svc.GetUser(context.Background())
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "Application/json")
	return json.NewEncoder(w).Encode(v)
}
