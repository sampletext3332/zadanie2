package main

import (
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const port = 8080

//go:embed index.html
var indexHTML []byte

var client = http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("starting weather app")
	log.Println("author : Mateusz Ł 101619")
	log.Println("port =", port)

	if len(os.Args) > 1 && os.Args[1] == "--health" {
		runHealthCheck()
	}

	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/weather", handleWeatherRequest)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /health")

	w.Write([]byte("OK"))
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")

	log.Println("GET /weather location =", location)

	if location == "" {
		http.Error(w, "missing location", http.StatusBadRequest)
		return
	}

	coords := findCoordinates(location)

	if coords == "" {
		log.Println("location not found:", location)
		http.Error(w, "location not found", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(coords, ",", 2)

	if len(parts) != 2 {
		http.Error(w, "invalid coordinates", http.StatusInternalServerError)
		return
	}

	data, err := requestWeatherData(parts[0], parts[1])

	if err != nil {
		log.Println("weather api error:", err)
		http.Error(w, "weather api error", http.StatusInternalServerError)
		return
	}

	log.Println("weather fetched successfully")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(`
		<h1>Weather Result</h1>

		<p>Temperature: ` + data[0] + ` °C</p>
		<p>Wind speed: ` + data[2] + ` km/h</p>
		<p>Humidity: ` + data[3] + `%</p>
		<p>Cloud cover: ` + data[4] + `%</p>

		<br>

		<a href="/">Back</a>
	`))
}

func findCoordinates(location string) string {
	switch location {
	case "Lublin,Poland":
		return "51.2465,22.5684"

	case "Kielce,Poland":
		return "50.8661,20.6286"

	case "Warsaw,Poland":
		return "52.2297,21.0122"
	}

	return ""
}

func requestWeatherData(lat, lon string) ([5]string, error) {
	log.Println("requestWeatherData:", lat, lon)

	resp, err := client.Get(
		"https://api.open-meteo.com/v1/forecast" +
			"?latitude=" + lat +
			"&longitude=" + lon +
			"&current=temperature_2m,precipitation,wind_speed_10m,relative_humidity_2m,cloud_cover",
	)

	if err != nil {
		return [5]string{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return [5]string{}, err
	}

	jsonStr := string(body)

	start := strings.Index(jsonStr, `"current":{`)

	if start >= 0 {
		jsonStr = jsonStr[start:]
	}

	return [5]string{
		extractJSONValue(jsonStr, "temperature_2m"),
		extractJSONValue(jsonStr, "precipitation"),
		extractJSONValue(jsonStr, "wind_speed_10m"),
		extractJSONValue(jsonStr, "relative_humidity_2m"),
		extractJSONValue(jsonStr, "cloud_cover"),
	}, nil
}

func extractJSONValue(data, key string) string {
	searchKey := `"` + key + `":`

	index := strings.Index(data, searchKey)

	if index < 0 {
		return ""
	}

	data = strings.TrimLeft(data[index+len(searchKey):], " ")

	if len(data) > 0 && data[0] == '"' {
		data = data[1:]

		if end := strings.Index(data, `"`); end >= 0 {
			return data[:end]
		}

		return ""
	}

	if end := strings.IndexAny(data, ",}"); end >= 0 {
		return data[:end]
	}

	return data
}

// endpoint do healthcheck
func runHealthCheck() {
	log.Println("healthcheck running")

	resp, err := client.Get(
		"http://localhost:" + strconv.Itoa(port) + "/health",
	)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("healthcheck failed")
		os.Exit(1)
	}

	log.Println("healthcheck OK")

	os.Exit(0)
}
