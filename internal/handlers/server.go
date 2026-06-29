package handlers

import "plaza-mia-api/internal/service"

type Server struct {
	Canchas      *service.CanchaService
	Reservas     *service.ReservaService
	Fidelizacion *service.FidelizacionService
	Auth         *service.AuthService
}

func NewServer(
	canchas *service.CanchaService,
	reservas *service.ReservaService,
	fidelizacion *service.FidelizacionService,
	auth *service.AuthService,
) *Server {
	return &Server{
		Canchas:      canchas,
		Reservas:     reservas,
		Fidelizacion: fidelizacion,
		Auth:         auth,
	}
}
