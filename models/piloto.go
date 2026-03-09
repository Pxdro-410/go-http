package models

// Piloto en este caso representa la estructura de nuestro recurso de Fórmula 1
type Piloto struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Equipo      string `json:"equipo"`
	NumeroAuto  int    `json:"numero_auto"`
	Campeonatos int    `json:"campeonatos"`
	Activo      bool   `json:"activo"`
}

// RespuestaError sirve para manejar los errores estructurados
type RespuestaError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}