package models

type Deporte string

const (
	DeportePadel    Deporte = "padel"
	DeportePingPong Deporte = "ping_pong"
)

type Cancha struct {
	ID         uint    `json:"id"          gorm:"primaryKey"`
	Nombre     string  `json:"nombre"      gorm:"not null"`
	Deporte    Deporte `json:"deporte"     gorm:"not null"`
	PrecioHora float64 `json:"precio_hora" gorm:"not null"`
	Activa     bool    `json:"activa"      gorm:"not null;default:true"`
}
