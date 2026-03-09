package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-f1/models"
	"api-f1/storage"
)

// funcion para escribir respuestas JSON de forma consistente
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// función para escribir errores de forma estructurada
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.RespuestaError{Error: message, Status: status})
}

// funcion que maneja rutas
func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetItems(w, r)
	case http.MethodPost:
		handleCreateItem(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
	}
}

// handleGetItems maneja la obtención de pilotos
func handleGetItems(w http.ResponseWriter, r *http.Request) {
	//  se leen los query parameters  
	query := r.URL.Query()
	idParam := query.Get("id")

	// si se proporciona un ID se busca ese piloto específico
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			writeError(w, http.StatusBadRequest, "El ID debe ser un numero entero dentro del rango permitido")
			return
		}

		// Buscar piloto por ID
		for _, piloto := range storage.Pilotos {
			if piloto.ID == id {
				writeJSON(w, http.StatusOK, piloto)
				return
			}
		}
		writeError(w, http.StatusNotFound, "No se encontro ninguna coincidencia para el id proporcionado")
		return
	}

	// si no se proporcionan parametros, se devuelve la lista completa de pilotos
	writeJSON(w, http.StatusOK, storage.Pilotos)
}

// esta función maneja la creación de un nuevo piloto
func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var nuevoPiloto models.Piloto

	// se decodifica el cuerpo de la petición JSON en la estructura Piloto
	if err := json.NewDecoder(r.Body).Decode(&nuevoPiloto); err != nil {
		writeError(w, http.StatusBadRequest, "El formato del JSON es incorrecto")
		return
	}

	// se hace una validacion de datos
	if nuevoPiloto.Nombre == "" || nuevoPiloto.Equipo == "" {
		writeError(w, http.StatusBadRequest, "Los campos de nombre y equipo son obligatorios")
		return
	}

	// se asigna un ID unico al nuevo piloto mediante autoincremento
	nuevoPiloto.ID = generarSiguienteID()

	// se agrega el nuevo piloto a la base de datos local
	storage.Pilotos = append(storage.Pilotos, nuevoPiloto)
	if err := storage.SavePilotos(); err != nil {
		writeError(w, http.StatusInternalServerError, "Error interno al guardar el piloto, intente nuevamente")
		return
	}

	// se devulve el codigo 2xx con el nuevo piloto creado
	writeJSON(w, http.StatusCreated, nuevoPiloto)
}

// funcion para autoincrementar el id cada vez que se crea un piloto
func generarSiguienteID() int {
	maxID := 0
	for _, piloto := range storage.Pilotos {
		if piloto.ID > maxID {
			maxID = piloto.ID
		}
	}
	return maxID + 1
}