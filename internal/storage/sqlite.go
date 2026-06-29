package storage

import (
	"time"

	"gorm.io/gorm"

	"plaza-mia-api/internal/models"
)

type AlmacenSQLite struct {
	db *gorm.DB
}

func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

// ── CanchaRepository ──────────────────────────────────────────────────────────

func (a *AlmacenSQLite) ListarCanchas() []models.Cancha {
	var canchas []models.Cancha
	a.db.Where("activa = ?", true).Find(&canchas)
	return canchas
}

func (a *AlmacenSQLite) BuscarCanchaPorID(id uint) (models.Cancha, bool) {
	var c models.Cancha
	if err := a.db.First(&c, id).Error; err != nil {
		return models.Cancha{}, false
	}
	return c, true
}

func (a *AlmacenSQLite) CrearCancha(c models.Cancha) models.Cancha {
	a.db.Create(&c)
	return c
}

func (a *AlmacenSQLite) ActualizarCancha(id uint, datos models.Cancha) (models.Cancha, bool) {
	var existente models.Cancha
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Cancha{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarCancha(id uint) bool {
	res := a.db.Delete(&models.Cancha{}, id)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) ListarCanchasDisponibles(fecha time.Time, horaInicio, horaFin time.Time) []models.Cancha {
	var ocupadas []uint
	a.db.Model(&models.Reserva{}).
		Where("fecha = ? AND estado IN ? AND hora_inicio < ? AND hora_fin > ?",
			fecha, []string{"pendiente", "confirmada"}, horaFin, horaInicio).
		Pluck("cancha_id", &ocupadas)

	var canchas []models.Cancha
	q := a.db.Where("activa = ?", true)
	if len(ocupadas) > 0 {
		q = q.Where("id NOT IN ?", ocupadas)
	}
	q.Find(&canchas)
	return canchas
}

// ── ReservaRepository ─────────────────────────────────────────────────────────

func (a *AlmacenSQLite) ListarReservas() []models.Reserva {
	var reservas []models.Reserva
	a.db.Preload("Cancha").Find(&reservas)
	return reservas
}

func (a *AlmacenSQLite) BuscarReservaPorID(id uint) (models.Reserva, bool) {
	var r models.Reserva
	if err := a.db.Preload("Cancha").First(&r, id).Error; err != nil {
		return models.Reserva{}, false
	}
	return r, true
}

func (a *AlmacenSQLite) CrearReserva(r models.Reserva) models.Reserva {
	r.CreadoEn = time.Now()
	a.db.Create(&r)
	return r
}

func (a *AlmacenSQLite) CancelarReserva(id uint) bool {
	res := a.db.Model(&models.Reserva{}).Where("id = ?", id).
		Update("estado", models.EstadoCancelada)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) ListarReservasPorUsuario(usuarioID uint) []models.Reserva {
	var reservas []models.Reserva
	a.db.Preload("Cancha").Where("usuario_id = ?", usuarioID).Find(&reservas)
	return reservas
}

func (a *AlmacenSQLite) ExisteConflicto(canchaID uint, fecha time.Time, horaInicio, horaFin time.Time, excluirID uint) bool {
	var count int64
	q := a.db.Model(&models.Reserva{}).
		Where("cancha_id = ? AND fecha = ? AND estado IN ? AND hora_inicio < ? AND hora_fin > ?",
			canchaID, fecha, []string{"pendiente", "confirmada"}, horaFin, horaInicio)
	if excluirID > 0 {
		q = q.Where("id != ?", excluirID)
	}
	q.Count(&count)
	return count > 0
}

// ── FidelizacionRepository ────────────────────────────────────────────────────

func (a *AlmacenSQLite) BuscarJugadorPorUsuarioID(usuarioID uint) (models.Jugador, bool) {
	var j models.Jugador
	if err := a.db.Where("usuario_id = ?", usuarioID).First(&j).Error; err != nil {
		return models.Jugador{}, false
	}
	return j, true
}

func (a *AlmacenSQLite) CrearJugador(j models.Jugador) models.Jugador {
	a.db.Create(&j)
	return j
}

func (a *AlmacenSQLite) ActualizarJugador(id uint, datos models.Jugador) (models.Jugador, bool) {
	var existente models.Jugador
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Jugador{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) ListarRecompensas() []models.Recompensa {
	var recompensas []models.Recompensa
	a.db.Where("activa = ?", true).Find(&recompensas)
	return recompensas
}

func (a *AlmacenSQLite) BuscarRecompensaPorID(id uint) (models.Recompensa, bool) {
	var r models.Recompensa
	if err := a.db.First(&r, id).Error; err != nil {
		return models.Recompensa{}, false
	}
	return r, true
}

func (a *AlmacenSQLite) RegistrarCanje(c models.CanjeRecompensa) (models.CanjeRecompensa, error) {
	c.CanjeadoEn = time.Now()
	if err := a.db.Create(&c).Error; err != nil {
		return models.CanjeRecompensa{}, err
	}
	return c, nil
}

func (a *AlmacenSQLite) ListarCanjesPorJugador(jugadorID uint) []models.CanjeRecompensa {
	var canjes []models.CanjeRecompensa
	a.db.Preload("Recompensa").Where("jugador_id = ?", jugadorID).Find(&canjes)
	return canjes
}

func (a *AlmacenSQLite) ContarReservasUltimos30Dias(usuarioID uint) int {
	var count int64
	hace30dias := time.Now().AddDate(0, 0, -30)
	a.db.Model(&models.Reserva{}).
		Where("usuario_id = ? AND estado = ? AND creado_en >= ?",
			usuarioID, models.EstadoCompletada, hace30dias).
		Count(&count)
	return int(count)
}

// ── UserRepository ────────────────────────────────────────────────────────────

func (a *AlmacenSQLite) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if err := a.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (a *AlmacenSQLite) BuscarUsuarioPorID(id uint) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.First(&u, id).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

// ── Seeds ─────────────────────────────────────────────────────────────────────

func (a *AlmacenSQLite) SembrarSiVacio() {
	var n int64
	a.db.Model(&models.Cancha{}).Count(&n)
	if n > 0 {
		return
	}
	canchas := []models.Cancha{
		{ID: 1, Nombre: "Padel 1", Deporte: models.DeportePadel, PrecioHora: 15.00, Activa: true},
		{ID: 2, Nombre: "Padel 2", Deporte: models.DeportePadel, PrecioHora: 15.00, Activa: true},
		{ID: 3, Nombre: "Padel 3", Deporte: models.DeportePadel, PrecioHora: 15.00, Activa: true},
		{ID: 4, Nombre: "Ping Pong 1", Deporte: models.DeportePingPong, PrecioHora: 8.00, Activa: true},
		{ID: 5, Nombre: "Ping Pong 2", Deporte: models.DeportePingPong, PrecioHora: 8.00, Activa: true},
	}
	a.db.Create(&canchas)

	recompensas := []models.Recompensa{
		{ID: 1, Nombre: "1 hora gratis Padel", Descripcion: "Canjea 100 puntos por 1 hora en cancha de padel", PuntosRequeridos: 100, DescuentoPct: 100, Activa: true},
		{ID: 2, Nombre: "50% descuento reserva", Descripcion: "Canjea 50 puntos por 50% en tu proxima reserva", PuntosRequeridos: 50, DescuentoPct: 50, Activa: true},
		{ID: 3, Nombre: "Ping Pong gratis", Descripcion: "Canjea 40 puntos por 1 hora de ping pong", PuntosRequeridos: 40, DescuentoPct: 100, Activa: true},
	}
	a.db.Create(&recompensas)
}

var _ Almacen = (*AlmacenSQLite)(nil)
var _ UserRepository = (*AlmacenSQLite)(nil)
