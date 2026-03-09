package main

import (
	"log"
	"net/http"
	"api-f1/handlers"
	"api-f1/storage"
)

func main() {
	// Cargar los datos desde la base de datos local creadaa
	err := storage.LoadPilotos()
	if err != nil {
		log.Fatalf("Error al cargar el archivo JSON de pilotos: %v", err)
	}

	// registrar los endpoints y sus handlers correspondientes, se define con api/items 
	http.HandleFunc("/api/items", handlers.ItemsHandler)

	// se confira el puerto en el puerto 41286 
	// favor de tomar en cuenta que mi carnet es el 241286, pero el puerto 241286 no es valido 
	puerto := "41286"
	log.Printf("API de Fórmula 1 Iniciada")
	log.Printf("Servidor corriendo en el puerto :%s", puerto)

	// se levantara el servidor HTTP 
	err = http.ListenAndServe(":"+puerto, nil)
	if err != nil {
		log.Fatalf("Error critico al iniciar el servidor: %v", err)
	}
}