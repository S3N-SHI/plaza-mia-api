-- Schema de Plaza Mia API
-- Este archivo es referencia documental del esquema.
-- La creacion real de tablas la hace GORM con AutoMigrate en main.go

CREATE TABLE IF NOT EXISTS usuarios (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre        TEXT    NOT NULL,
    email         TEXT    NOT NULL UNIQUE,
    password_hash TEXT    NOT NULL,
    rol           TEXT    NOT NULL DEFAULT 'jugador',
    creado_en     DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS canchas (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre      TEXT    NOT NULL,
    deporte     TEXT    NOT NULL,
    precio_hora REAL    NOT NULL,
    activa      INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS reservas (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    cancha_id        INTEGER NOT NULL REFERENCES canchas(id),
    usuario_id       INTEGER NOT NULL REFERENCES usuarios(id),
    fecha            DATETIME NOT NULL,
    hora_inicio      DATETIME NOT NULL,
    hora_fin         DATETIME NOT NULL,
    estado           TEXT     NOT NULL DEFAULT 'pendiente',
    puntos_otorgados INTEGER  NOT NULL DEFAULT 0,
    creado_en        DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jugadores (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id     INTEGER NOT NULL UNIQUE REFERENCES usuarios(id),
    nivel          TEXT    NOT NULL DEFAULT 'bronze',
    puntos_totales INTEGER NOT NULL DEFAULT 0,
    racha_semanas  INTEGER NOT NULL DEFAULT 0,
    ultima_reserva DATETIME
);

CREATE TABLE IF NOT EXISTS recompensas (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre            TEXT    NOT NULL,
    descripcion       TEXT,
    puntos_requeridos INTEGER NOT NULL,
    descuento_pct     REAL    NOT NULL DEFAULT 0,
    activa            INTEGER NOT NULL DEFAULT 1,
    vigencia_hasta    DATETIME
);

CREATE TABLE IF NOT EXISTS canje_recompensas (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    jugador_id    INTEGER  NOT NULL REFERENCES jugadores(id),
    recompensa_id INTEGER  NOT NULL REFERENCES recompensas(id),
    reserva_id    INTEGER  REFERENCES reservas(id),
    canjeado_en   DATETIME DEFAULT CURRENT_TIMESTAMP
);
