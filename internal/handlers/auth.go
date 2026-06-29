package handlers

import (
	"encoding/json"
	"net/http"

	"plaza-mia-api/internal/models"
)

type registrarRequest struct {
	Nombre   string     `json:"nombre"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Rol      models.Rol `json:"rol"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
	var req registrarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	usuario, err := s.Auth.Registrar(req.Nombre, req.Email, req.Password, req.Rol)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, usuario)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	token, err := s.Auth.Login(req.Email, req.Password)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
