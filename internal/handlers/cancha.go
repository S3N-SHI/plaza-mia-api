package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"plaza-mia-api/internal/models"
)

func (s *Server) ListarCanchas(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, s.Canchas.Listar())
}

func (s *Server) ObtenerCancha(w http.ResponseWriter, r *http.Request) {
	id, err := parsearID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un numero entero positivo")
		return
	}
	cancha, err := s.Canchas.Obtener(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, cancha)
}

func (s *Server) CrearCancha(w http.ResponseWriter, r *http.Request) {
	var nueva models.Cancha
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	creada, err := s.Canchas.Crear(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

func (s *Server) ActualizarCancha(w http.ResponseWriter, r *http.Request) {
	id, err := parsearID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	var datos models.Cancha
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido: "+err.Error())
		return
	}
	actualizada, err := s.Canchas.Actualizar(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, actualizada)
}

func (s *Server) BorrarCancha(w http.ResponseWriter, r *http.Request) {
	id, err := parsearID(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id invalido")
		return
	}
	if err := s.Canchas.Borrar(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusNoContent, nil)
}

// CanchasDisponibles atiende GET /api/v1/canchas/disponibilidad
// Ejemplo: ?fecha=2026-07-01&hora_inicio=10:00&hora_fin=11:00
func (s *Server) CanchasDisponibles(w http.ResponseWriter, r *http.Request) {
	fechaStr := r.URL.Query().Get("fecha")
	hiStr := r.URL.Query().Get("hora_inicio")
	hfStr := r.URL.Query().Get("hora_fin")
	if fechaStr == "" || hiStr == "" || hfStr == "" {
		RespondError(w, http.StatusBadRequest, "parametros requeridos: fecha, hora_inicio, hora_fin")
		return
	}
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "fecha invalida, formato: YYYY-MM-DD")
		return
	}
	hi, err := time.Parse("15:04", hiStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "hora_inicio invalida, formato: HH:MM")
		return
	}
	hf, err := time.Parse("15:04", hfStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "hora_fin invalida, formato: HH:MM")
		return
	}
	canchas, err := s.Canchas.BuscarDisponibles(fecha, hi, hf)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, canchas)
}

func parsearID(r *http.Request) (uint, error) {
	n, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	return uint(n), err
}
