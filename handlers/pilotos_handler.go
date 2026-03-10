package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"api-f1/models"
	"api-f1/storage"
)

// funciones  para escribir respuestas JSON y errores
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.RespuestaError{Error: message, Status: status})
}

// función principal que maneja todas las solicitudes 
func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	// se limpia el path para extraer el ID 
	path := strings.TrimPrefix(r.URL.Path, "/api/items")
	path = strings.TrimPrefix(path, "/")

	// se maneja el path para determinar si es una solicitud para un recurso específico o para la colección completa
	if path != "" {
		id, err := strconv.Atoi(path)
		if err != nil {
			writeError(w, http.StatusBadRequest, "El ID debe ser un entero validp")
			return
		}

		// se crean los metodos para manejar GET, PATCH y DELETE para un recurso específico
		switch r.Method {
		case http.MethodGet:
			handleGetItemByID(w, r, id)
		case http.MethodDelete:
			handleDeleteItem(w, r, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "metodo no permitido para este endpoint")
		}
		return
	}

	// Si no hay ID en la ruta, se maneja la solicitud para la coleccion con get y post
	switch r.Method {
	case http.MethodGet:
		handleGetItems(w, r)
	case http.MethodPost:
		handleCreateItem(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Metodo no permitido")
	}
}

// funcion para manejar GET sin ID, es decir, con filtros por query parameters
func handleGetItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	equipoFilter := query.Get("equipo")
	activoFilter := query.Get("activo")
	idFilter := query.Get("id")

	var resultados []models.Piloto

	for _, piloto := range storage.Pilotos {
		match := true

		// Filtro por ID
		if idFilter != "" {
			id, _ := strconv.Atoi(idFilter)
			if piloto.ID != id { match = false }
		}
		// filtro por equipo
		if equipoFilter != "" {
			// strings.EqualFold compara sin importar mayúsculas/minúsculas
			if !strings.EqualFold(piloto.Equipo, equipoFilter) { match = false }
		}
		// Filtro por Estado 
		if activoFilter != "" {
			activo, _ := strconv.ParseBool(activoFilter)
			if piloto.Activo != activo { match = false }
		}

		// si hay match con los filtros, se agrega a los resultados para mostrarlo en el writeJSON
		if match {
			resultados = append(resultados, piloto)
		}
	}

	// Si no hay resultados, se duevuelve un arreglo vacio en lugar de null
	if resultados == nil {
		resultados = []models.Piloto{}
	}

	writeJSON(w, http.StatusOK, resultados)
}

// funcion para manejar POST, es decir, crear un nuevo piloto
func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var nuevoPiloto models.Piloto
	if err := json.NewDecoder(r.Body).Decode(&nuevoPiloto); err != nil {
		writeError(w, http.StatusBadRequest, "formato JSON invalido")
		return
	}

	if nuevoPiloto.Nombre == "" || nuevoPiloto.Equipo == "" {
		writeError(w, http.StatusBadRequest, "Debe proporcionar al menos el nombre y el equipo del piloto")
		return
	}

	nuevoPiloto.ID = generarSiguienteID()
	storage.Pilotos = append(storage.Pilotos, nuevoPiloto)
	
	if err := storage.SavePilotos(); err != nil {
		writeError(w, http.StatusInternalServerError, "Error interno al guardar el piloto")
		return
	}

	writeJSON(w, http.StatusCreated, nuevoPiloto)
}

// funcion para manejar GET con ID
func handleGetItemByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, piloto := range storage.Pilotos {
		if piloto.ID == id {
			writeJSON(w, http.StatusOK, piloto)
			return
		}
	}
	writeError(w, http.StatusNotFound, "No se encontraron coindidencias para el ID proporcionado")
}

// funcion para manejar DELETE y eliminar un piloto existente
func handleDeleteItem(w http.ResponseWriter, r *http.Request, id int) {
	index := -1
	for i, p := range storage.Pilotos {
		if p.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		writeError(w, http.StatusNotFound, "El piloto que intenta eliminar no existe")
		return
	}

	// se elimina el piloto del slice utilizando append para crear un nuevo slice sin el elemento a eliminar
	storage.Pilotos = append(storage.Pilotos[:index], storage.Pilotos[index+1:]...)

	if err := storage.SavePilotos(); err != nil {
		writeError(w, http.StatusInternalServerError, "Error al guardar los cambios")
		return
	}

	// se devuelve un mensaje de exito en lugar de un cuerpo vacio para confirmar que la eliminacion fue exitosa
	writeJSON(w, http.StatusOK, map[string]string{"mensaje": "Piloto eliminado exitosamente"})
}

// funcion para generar un ID autoincrementable para los nuevos pilotos
func generarSiguienteID() int {
	maxID := 0
	for _, piloto := range storage.Pilotos {
		if piloto.ID > maxID {
			maxID = piloto.ID
		}
	}
	return maxID + 1
}