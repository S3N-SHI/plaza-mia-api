-- Queries de referencia para Plaza Mia API
-- Estas queries documentan las operaciones principales del sistema.
-- La ejecucion real se hace via GORM en internal/storage/sqlite.go

-- ── CANCHAS ──────────────────────────────────────────────────────────────────

-- Listar canchas activas
SELECT * FROM canchas WHERE activa = 1;

-- Buscar canchas disponibles en una franja horaria (regla no-CRUD)
-- Retorna canchas que NO tienen reserva activa que se solape con la franja pedida
SELECT * FROM canchas
WHERE activa = 1
  AND id NOT IN (
    SELECT cancha_id FROM reservas
    WHERE fecha = :fecha
      AND estado IN ('pendiente', 'confirmada')
      AND hora_inicio < :hora_fin
      AND hora_fin   > :hora_inicio
  );

-- ── RESERVAS ─────────────────────────────────────────────────────────────────

-- Crear reserva
INSERT INTO reservas (cancha_id, usuario_id, fecha, hora_inicio, hora_fin, estado, puntos_otorgados, creado_en)
VALUES (:cancha_id, :usuario_id, :fecha, :hora_inicio, :hora_fin, 'pendiente', :puntos, CURRENT_TIMESTAMP);

-- Verificar conflicto de horario antes de insertar
SELECT COUNT(*) FROM reservas
WHERE cancha_id   = :cancha_id
  AND fecha       = :fecha
  AND estado IN ('pendiente', 'confirmada')
  AND hora_inicio < :hora_fin
  AND hora_fin    > :hora_inicio;

-- Cancelar reserva
UPDATE reservas SET estado = 'cancelada' WHERE id = :id;

-- Reservas de un usuario
SELECT r.*, c.nombre AS cancha_nombre, c.deporte
FROM reservas r
JOIN canchas c ON c.id = r.cancha_id
WHERE r.usuario_id = :usuario_id
ORDER BY r.fecha DESC;

-- Contar reservas completadas en los ultimos 30 dias (para calcular nivel)
SELECT COUNT(*) FROM reservas
WHERE usuario_id = :usuario_id
  AND estado     = 'completada'
  AND creado_en >= DATE('now', '-30 days');

-- ── FIDELIZACION ──────────────────────────────────────────────────────────────

-- Obtener perfil de jugador
SELECT j.*, u.nombre, u.email
FROM jugadores j
JOIN usuarios u ON u.id = j.usuario_id
WHERE j.usuario_id = :usuario_id;

-- Actualizar puntos y nivel del jugador
UPDATE jugadores
SET puntos_totales = :puntos_totales,
    nivel          = :nivel,
    ultima_reserva = CURRENT_TIMESTAMP
WHERE id = :id;

-- Listar recompensas activas
SELECT * FROM recompensas WHERE activa = 1 ORDER BY puntos_requeridos ASC;

-- Registrar canje de recompensa
INSERT INTO canje_recompensas (jugador_id, recompensa_id, reserva_id, canjeado_en)
VALUES (:jugador_id, :recompensa_id, :reserva_id, CURRENT_TIMESTAMP);

-- Historial de canjes de un jugador
SELECT cr.*, r.nombre AS recompensa_nombre, r.descuento_pct
FROM canje_recompensas cr
JOIN recompensas r ON r.id = cr.recompensa_id
WHERE cr.jugador_id = :jugador_id
ORDER BY cr.canjeado_en DESC;

-- ── AUTH ─────────────────────────────────────────────────────────────────────

-- Registrar usuario
INSERT INTO usuarios (nombre, email, password_hash, rol, creado_en)
VALUES (:nombre, :email, :password_hash, :rol, CURRENT_TIMESTAMP);

-- Buscar usuario por email (para login)
SELECT * FROM usuarios WHERE email = :email LIMIT 1;
