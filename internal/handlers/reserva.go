package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"plaza-mia-api/internal/middleware"
	"plaza-mia-api/internal/models"
	"plaza-mia-api/internal/service"
)

type crearReservaRequest struct {
	CanchaID   uint   `json:"cancha_id"`
	Fecha      string `json:"fecha"`        // "2026-07-01"
	HoraInicio string `json:"hora_inicio"`  // "10:00"
	HoraFin    string `json:"hora_fin"`     // "11:00"
}

func (s *Server) ListarReservas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Reservas.Listar())
}

func (s *Server) MisReservas(w http.ResponseWriter, r *http.Request) {
	usuarioID := r.Context().Value(middleware.ClaveUsuarioID).(uint)
	RespondJSON(w, http.StatusOK, s.Reservas.ListarPorUsuario(usuarioID))
}

func (s *Server) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id, err := parsearID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	reserva, err := s.Reservas.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, reserva)
}

func (s *Server) CrearReserva(w http.ResponseWriter, r *http.Request) {
	usuarioID := r.Context().Value(middleware.ClaveUsuarioID).(uint)
	var req crearReservaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "fecha invalida, formato: YYYY-MM-DD")
		return
	}
	hi, err := time.Parse("15:04", req.HoraInicio)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "hora_inicio invalida, formato: HH:MM")
		return
	}
	hf, err := time.Parse("15:04", req.HoraFin)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "hora_fin invalida, formato: HH:MM")
		return
	}
	nueva := models.Reserva{
		CanchaID:   req.CanchaID,
		UsuarioID:  usuarioID,
		Fecha:      fecha,
		HoraInicio: hi,
		HoraFin:    hf,
	}
	reserva, err := s.Reservas.Crear(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	// Regla de fidelizacion: calcular y acreditar puntos
	puntos := service.CalcularPuntos(fecha, hi)
	reserva.PuntosOtorgados = puntos
	s.Fidelizacion.AgregarPuntos(usuarioID, puntos)
	RespondJSON(w, http.StatusCreated, reserva)
}

func (s *Server) CancelarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := parsearID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	if err := s.Reservas.Cancelar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"mensaje": "reserva cancelada"})
}
