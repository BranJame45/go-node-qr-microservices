package models

// QRRequest representa la solicitud de factorización QR
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// Statistics representa las métricas estadísticas
type Statistics struct {
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	Average    float64 `json:"average"`
	Sum        float64 `json:"sum"`
	IsDiagonal bool    `json:"isDiagonal"`
}

// QRResponse representa la respuesta completa
type QRResponse struct {
	Q           [][]float64          `json:"Q"`
	R           [][]float64          `json:"R"`
	Statistics  *StatisticsResponse `json:"statistics"`
}

// StatisticsResponse estructura de estadísticas para Go
type StatisticsResponse struct {
	Global    Statistics            `json:"global"`
	PerMatrix map[string]Statistics `json:"perMatrix"`
}

// NodeStatsRequest para enviar a Node API
type NodeStatsRequest struct {
	Matrices map[string][][]float64 `json:"matrices"`
}

// LoginRequest
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse
type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
