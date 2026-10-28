# Factorización QR — Microservicios Go + Node

Dos APIs RESTful que se comunican por HTTP, contenerizadas con Docker y desplegadas en AWS:

- **API Go (Fiber):** recibe una matriz rectangular, calcula su **factorización QR** (matrices Q y R) con `gonum`, reenvía Q y R a la API Node y devuelve Q, R y las estadísticas. Actúa como orquestador.
- **API Node (Express):** recibe Q y R y calcula estadísticas: **máximo, mínimo, promedio, suma total** y si **alguna matriz es diagonal**.
- **Frontend (Vite + React):** interfaz con login (JWT) para ingresar una matriz y ver Q, R y las estadísticas.
- **Seguridad:** JWT (HS256) protege los endpoints de ambas APIs con un secreto compartido.

> **Nota sobre el enunciado:** la sección de arquitectura menciona "rotación de la matriz", pero la funcionalidad requerida pide explícitamente **factorización QR**. Se implementa QR por ser el requisito funcional concreto; la "rotación" es texto heredado.

## Arquitectura

```
┌──────────┐      ┌─────────────────────────┐      ┌──────────────┐
│ Frontend │ ───▶ │        API Go           │ ───▶ │   API Node   │
│  (nginx) │      │  Fiber · QR con gonum   │ HTTP │   Express    │
│  :80     │ ◀─── │  orquesta · JWT · :8080 │ ◀─── │ stats · :3000│
└──────────┘      └─────────────────────────┘      └──────────────┘
```

El cliente llama a la API Go. Esta calcula la QR, envía Q y R a la API Node (reenviando el JWT), recibe las estadísticas y devuelve todo junto. La API Node es interna: solo la consume la API Go.

## Estructura

```
/
├── docker-compose.yml
├── .env.example
├── README.md
├── go-api/            # API Go (Fiber)
│   ├── Dockerfile
│   ├── main.go
│   └── internal/{qr,auth,handlers,models}
├── node-api/          # API Node (Express)
│   ├── Dockerfile
│   └── src/{index.js,routes,services,middleware}
└── frontend/          # Vite + React
    ├── Dockerfile
    └── src/
```

## Variables de entorno

Copiar `.env.example` a `.env`:

| Variable | Descripción | Valor por defecto |
|---|---|---|
| `JWT_SECRET` | Secreto compartido HS256 (Go y Node) | `change-me-super-secret` |
| `GO_PORT` | Puerto de la API Go | `8080` |
| `NODE_API_URL` | URL de la API Node (nombre de servicio en Docker) | `http://node-api:3000` |
| `DEMO_USER` | Usuario para el login demo | `admin` |
| `DEMO_PASSWORD` | Contraseña para el login demo | `admin123` |
| `NODE_PORT` | Puerto de la API Node | `3000` |
| `VITE_GO_API_URL` | URL de la API Go para el frontend (build-time) | `http://localhost:8080` |

> `VITE_GO_API_URL` se resuelve en **tiempo de build** del frontend. Para desplegar en un servidor, debe apuntar a la IP/host público (ej. `http://<IP_PUBLICA>:8080`) **antes** de construir la imagen.

## Cómo correr

### Con Docker (recomendado)

```bash
cp .env.example .env
docker-compose up -d --build
```

- Frontend: http://localhost
- API Go: http://localhost:8080
- API Node: http://localhost:3000

### En local (sin Docker)

```bash
# API Node
cd node-api && npm install && npm start

# API Go (en otra terminal)
cd go-api && go run main.go

# Frontend (en otra terminal)
cd frontend && npm install && npm run dev
```

## API Reference

### API Go (base: `http://<host>:8080`)

#### `POST /api/login`
Emite un JWT. No requiere autenticación.

- **Body:** `{ "username": "admin", "password": "admin123" }`
- **200:** `{ "token": "<jwt>" }`
- **401:** credenciales inválidas.

