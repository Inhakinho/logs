# Loki Logging Test in Go

Este proyecto es una pequeña prueba de cómo:

- Enviar logs a **Grafana Loki** mediante HTTP.
- Consultar logs históricos (rango de tiempo).
- Consultar logs recientes (consulta instantánea).
- Todo ello implementado en Go, de forma sencilla y funcional.

---

## 📑 Tabla de contenidos

- [📦 Estructura del proyecto](#-estructura-del-proyecto)
- [🚀 Cómo ejecutar](#-cómo-ejecutar)
- [🐳 Docker Compose](#-docker-compose)
- [🧪 Consultas LogQL](#-consultas-logql)
- [📘 Requisitos](#-requisitos)
- [📦 Licencia](#-licencia)

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

1. Clona este repositorio:

```bash
git clone https://github.com/tuusuario/logs.git
cd logs
```

2. Levanta **Grafana** y **Loki** con Docker Compose:

```bash
docker-compose up -d
```

Esto iniciará:

- Loki en `http://localhost:3100`
- Grafana en `http://localhost:3000` (usuario: `admin`, contraseña: `admin`)

3. Ejecuta el programa Go para enviar y consultar logs:

```bash
go run main.go
```

Verás en consola:
- Mensajes enviados a Loki.
- Logs recuperados con consultas LogQL.

---

## 🐳 Docker Compose

Este es el archivo `docker-compose.yml` incluido en el proyecto:

```yaml
version: '3'

services:
  loki:
    image: grafana/loki:2.9.0
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml

  grafana:
    image: grafana/grafana:10.4.2
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

> **Nota:** Una vez iniciado Grafana, debes configurar Loki como "Data Source" apuntando a `http://loki:3100`.

---

## 🧪 Consultas LogQL

Ejemplo de consulta usada en el código:

```logql
{app="test"}
```

Esto selecciona todos los logs etiquetados con `app="test"`.

---

## 📘 Requisitos

- Go 1.18 o superior
- Docker y Docker Compose
- Acceso local a los puertos `3100` (Loki) y `3000` (Grafana)

---

## 📦 Licencia

MIT
