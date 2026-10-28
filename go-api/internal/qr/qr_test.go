package qr

import (
	"testing"
)

func TestQRFactorization(t *testing.T) {
	// Matriz de prueba
	A := [][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}

	Q, R, err := Factorize(A)
	if err != nil {
		t.Fatalf("Error al factorizar: %v", err)
	}

	// Verificar dimensiones
	if len(Q) != 3 || len(Q[0]) != 3 {
		t.Fatalf("Dimensiones de Q incorrectas")
	}
	if len(R) != 3 || len(R[0]) != 3 {
		t.Fatalf("Dimensiones de R incorrectas")
	}

	// Verificar Q*R ≈ A
	QR := Multiply(Q, R)
	if !AlmostEqual(A, QR, 1e-9) {
		t.Fatalf("Q*R no es igual a A dentro de tolerancia")
	}
}

func TestQRRectangular(t *testing.T) {
	A := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	Q, R, err := Factorize(A)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	QR := Multiply(Q, R)
	if !AlmostEqual(A, QR, 1e-9) {
		t.Fatalf("No coincide")
	}
}
