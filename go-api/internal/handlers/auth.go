package handlers

import (
	"os"
	"strings"

	"go-api/internal/auth"
	"go-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

// LoginHandler maneja POST /api/login
func LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "body inválido"})
	}

	demoUser := os.Getenv("DEMO_USER")
	if demoUser == "" {
		demoUser = "admin"
	}
	demoPass := os.Getenv("DEMO_PASSWORD")
	if demoPass == "" {
		demoPass = "admin123"
	}

	if req.Username != demoUser || req.Password != demoPass {
		return c.Status(401).JSON(models.ErrorResponse{Error: "credenciales inválidas"})
	}

	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		return c.Status(500).JSON(models.ErrorResponse{Error: "error al generar token"})
	}

	return c.Status(200).JSON(models.LoginResponse{Token: token})
}

// JWTMiddleware valida token JWT
func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(models.ErrorResponse{Error: "No autorizado"})
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Token inválido"})
	}
	token, err := auth.ValidateToken(parts[1])
	if err != nil || !token.Valid {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Token inválido"})
	}
	return c.Next()
}
