package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"plaza-mia-api/internal/models"
	"plaza-mia-api/internal/storage"
)

var secretoJWT = []byte("plaza-mia-secreto-demo-cambiar-en-produccion")

const duracionToken = 24 * time.Hour

type Claims struct {
	UsuarioID uint       `json:"uid"`
	Rol       models.Rol `json:"rol"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Registrar(nombre, email, password string, rol models.Rol) (models.Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(strings.ToLower(email))
	if nombre == "" || email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrCampoRequerido
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}
	if rol == "" {
		rol = models.RolJugador
	}
	return s.repo.CrearUsuario(models.Usuario{
		Nombre:       nombre,
		Email:        email,
		PasswordHash: string(hash),
		Rol:          rol,
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}
	return s.generarToken(u)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		Rol:       u.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWT)
}

func (s *AuthService) ValidarToken(tokenStr string) (uint, models.Rol, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretoJWT, nil
	})
	if err != nil || !token.Valid {
		return 0, "", ErrCredencialesInvalidas
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", ErrCredencialesInvalidas
	}
	return claims.UsuarioID, claims.Rol, nil
}
