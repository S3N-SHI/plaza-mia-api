package handlers

import (
	"encoding/json"
	"net/http"

	"plaza-mia-api/internal/middleware"
)

type canjearRequest struct {
	RecompensaID uint `json:"recompensa_id"`
}

func (s *Server) MiPerfil(w http.ResponseWriter, r *http.Request) {
	usuarioID := r.Context().Value(middleware.ClaveUsuarioID).(uint)
	perfil, err := s.Fidelizacion.ObtenerPerfil(usuarioID)
	if err != nil {
		perfil = s.Fidelizacion.ObtenerOCrearPerfil(usuarioID)
	}
	RespondJSON(w, http.StatusOK, perfil)
}

func (s *Server) ListarRecompensas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Fidelizacion.ListarRecompensas())
}

func (s *Server) CanjearRecompensa(w http.ResponseWriter, r *http.Request) {
	usuarioID := r.Context().Value(middleware.ClaveUsuarioID).(uint)
	var req canjearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	if req.RecompensaID == 0 {
		RespondError(w, http.StatusBadRequest, "recompensa_id es requerido")
		return
	}
	canje, err := s.Fidelizacion.CanjearRecompensa(usuarioID, req.RecompensaID)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, canje)
}

func (s *Server) HistorialCanjes(w http.ResponseWriter, r *http.Request) {
	usuarioID := r.Context().Value(middleware.ClaveUsuarioID).(uint)
	RespondJSON(w, http.StatusOK, s.Fidelizacion.HistorialCanjes(usuarioID))
}
