package service

import (
	"plaza-mia-api/internal/models"
	"plaza-mia-api/internal/storage"
)

type FidelizacionService struct {
	repo     storage.FidelizacionRepository
	reservaR storage.ReservaRepository
}

func NuevoFidelizacionService(repo storage.FidelizacionRepository, reservaR storage.ReservaRepository) *FidelizacionService {
	return &FidelizacionService{repo: repo, reservaR: reservaR}
}

func (s *FidelizacionService) ObtenerPerfil(usuarioID uint) (models.Jugador, error) {
	j, ok := s.repo.BuscarJugadorPorUsuarioID(usuarioID)
	if !ok {
		return models.Jugador{}, ErrNoEncontrado
	}
	return j, nil
}

func (s *FidelizacionService) ObtenerOCrearPerfil(usuarioID uint) models.Jugador {
	j, ok := s.repo.BuscarJugadorPorUsuarioID(usuarioID)
	if ok {
		return j
	}
	return s.repo.CrearJugador(models.Jugador{
		UsuarioID: usuarioID,
		Nivel:     models.NivelBronze,
	})
}

func (s *FidelizacionService) AgregarPuntos(usuarioID uint, puntos int) (models.Jugador, error) {
	j := s.ObtenerOCrearPerfil(usuarioID)
	j.PuntosTotales += puntos
	reservas30d := s.repo.ContarReservasUltimos30Dias(usuarioID)
	j.Nivel = calcularNivel(reservas30d)
	actualizado, ok := s.repo.ActualizarJugador(j.ID, j)
	if !ok {
		return models.Jugador{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *FidelizacionService) ListarRecompensas() []models.Recompensa {
	return s.repo.ListarRecompensas()
}

func (s *FidelizacionService) CanjearRecompensa(usuarioID, recompensaID uint) (models.CanjeRecompensa, error) {
	j := s.ObtenerOCrearPerfil(usuarioID)
	r, ok := s.repo.BuscarRecompensaPorID(recompensaID)
	if !ok || !r.Activa {
		return models.CanjeRecompensa{}, ErrNoEncontrado
	}
	if j.PuntosTotales < r.PuntosRequeridos {
		return models.CanjeRecompensa{}, ErrPuntosInsuficientes
	}
	j.PuntosTotales -= r.PuntosRequeridos
	s.repo.ActualizarJugador(j.ID, j)
	return s.repo.RegistrarCanje(models.CanjeRecompensa{
		JugadorID:    j.ID,
		RecompensaID: recompensaID,
	})
}

func (s *FidelizacionService) HistorialCanjes(usuarioID uint) []models.CanjeRecompensa {
	j := s.ObtenerOCrearPerfil(usuarioID)
	return s.repo.ListarCanjesPorJugador(j.ID)
}

// calcularNivel es funcion pura: facil de testear.
// <5 reservas/mes = Bronze | 5-9 = Silver | 10+ = Gold
func calcularNivel(reservas30d int) models.NivelJugador {
	switch {
	case reservas30d >= 10:
		return models.NivelGold
	case reservas30d >= 5:
		return models.NivelSilver
	default:
		return models.NivelBronze
	}
}
