package service

import "errors"

var (
	ErrCampoRequerido        = errors.New("campo requerido ausente o vacio")
	ErrPrecioNegativo        = errors.New("el precio no puede ser negativo")
	ErrFechaInvalida         = errors.New("la fecha o franja horaria es invalida")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrConflictoHorario      = errors.New("la cancha ya esta reservada en ese horario")
	ErrEmailEnUso            = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
	ErrPuntosInsuficientes   = errors.New("puntos insuficientes para canjear esta recompensa")
	ErrSinPermiso            = errors.New("no tienes permiso para realizar esta accion")
)
