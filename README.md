# Loki Logging Test in Go

Este proyecto es una pequeña prueba de cómo:

- Generar logs tanto desde código Go como desde consola.
- Enviar logs a **Grafana Loki** mediante HTTP.
- Consultar logs históricos (rango de tiempo).
- Consultar logs recientes (consulta instantánea).
- Usar **Promtail** para enviar automáticamente logs de archivos locales a Loki.
- Permitir recuperar logs filtrados por un `fileUUID`.


---

## 📑 Tabla de contenidos

- [📦 Estructura del proyecto](#-estructura-del-proyecto)
- [🚀 Cómo ejecutar](#-cómo-ejecutar)
- [🐳 Docker Compose](#-docker-compose)
- [📂 Promtail](#-promtail)
- [🧪 Consultas LogQL](#-consultas-logql)
- [📘 Requisitos](#-requisitos)
- [📦 Licencia](#-licencia)

---

## 🧱 Estructura del proyecto

### `main.go`

Este archivo contiene funciones clave para:

#### 🔸 `sendLogToLoki(message, level, app)`

Envía logs directamente a Loki vía HTTP, con etiquetas:
- `level`: nivel de log (`info`, `debug`, `error`, etc.)
- `app`: nombre de la app (por ejemplo `"test"`)
- `file_uuid`: generado automáticamente por `ulid`

#### 🔸 `getLogsByUUID(uuid)`

Consulta logs históricos desde Loki usando LogQL con la etiqueta `file_uuid`.

#### 🔸 `generateTestLogs()`

Genera múltiples logs de distintos niveles y los envía a Loki.

#### 🔸 `handleGetLogsByUUID`

API HTTP que permite recuperar logs usando un `file_uuid`.

---

## 🚀 Cómo ejecutar

1. Clona este repositorio:

```bash
git clone https://github.com/tuusuario/logs.git
cd logs
```

2. Levanta **Grafana**, **Loki** y **Promtail** con Docker Compose:

```bash
docker-compose up -d
```

Esto iniciará:

- Loki en `http://localhost:3100`
- Grafana en `http://localhost:3000` (usuario: `admin`, contraseña: `admin`)
- Promtail levantado en `http://localhost:9080` recolectando logs de `./promtail-test/logs`

3. Ejecuta el programa Go para enviar y consultar logs:

```bash
go run main.go
```

Verás en consola:
- Mensajes enviados a Loki.
- Logs recuperados con consultas LogQL.

4. También puedes generar logs directamente desde consola:

```bash
echo '{"file_uuid":"01HYTESTUUID1234567877","level":"info","msg":"EJEMPLO DE LOG1"}' >> promtail-test/logs/testfile.log
echo '{"file_uuid":"01HYTESTUUID1234567877","level":"debug","msg":"EJEMPLO DE LOG1"}' >> promtail-test/logs/testfile.log
echo '{"file_uuid":"01HYTESTUUID1234567877","level":"warn","msg":"EJEMPLO DE LOG1"}' >> promtail-test/logs/testfile.log
echo '{"file_uuid":"01HYTESTUUID1234567877","level":"error","msg":"EJEMPLO DE LOG1"}' >> promtail-test/logs/testfile.log
```
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
    volumes:
      - ./loki-data:/loki

  grafana:
    image: grafana/grafana:10.4.2
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

  promtail:
    image: grafana/promtail:2.9.0
    volumes:
      - ./promtail-test/logs:/var/log/mylogs
      - ./promtail-config.yml:/etc/promtail/config.yml
    command: -config.file=/etc/promtail/config.yml
    depends_on:
      - loki
```

> **Nota:** Una vez iniciado Grafana, debes configurar Loki como "Data Source" apuntando a `http://loki:3100`.

---

## 📂 Promtail

Archivo de configuración típico: `promtail-config.yml`

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: local-logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: local-logs
          __path__: /var/log/mylogs/*.log
    pipeline_stages:
      - json:
          expressions:
            file_uuid: ""
            level: ""
            msg: ""
      - labels:
          file_uuid:
          level:
```

Este Promtail recoge automáticamente logs en formato JSON desde archivos `.log` y extrae etiquetas como `file_uuid` y `level`.

---

## 🧪 Consultas LogQL

Ejemplo de consulta usada en el código:

Consulta básica por app:
```logql
{app="test"}
```

Esto selecciona todos los logs etiquetados con `app="test"`.

Consulta específica por UUID:

```logql
{file_uuid="01HYTESTUUID1234567894"}
```
Esto selecciona todos los logs etiquetados con `file_uuid="01HYTESTUUID1234567894"`.

---

## 📘 Requisitos

- Go 1.21.6 o superior
- Docker y Docker Compose
- Acceso local a los puertos `9080`(Promtail) `3100` (Loki) y `3000` (Grafana) 

---

## 📦 Licencia

MIT
