package storage

import (
	"encoding/json"
	"os"
	"sync"

	// Importamos el paquete de modelos para usar la estructura Piloto
	"api-f1/models" 
)

var (
	// Pilotos es la variable global que almacenará nuestros pilotos en memoria, se maneja como una base de datos simple
	Pilotos []models.Piloto
	
	// se usa un mutex para evitar condiciones de carrera al acceder a la variable global Pilotos 
	mutex sync.Mutex 
)

// ruta del archivo JSON
const filePath = "./data/pilotos.json"

// LoadPilotos lee el archivo JSON al iniciar el servidor
func LoadPilotos() error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Pilotos)
	if err != nil {
		return err
	}
	return nil
}

// SavePilotos sobrescribe el archivo JSON con los datos actuales
func SavePilotos() error {
	mutex.Lock()
	defer mutex.Unlock()

	// usamos MarshalIndent para escribir el JSON de forma legible
	data, err := json.MarshalIndent(Pilotos, "", "  ")
	if err != nil {
		return err
	}

	// 0644 son los permisos estándar de lectura/escritura para el archivo
	return os.WriteFile(filePath, data, 0644)
}