# 📂 Carpeta de Logs

Este directorio contiene los archivos `.log` que son recolectados automáticamente por **Promtail** y enviados a **Grafana Loki**.

## 📌 Detalles:

- Promtail está configurado para leer archivos `.log` ubicados en esta carpeta.
- Cada línea debe estar en formato JSON con las siguientes claves:
  - `file_uuid`: identificador único del archivo o evento.
  - `level`: nivel del log (`DEBUG`, `INFO`, `ERROR`, etc.).
  - `msg`: mensaje del log.

## 📝 Ejemplo de línea válida:

```json
{"file_uuid":"01HYTESTUUID1234567894","level":"info","msg":"EJEMPLO DE LOG1"}
{"file_uuid":"01HYTESTUUID1234567894","level":"debug","msg":"EJEMPLO DE LOG1"}
{"file_uuid":"01HYTESTUUID1234567894","level":"warn","msg":"EJEMPLO DE LOG1"}
{"file_uuid":"01HYTESTUUID1234567894","level":"error","msg":"EJEMPLO DE LOG1"}
```
