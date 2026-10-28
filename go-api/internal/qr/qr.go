package qr

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Factorize realiza factorización QR usando gonum
func Factorize(A [][]float64) ([][]float64, [][]float64, error) {
	m := len(A)
	if m == 0 {
		return nil, nil, fmt.Errorf("matriz vacía")
	}
	n := len(A[0])
	if n == 0 {
		return nil, nil, fmt.Errorf("matriz inválida")
	}
	for i := 0; i < m; i++ {
		if len(A[i]) != n {
			return nil, nil, fmt.Errorf("matriz no rectangular")
		}
	}
	if m < n {
		return nil, nil, fmt.Errorf("matriz debe tener al menos tantas filas como columnas")
	}

	data := make([]float64, 0, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			data = append(data, A[i][j])
		}
	}
	matA := mat.NewDense(m, n, data)

	var qr mat.QR
	qr.Factorize(matA)

	// Obtener Q y R completos
	Qm := new(mat.Dense)
	Rm := new(mat.Dense)
	qr.QTo(Qm)
	qr.RTo(Rm)

	// Para economía: Q es m x n, R es n x n
	Q := mat.NewDense(m, n, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			Q.Set(i, j, Qm.At(i, j))
		}
	}
	R := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			R.Set(i, j, Rm.At(i, j))
		}
	}

	qRes := make([][]float64, m)
	for i := 0; i < m; i++ {
		qRes[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			qRes[i][j] = Q.At(i, j)
		}
	}

	rRes := make([][]float64, n)
	for i := 0; i < n; i++ {
		rRes[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			rRes[i][j] = R.At(i, j)
		}
	}

	return qRes, rRes, nil
}

func Multiply(Q, R [][]float64) [][]float64 {
	m := len(Q)
	n := len(Q[0])
	p := len(R[0])
	res := make([][]float64, m)
	for i := 0; i < m; i++ {
		res[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			var sum float64
			for k := 0; k < n; k++ {
				sum += Q[i][k] * R[k][j]
			}
			res[i][j] = sum
		}
	}
	return res
}

func AlmostEqual(a, b [][]float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if math.Abs(a[i][j]-b[i][j]) > tol {
				return false
			}
		}
	}
	return true
}
