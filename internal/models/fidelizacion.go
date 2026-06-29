package models

import "time"

type NivelJugador string

const (
	NivelBronze NivelJugador = "bronze"
	NivelSilver NivelJugador = "silver"
	NivelGold   NivelJugador = "gold"
)

type Jugador struct {
	ID            uint         `json:"id"              gorm:"primaryKey"`
	UsuarioID     uint         `json:"usuario_id"      gorm:"uniqueIndex;not null"`
	Nivel         NivelJugador `json:"nivel"           gorm:"not null;default:'bronze'"`
	PuntosTotales int          `json:"puntos_totales"  gorm:"default:0"`
	RachaSemanas  int          `json:"racha_semanas"   gorm:"default:0"`
	UltimaReserva *time.Time   `json:"ultima_reserva"`
	Usuario       *Usuario     `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}

type Recompensa struct {
	ID               uint       `json:"id"                gorm:"primaryKey"`
	Nombre           string     `json:"nombre"            gorm:"not null"`
	Descripcion      string     `json:"descripcion"`
	PuntosRequeridos int        `json:"puntos_requeridos" gorm:"not null"`
	DescuentoPct     float64    `json:"descuento_pct"     gorm:"default:0"`
	Activa           bool       `json:"activa"            gorm:"not null;default:true"`
	VigenciaHasta    *time.Time `json:"vigencia_hasta"`
}

type CanjeRecompensa struct {
	ID           uint        `json:"id"            gorm:"primaryKey"`
	JugadorID    uint        `json:"jugador_id"    gorm:"not null"`
	RecompensaID uint        `json:"recompensa_id" gorm:"not null"`
	ReservaID    *uint       `json:"reserva_id"`
	CanjeadoEn   time.Time   `json:"canjeado_en"`
	Jugador      *Jugador    `json:"jugador,omitempty"    gorm:"foreignKey:JugadorID"`
	Recompensa   *Recompensa `json:"recompensa,omitempty" gorm:"foreignKey:RecompensaID"`
}
