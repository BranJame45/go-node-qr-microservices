package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go-api/internal/handlers"
	"go-api/internal/models"
	"go-api/internal/qr"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func newApp() *fiber.App {
	app := fiber.New()

	// CORS: permite que el frontend (otro puerto/origen) consuma la API desde el navegador
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// POST /api/login (sin auth)
	app.Post("/api/login", handlers.LoginHandler)

	// POST /api/qr (requiere JWT)
	app.Post("/api/qr", handlers.JWTMiddleware, handleQR)

	return app
}

func main() {
	app := newApp()
	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Go corriendo en puerto %s", port)
	log.Fatal(app.Listen(":" + port))
}

func handleQR(c *fiber.Ctx) error {
	var req models.QRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "body inválido"})
	}

	// Validar matriz
	if len(req.Matrix) == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "matriz vacía"})
	}
	m := len(req.Matrix)
	n := len(req.Matrix[0])
	if n == 0 {
		return c.Status(400).JSON(models.ErrorResponse{Error: "matriz inválida"})
	}
	for i := 0; i < m; i++ {
		if len(req.Matrix[i]) != n {
			return c.Status(400).JSON(models.ErrorResponse{Error: "matriz no rectangular"})
		}
	}
	if m < n {
		return c.Status(400).JSON(models.ErrorResponse{Error: "matriz debe tener al menos tantas filas como columnas"})
	}

	// Calcular QR
	Q, R, err := qr.Factorize(req.Matrix)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	// Llamar a Node API reenviando token
	authHeader := c.Get("Authorization")
	stats, err := callNodeStats(Q, R, authHeader)
	if err != nil {
		return c.Status(502).JSON(models.ErrorResponse{Error: "la API Node no respondió"})
	}

	resp := models.QRResponse{
		Q:          Q,
		R:          R,
		Statistics: stats,
	}
	return c.Status(200).JSON(resp)
}

func callNodeStats(Q, R [][]float64, authHeader string) (*models.StatisticsResponse, error) {
	nodeURL := os.Getenv("NODE_API_URL")
	if nodeURL == "" {
		nodeURL = "http://localhost:3000"
	}
	reqBody := models.NodeStatsRequest{
		Matrices: map[string][][]float64{
			"Q": Q,
			"R": R,
		},
	}
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", nodeURL+"/api/statistics", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("node error: %s", string(body))
	}
	var stats models.StatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
