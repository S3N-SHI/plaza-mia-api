package service

import (
	"strings"
	"time"

	"plaza-mia-api/internal/models"
	"plaza-mia-api/internal/storage"
)

type CanchaService struct {
	repo storage.CanchaRepository
}

func NuevoCanchaService(repo storage.CanchaRepository) *CanchaService {
	return &CanchaService{repo: repo}
}

func (s *CanchaService) Listar() []models.Cancha {
	return s.repo.ListarCanchas()
}

func (s *CanchaService) Obtener(id uint) (models.Cancha, error) {
	c, ok := s.repo.BuscarCanchaPorID(id)
	if !ok {
		return models.Cancha{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *CanchaService) Crear(c models.Cancha) (models.Cancha, error) {
	if err := validarCancha(c); err != nil {
		return models.Cancha{}, err
	}
	return s.repo.CrearCancha(c), nil
}

func (s *CanchaService) Actualizar(id uint, datos models.Cancha) (models.Cancha, error) {
	if err := validarCancha(datos); err != nil {
		return models.Cancha{}, err
	}
	actualizada, ok := s.repo.ActualizarCancha(id, datos)
	if !ok {
		return models.Cancha{}, ErrNoEncontrado
	}
	return actualizada, nil
}

func (s *CanchaService) Borrar(id uint) error {
	if !s.repo.BorrarCancha(id) {
		return ErrNoEncontrado
	}
	return nil
}

func (s *CanchaService) BuscarDisponibles(fecha time.Time, horaInicio, horaFin time.Time) ([]models.Cancha, error) {
	if horaInicio.After(horaFin) || horaInicio.Equal(horaFin) {
		return nil, ErrFechaInvalida
	}
	return s.repo.ListarCanchasDisponibles(fecha, horaInicio, horaFin), nil
}

func validarCancha(c models.Cancha) error {
	if strings.TrimSpace(c.Nombre) == "" {
		return ErrCampoRequerido
	}
	if c.PrecioHora < 0 {
		return ErrPrecioNegativo
	}
	if c.Deporte != models.DeportePadel && c.Deporte != models.DeportePingPong {
		return ErrCampoRequerido
	}
	return nil
}
