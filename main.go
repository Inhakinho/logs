package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Level     string `json:"level,omitempty"`
	App       string `json:"app,omitempty"`
}

var totalSent int

func main() {

	go generateTestLogs()

	// getLogsFromLoki(`{app="test"}`)
	// getLatestLogsFromLokiInstant(`{app="test"}`)

	router := gin.Default()
	router.GET("/logs/:uuid", handleGetLogsByUUID)
	router.Run(":8085") // escucha en localhost:8085
}

func generateTestLogs() {
	for i := 1; i <= 5; i++ {
		fmt.Println("Grupo de logs:", i)
		sendLogToLoki("¡Hola amigo!", "info", "test")
		sendLogToLoki("¡Debug amigo!", "debug", "test")
		sendLogToLoki("¡Advertencia amigo!", "warn", "test")
		sendLogToLoki("¡Error amigo!", "error", "test")
		sendLogToLoki("¡Fatal amigo!", "fatal", "test")
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Total de logs enviados: %d\n", totalSent)
}

func sendLogToLoki(message, level, app string) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	entropy := ulid.Monotonic(rand.Reader, 0)
	uuid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"level":    level,
					"app":      app,
					"fileUUID": uuid,
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
		totalSent++
		fmt.Println("Log enviado:", message, "| Nivel:", level, "| UUID:", uuid, "| Status:", resp.Status)
	}
}

func handleGetLogsByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	logs, err := getLogsByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"uuid": uuid,
		"logs": logs,
	})
}

func getLogsByUUID(fileUUID string) ([]LogEntry, error) {
	baseUrl := "http://localhost:3100/loki/api/v1/query_range"

	now := time.Now()
	start := now.Add(-10 * time.Minute).UnixNano()
	end := now.UnixNano()

	logql := fmt.Sprintf(`{fileUUID="%s"}`, fileUUID)

	params := url.Values{}
	params.Set("query", logql)
	params.Set("start", fmt.Sprintf("%d", start))
	params.Set("end", fmt.Sprintf("%d", end))
	params.Set("limit", "100")
	params.Set("direction", "backward")

	fullUrl := baseUrl + "?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("fallo en HTTP GET: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("fallo en parsear JSON: %w", err)
	}

	results := parsed["data"].(map[string]interface{})["result"].([]interface{})
	var entries []LogEntry

	for _, r := range results {
		entry := r.(map[string]interface{})

		labels := entry["stream"].(map[string]interface{})
		app := labels["app"].(string)
		level := labels["level"].(string)

		values := entry["values"].([]interface{})
		for _, val := range values {
			pair := val.([]interface{})
			timestampStr := pair[0].(string)
			message := pair[1].(string)

			//timestamp de nanosegundos a string
			nanos, _ := strconv.ParseInt(timestampStr, 10, 64)
			formattedTime := time.Unix(0, nanos).Format("02/01/2006 15:04:05")

			entries = append(entries, LogEntry{
				Timestamp: formattedTime,
				Message:   message,
				App:       app,
				Level:     level,
			})
		}
	}

	return entries, nil
}

// func getLogsFromLoki(logql string) {
// 	baseUrl := "http://localhost:3100/loki/api/v1/query_range"

// 	now := time.Now()
// 	start := now.Add(time.Minute * -5).UnixNano()
// 	end := now.UnixNano()

// 	params := url.Values{}
// 	params.Set("query", logql)

// 	params.Set("limit", "10") //limit de logs
// 	params.Set("start", fmt.Sprintf("%d", start))
// 	params.Set("end", fmt.Sprintf("%d", end))
// 	params.Set("direction", "backward")

// 	fullUrl := baseUrl + "?" + params.Encode()

// 	resp, err := http.Get(fullUrl)
// 	if err != nil {
// 		fmt.Println("Error al consultar logs:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var parsed map[string]interface{}
// 	if err := json.Unmarshal(body, &parsed); err != nil {
// 		fmt.Println("Error al parsear respuesta:", err)
// 		return
// 	}

// 	fmt.Println("\nLogs recibidos desde Loki:")
// 	results := parsed["data"].(map[string]interface{})["result"].([]interface{})

// 	logCount := 0

// 	for _, r := range results {
// 		entry := r.(map[string]interface{})
// 		values := entry["values"].([]interface{})
// 		logCount += len(values)

// 		for _, val := range values {
// 			pair := val.([]interface{})
// 			timestamp := pair[0].(string)
// 			message := pair[1].(string)
// 			fmt.Printf("%s ms |====> %s\n", timestamp, message)
// 		}
// 	}

// 	fmt.Printf("Total logs recuperados (query_range): %d\n", logCount)
// }

// func getLatestLogsFromLokiInstant(logql string) {
// 	baseUrl := "http://localhost:3100/loki/api/v1/query"

// 	params := url.Values{}
// 	params.Set("query", logql)
// 	//params.Set("limit", "500") //controlar limite

// 	fullUrl := baseUrl + "?" + params.Encode()

// 	resp, err := http.Get(fullUrl)
// 	if err != nil {
// 		fmt.Println("Error al consultar logs:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)

// 	var parsed map[string]interface{}
// 	if err := json.Unmarshal(body, &parsed); err != nil {
// 		fmt.Println("Error al parsear respuesta:", err)
// 		return
// 	}

// 	fmt.Println("\nLogs recientes desde Loki (consulta instantánea):")
// 	results := parsed["data"].(map[string]interface{})["result"].([]interface{})

// 	logCount := 0

// 	for _, r := range results {
// 		entry := r.(map[string]interface{})
// 		values := entry["values"].([]interface{})
// 		logCount = logCount + len(values)

// 		for _, val := range values {
// 			pair := val.([]interface{})
// 			timestamp := pair[0].(string)
// 			message := pair[1].(string)
// 			fmt.Printf("%s ms |====> %s\n", timestamp, message)
// 		}
// 	}

// 	fmt.Printf("Total logs recuperados (consulta instantánea): %d\n", logCount)
// }
