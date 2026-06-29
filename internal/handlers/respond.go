package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"plaza-mia-api/internal/service"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

func statusDeError(err error) int {
	switch {
	case errors.Is(err, service.ErrNoEncontrado):
		return http.StatusNotFound
	case errors.Is(err, service.ErrEmailEnUso), errors.Is(err, service.ErrConflictoHorario):
		return http.StatusConflict
	case errors.Is(err, service.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrSinPermiso):
		return http.StatusForbidden
	case errors.Is(err, service.ErrPuntosInsuficientes):
		return http.StatusUnprocessableEntity
	case errors.Is(err, service.ErrCampoRequerido), errors.Is(err, service.ErrPrecioNegativo),
		errors.Is(err, service.ErrFechaInvalida):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
