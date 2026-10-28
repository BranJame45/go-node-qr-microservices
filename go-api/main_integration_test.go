package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-api/internal/auth"
	"go-api/internal/models"
)

func TestQRWithoutAuth(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DEMO_USER", "admin")
	os.Setenv("DEMO_PASSWORD", "admin123")

	app := newApp()

	reqBody := models.QRRequest{
		Matrix: [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestQRWithValidTokenAndValidMatrix(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DEMO_USER", "admin")
	os.Setenv("DEMO_PASSWORD", "admin123")

	// Mock Node API
	nodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/statistics" && r.Method == "POST" {
			stats := models.StatisticsResponse{
				Global: models.Statistics{Max: 10, Min: -10, Average: 1.5, Sum: 50, IsDiagonal: false},
				PerMatrix: map[string]models.Statistics{
					"Q": {Max: 1, Min: -1, Average: 0, Sum: 0, IsDiagonal: false},
					"R": {Max: 167, Min: -68, Average: 11, Sum: 33, IsDiagonal: false},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(stats)
			return
		}
		w.WriteHeader(404)
	}))
	defer nodeServer.Close()

	os.Setenv("NODE_API_URL", nodeServer.URL)

	app := newApp()

	token, err := auth.GenerateToken("admin")
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	reqBody := models.QRRequest{
		Matrix: [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var qrResp models.QRResponse
	if err := json.NewDecoder(resp.Body).Decode(&qrResp); err != nil {
		t.Fatalf("error decoding: %v", err)
	}
	if len(qrResp.Q) == 0 || len(qrResp.R) == 0 {
		t.Fatalf("Q or R missing")
	}
	if qrResp.Statistics == nil {
		t.Fatalf("statistics missing")
	}
}

func TestQRInvalidMatrixRectangular(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DEMO_USER", "admin")
	os.Setenv("DEMO_PASSWORD", "admin123")

	nodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats := models.StatisticsResponse{
			Global:    models.Statistics{Max: 0, Min: 0, Average: 0, Sum: 0, IsDiagonal: false},
			PerMatrix: map[string]models.Statistics{"Q": {Max: 0, Min: 0, Average: 0, Sum: 0, IsDiagonal: false}, "R": {Max: 0, Min: 0, Average: 0, Sum: 0, IsDiagonal: false}},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(stats)
	}))
	defer nodeServer.Close()
	os.Setenv("NODE_API_URL", nodeServer.URL)

	app := newApp()
	token, _ := auth.GenerateToken("admin")

	// Matriz no rectangular
	reqBody := models.QRRequest{
		Matrix: [][]float64{{1, 2}, {3}},
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestQRWithNodeErrorReturns502(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DEMO_USER", "admin")
	os.Setenv("DEMO_PASSWORD", "admin123")

	nodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}))
	defer nodeServer.Close()
	os.Setenv("NODE_API_URL", nodeServer.URL)

	app := newApp()
	token, _ := auth.GenerateToken("admin")

	reqBody := models.QRRequest{
		Matrix: [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
	}
	b, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/qr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != 502 {
		t.Fatalf("expected 502, got %d", resp.StatusCode)
	}
}
