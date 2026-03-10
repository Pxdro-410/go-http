# Api restful de f1
## Pedro Caso - 241286

Esta API fue contruida con GO y permite gestionar un catálogo de pilotos de Fórmula 1, demostrando el manejo de métodos HTTP, parámetros en la ruta, validaciones y persistencia de datos en un archivo local. 

## Características Implementadas
- Construido puramente con Go estándar.
- Contiene persistencia, es decir que los cambios que se realicen con POST y DELETE se guardan y actualizan en `data/pilotos.json`.
- Soporte para múltiples query parameters combinados.
- Rutas dinámicas (/api/items/1).
- Todas las respuestas de error devuelven un JSON consistente con su código HTTP.

## Cómo ejecutar el proyecto

### El proyecto está configurado para correr utilizando Docker. El puerto donde se ejecuta es mi carnet 241286, sin embargo, por ser un carnet de 6 digitos se utiliza el puerto **41286** (Autorizado por Erick).
Pasos a seguir:
1. Clona este repositorio.
2. Abre una terminal en la raíz del proyecto.
3. Ejecuta el siguiente comando: docker compose up

## Documentación de Endpoints y Ejemplos de Uso
la estructura de pilotos.json es la siguiente:
{
  "id": 1,
  "nombre": "Max Verstappen",
  "equipo": "Red Bull Racing",
  "numero_auto": 1,
  "campeonatos": 3,
  "activo": true
}

## Utilizando postman, se realizan las peticiones HTTP
### utilizar metodo GET

1. Obtener todos los pilotos
Ruta: /api/items
Esto devuelve la lista completa de pilotos registrados.

2. Obtener un piloto por ID (Path Parameter)
Ruta: /api/items/1
Devuelve únicamente el piloto que coincida con el ID proporcionado en la URL.

3. Filtrar pilotos (Query Parameters Combinados)
Ruta: /api/items?equipo=Ferrari&activo=true
Permite filtrar la lista usando múltiples parámetros.
nota: solo soporta los parametros id, equipo y activo.

### utilizar método POST
1. Crear un nuevo piloto
Ruta: /api/items
Registra un nuevo piloto y lo guarda en el archivo JSON. Valida que el nombre y equipo no estén vacíos. En este caso el id se asigna automáticamente.

ejemplo de peticion:
{
  "nombre": "Franco Colapinto",
  "equipo": "Alpine",
  "numero_auto": 43,
  "campeonatos": 0,
  "activo": true
}

### utilizar método DELETE
Ruta: /api/items/9 (ejemplo)
Esto elimina físicamente el registro del piloto con el ID especificado.

### Para dejar de ejecutar el contenedor
1. Abre una terminal en la raíz del proyecto.
2. Ejecuta el siguiente comando: docker compose down
