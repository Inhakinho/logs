package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	sendLogToLoki("¡Hola amigo!", "info", "test")
	sendLogToLoki("¡Debug amigo!", "debug", "test")
	sendLogToLoki("¡Advertencia amigo!", "warn", "test")
	sendLogToLoki("¡Error amigo!", "error", "test")
	sendLogToLoki("¡Fatal amigo!", "fatal", "test")

	getLogsFromLoki(`{app="test"}`)
	getLatestLogsFromLokiInstant(`{app="test"}`)
}

func sendLogToLoki(message, level, app string) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"level": level,
					"app":   app,
				},
				"values": [][]string{
					{timestamp, message},
				},
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:3100/loki/api/v1/push", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error al enviar a Loki:", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("Log enviado:", message, "| Nivel:", level, "| Status:", resp.Status)
	}
}

func getLogsFromLoki(logql string) {
	baseUrl := "http://localhost:3100/loki/api/v1/query_range"

	now := time.Now()
	start := now.Add(time.Minute * -5).UnixNano()
	end := now.UnixNano()

	params := url.Values{}
	params.Set("query", logql)

	params.Set("limit", "10")
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("end", fmt.Sprintf("%d", end))
	params.Set("direction", "backward")

	fullUrl := baseUrl + "?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Error al consultar logs:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		fmt.Println("Error al parsear respuesta:", err)
		return
	}

	fmt.Println("\nLogs recibidos desde Loki:")

	results := parsed["data"].(map[string]interface{})["result"].([]interface{})
	for _, r := range results {
		entry := r.(map[string]interface{})
		values := entry["values"].([]interface{})
		for _, val := range values {
			pair := val.([]interface{})
			timestamp := pair[0].(string)
			message := pair[1].(string)
			fmt.Printf("%s ms |====> %s\n", timestamp, message)
		}
	}
}

func getLatestLogsFromLokiInstant(logql string) {
	baseUrl := "http://localhost:3100/loki/api/v1/query"

	params := url.Values{}
	params.Set("query", logql)
	//params.Set("limit", "500") //controlar limite

	fullUrl := baseUrl + "?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Error al consultar logs:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		fmt.Println("Error al parsear respuesta:", err)
		return
	}

	fmt.Println("\nLogs recientes desde Loki (consulta instantánea):")

	results := parsed["data"].(map[string]interface{})["result"].([]interface{})
	for _, r := range results {
		entry := r.(map[string]interface{})
		values := entry["values"].([]interface{})
		for _, val := range values {
			pair := val.([]interface{})
			timestamp := pair[0].(string)
			message := pair[1].(string)
			fmt.Printf("%s ms |====> %s\n", timestamp, message)
		}
	}
}
