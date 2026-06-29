package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/glebarez/sqlite"
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"

    "plaza-mia-api/internal/handlers"
    "plaza-mia-api/internal/middleware"
    "plaza-mia-api/internal/models"
    "plaza-mia-api/internal/service"
    "plaza-mia-api/internal/storage"
)

func main() {
    // 1. Abrir BD
    db, err := gorm.Open(sqlite.Open("plaza_mia.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Error abriendo BD:", err)
    }

    // 2. AutoMigrate (crea las tablas)
    db.AutoMigrate(
        &models.Usuario{},
        &models.Cancha{},
        &models.Reserva{},
        &models.Jugador{},
        &models.Recompensa{},
        &models.CanjeRecompensa{},
    )

    // 3. Construir capas
    almacen := storage.NuevoAlmacenSQLite(db)
    almacen.SembrarSiVacio()

    canchasSvc      := service.NuevoCanchaService(almacen)
    reservasSvc     := service.NuevoReservaService(almacen)
    fidelizacionSvc := service.NuevoFidelizacionService(almacen, almacen)
    authSvc         := service.NuevoAuthService(almacen)

    srv := handlers.NewServer(canchasSvc, reservasSvc, fidelizacionSvc, authSvc)

    // 4. Router
    r := chi.NewRouter()
    r.Use(middleware.CORS)

    r.Post("/api/v1/auth/registro", srv.Registrar)
    r.Post("/api/v1/auth/login",    srv.Login)

    r.Group(func(r chi.Router) {
        r.Use(middleware.Auth(authSvc))
        r.Get("/api/v1/canchas",                  srv.ListarCanchas)
        r.Get("/api/v1/canchas/disponibilidad",   srv.CanchasDisponibles)
        r.Get("/api/v1/canchas/{id}",             srv.ObtenerCancha)
        r.Post("/api/v1/canchas",                 srv.CrearCancha)
        r.Put("/api/v1/canchas/{id}",             srv.ActualizarCancha)
        r.Delete("/api/v1/canchas/{id}",          srv.BorrarCancha)
        r.Get("/api/v1/reservas",                 srv.ListarReservas)
        r.Get("/api/v1/reservas/{id}",            srv.ObtenerReserva)
        r.Post("/api/v1/reservas",                srv.CrearReserva)
        r.Delete("/api/v1/reservas/{id}",         srv.CancelarReserva)
        r.Get("/api/v1/reservas/mis-reservas",    srv.MisReservas)
        r.Get("/api/v1/fidelizacion/perfil",      srv.MiPerfil)
        r.Get("/api/v1/fidelizacion/recompensas", srv.ListarRecompensas)
        r.Post("/api/v1/fidelizacion/canjear",    srv.CanjearRecompensa)
        r.Get("/api/v1/fidelizacion/canjes",      srv.HistorialCanjes)
    })

    fmt.Println("Servidor en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}