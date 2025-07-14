# Loki Logging Test in Go

Este proyecto es una pequeña prueba de cómo:

- Enviar logs a **Grafana Loki** mediante HTTP.
- Consultar logs históricos (rango de tiempo).
- Consultar logs recientes (consulta instantánea).
- Todo ello implementado en Go.

---

## 🧱 Estructura del proyecto

### `main.go`

Este archivo contiene 3 funciones principales:

#### 🔸 `sendLogToLoki(message, level, app)`

Envía logs a Loki en formato JSON con las siguientes etiquetas:
- `level`: nivel de log (info, debug, error, etc.)
- `app`: nombre de la app (por ejemplo `"test"`)

#### 🔸 `getLogsFromLoki(logql)`

Consulta los logs de los últimos 5 minutos usando la API `/loki/api/v1/query_range`.

- Usa el lenguaje LogQL.
- Muestra todos los mensajes recuperados de forma estructurada.

#### 🔸 `getLatestLogsFromLokiInstant(logql)`

Consulta los logs **recientes** (último estado), sin definir rango de tiempo, mediante `/loki/api/v1/query`.

---

## 🚀 Cómo ejecutar

1. Clona este repo:
   ```bash
   git clone https://github.com/tuusuario/logs.git
   cd logs
