package models

import "time"

type EstadoReserva string

const (
	EstadoPendiente  EstadoReserva = "pendiente"
	EstadoConfirmada EstadoReserva = "confirmada"
	EstadoCancelada  EstadoReserva = "cancelada"
	EstadoCompletada EstadoReserva = "completada"
)

type Reserva struct {
	ID              uint          `json:"id"               gorm:"primaryKey"`
	CanchaID        uint          `json:"cancha_id"        gorm:"not null"`
	UsuarioID       uint          `json:"usuario_id"       gorm:"not null"`
	Fecha           time.Time     `json:"fecha"            gorm:"not null"`
	HoraInicio      time.Time     `json:"hora_inicio"      gorm:"not null"`
	HoraFin         time.Time     `json:"hora_fin"         gorm:"not null"`
	Estado          EstadoReserva `json:"estado"           gorm:"not null;default:'pendiente'"`
	PuntosOtorgados int           `json:"puntos_otorgados" gorm:"default:0"`
	CreadoEn        time.Time     `json:"creado_en"`
	Cancha          *Cancha       `json:"cancha,omitempty"  gorm:"foreignKey:CanchaID"`
	Usuario         *Usuario      `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
