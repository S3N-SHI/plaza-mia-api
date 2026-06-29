package storage

import (
	"plaza-mia-api/internal/models"
	"time"
)

type CanchaRepository interface {
	ListarCanchas() []models.Cancha
	BuscarCanchaPorID(id uint) (models.Cancha, bool)
	CrearCancha(c models.Cancha) models.Cancha
	ActualizarCancha(id uint, datos models.Cancha) (models.Cancha, bool)
	BorrarCancha(id uint) bool
	ListarCanchasDisponibles(fecha time.Time, horaInicio, horaFin time.Time) []models.Cancha
}

type ReservaRepository interface {
	ListarReservas() []models.Reserva
	BuscarReservaPorID(id uint) (models.Reserva, bool)
	CrearReserva(r models.Reserva) models.Reserva
	CancelarReserva(id uint) bool
	ListarReservasPorUsuario(usuarioID uint) []models.Reserva
	ExisteConflicto(canchaID uint, fecha time.Time, horaInicio, horaFin time.Time, excluirID uint) bool
}

type FidelizacionRepository interface {
	BuscarJugadorPorUsuarioID(usuarioID uint) (models.Jugador, bool)
	CrearJugador(j models.Jugador) models.Jugador
	ActualizarJugador(id uint, datos models.Jugador) (models.Jugador, bool)
	ListarRecompensas() []models.Recompensa
	BuscarRecompensaPorID(id uint) (models.Recompensa, bool)
	RegistrarCanje(c models.CanjeRecompensa) (models.CanjeRecompensa, error)
	ListarCanjesPorJugador(jugadorID uint) []models.CanjeRecompensa
	ContarReservasUltimos30Dias(usuarioID uint) int
}

type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
	BuscarUsuarioPorID(id uint) (models.Usuario, bool)
}

type Almacen interface {
	CanchaRepository
	ReservaRepository
	FidelizacionRepository
}