#### `POST /api/qr` 🔒 requiere `Authorization: Bearer <token>`
Calcula la factorización QR y devuelve Q, R y las estadísticas (obtenidas de la API Node).

- **Body:** `{ "matrix": [[12,-51,4],[6,167,-68],[-4,24,-41]] }`
  (matriz rectangular de números, con filas ≥ columnas)
- **200:**
```json
{
  "Q": [[...], ...],
  "R": [[...], ...],
  "statistics": {
    "global":    { "max": 0, "min": 0, "average": 0, "sum": 0, "isDiagonal": false },
    "perMatrix": {
      "Q": { "max": 0, "min": 0, "average": 0, "sum": 0, "isDiagonal": false },
      "R": { "max": 0, "min": 0, "average": 0, "sum": 0, "isDiagonal": false }
    }
  }
}
```
- **400:** matriz inválida (vacía, no rectangular, no numérica, o filas < columnas).
- **401:** token ausente o inválido.
- **502:** la API Node no respondió.

#### `GET /health`
`{ "status": "ok" }` — sin autenticación.

### API Node (base: `http://<host>:3000`, interna)

#### `POST /api/statistics` 🔒 requiere `Authorization: Bearer <token>`
Calcula estadísticas sobre las matrices recibidas.

- **Body:** `{ "matrices": { "Q": [[...]], "R": [[...]] } }`
- **200:** mismo objeto `{ global, perMatrix }` que se muestra arriba.
  - `max`, `min`, `sum`, `average` sobre el conjunto de valores.
  - `isDiagonal`: `true` si alguna matriz es cuadrada y todo fuera de la diagonal es ≈ 0.
- **400:** body inválido. **401:** token ausente o inválido.

#### `GET /health`
`{ "status": "ok" }`.

### Ejemplo con curl

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r .token)

# 2. Factorización QR + estadísticas
curl -X POST http://localhost:8080/api/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix":[[12,-51,4],[6,167,-68],[-4,24,-41]]}'
```

## Seguridad (JWT)

- La API Go emite el token en `/api/login` (credenciales demo por variables de entorno).
- `/api/qr` (Go) y `/api/statistics` (Node) exigen `Authorization: Bearer <token>`.
- Ambos servicios validan la firma con el **mismo** `JWT_SECRET`. La API Go reenvía el token a la API Node.
- La API Node no se expone al exterior en el despliegue (su puerto no se abre en el grupo de seguridad): solo la consume la API Go por la red interna de Docker.

## Pruebas

```bash
# API Go: unitarias (QR) + integración del handler /api/qr
cd go-api && go test ./... -count=1

# API Node: endpoint y cálculo de estadísticas (Jest + Supertest)
cd node-api && npm test
```

## Despliegue en AWS EC2

1. Lanzar una instancia **EC2 Ubuntu 22.04** (`t2.micro`) con un *security group* que abra los puertos **22**, **80** y **8080**.
2. Conectarse por SSH e instalar Docker:
```bash
sudo apt-get update
sudo apt-get install -y docker.io docker-compose git
```
3. Clonar el repo, configurar la IP pública en el frontend y levantar:
```bash
git clone https://github.com/BranJame45/go-node-qr-microservices.git
cd go-node-qr-microservices
cp .env.example .env
sed -i "s|VITE_GO_API_URL=.*|VITE_GO_API_URL=http://<IP_PUBLICA>:8080|" .env
sudo docker-compose up -d --build
```

El frontend queda en `http://<IP_PUBLICA>` y la API Go en `http://<IP_PUBLICA>:8080`.

## Notas

- **QR vs rotación:** se implementa factorización QR (algoritmo de `gonum`). Q es ortogonal (m×n), R triangular superior (n×n), con A = Q·R.
- **CORS:** la API Go habilita CORS para que el frontend (servido en otro puerto) pueda consumirla desde el navegador.
- Sin secretos hardcodeados: toda la configuración viene por variables de entorno.
