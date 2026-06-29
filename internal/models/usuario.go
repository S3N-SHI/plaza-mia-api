package models

import "time"

type Rol string

const (
	RolJugador       Rol = "jugador"
	RolAdministrador Rol = "administrador"
)

type Usuario struct {
	ID           uint      `json:"id"         gorm:"primaryKey"`
	Nombre       string    `json:"nombre"     gorm:"not null"`
	Email        string    `json:"email"      gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-"          gorm:"not null"`
	Rol          Rol       `json:"rol"        gorm:"not null;default:'jugador'"`
	CreadoEn     time.Time `json:"creado_en"`
}
