package service

import (
	"time"

	"plaza-mia-api/internal/models"
	"plaza-mia-api/internal/storage"
)

type ReservaService struct {
	repo storage.ReservaRepository
}

func NuevoReservaService(repo storage.ReservaRepository) *ReservaService {
	return &ReservaService{repo: repo}
}

func (s *ReservaService) Listar() []models.Reserva {
	return s.repo.ListarReservas()
}

func (s *ReservaService) ListarPorUsuario(usuarioID uint) []models.Reserva {
	return s.repo.ListarReservasPorUsuario(usuarioID)
}

func (s *ReservaService) Obtener(id uint) (models.Reserva, error) {
	r, ok := s.repo.BuscarReservaPorID(id)
	if !ok {
		return models.Reserva{}, ErrNoEncontrado
	}
	return r, nil
}

func (s *ReservaService) Crear(r models.Reserva) (models.Reserva, error) {
	if err := validarReserva(r); err != nil {
		return models.Reserva{}, err
	}
	if s.repo.ExisteConflicto(r.CanchaID, r.Fecha, r.HoraInicio, r.HoraFin, 0) {
		return models.Reserva{}, ErrConflictoHorario
	}
	r.Estado = models.EstadoPendiente
	return s.repo.CrearReserva(r), nil
}

func (s *ReservaService) Cancelar(id uint) error {
	r, ok := s.repo.BuscarReservaPorID(id)
	if !ok {
		return ErrNoEncontrado
	}
	if r.Estado == models.EstadoCancelada || r.Estado == models.EstadoCompletada {
		return ErrCampoRequerido
	}
	s.repo.CancelarReserva(id)
	return nil
}

// CalcularPuntos es la regla de negocio principal del modulo Reservas.
// Martes y miercoles antes de las 18h otorgan el doble de puntos (baja demanda).
func CalcularPuntos(fecha time.Time, horaInicio time.Time) int {
	base := 10
	dia := fecha.Weekday()
	esBajaDemanda := (dia == time.Tuesday || dia == time.Wednesday) && horaInicio.Hour() < 18
	if esBajaDemanda {
		return base * 2
	}
	return base
}

func validarReserva(r models.Reserva) error {
	if r.CanchaID == 0 || r.UsuarioID == 0 {
		return ErrCampoRequerido
	}
	if r.Fecha.IsZero() || r.HoraInicio.IsZero() || r.HoraFin.IsZero() {
		return ErrFechaInvalida
	}
	if !r.HoraInicio.Before(r.HoraFin) {
		return ErrFechaInvalida
	}
	return nil
}
