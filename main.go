package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables file")
	}

	port := os.Getenv("PORT")
	apiKey := os.Getenv("API_KEY")
	tokenUrl := os.Getenv("TOKEN_URL")
	apiUrl := os.Getenv("API_URL")
	initialToken := ""

	cfg := cfg{
		port:     port,
		apiKey:   apiKey,
		tokenUrl: tokenUrl,
		apiUrl:   apiUrl,
	}

	client := &http.Client{
		Transport: &authRoundTripper{
			next:       &loggingRoundTripper{next: http.DefaultTransport, logger: os.Stdout},
			token:      initialToken,
			maxRetries: 3,
			retryDelay: time.Second * 2,
			cfg:        &cfg,
		},
	}

	app := app{
		client: client,
		cfg:    &cfg,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /tempoEspera/Estacao/todos", app.handleAllWaitingTimes)
	mux.HandleFunc("GET /infoEstacao/{estacao}", app.handleStationInfo)
	mux.HandleFunc("GET /infoEstacao/todos", app.handleAllStations)
	mux.HandleFunc("GET /estadoLinha/todos", app.handleAllLinesInfo)
	mux.HandleFunc("GET /estadoLinha/{linha}", app.handleLineInfo)
	mux.HandleFunc("GET /tempoEspera/Linha/{linha}", app.handleLineWaitingTimes)
	mux.HandleFunc("GET /tempoEspera/Estacao/{estacao}", app.handleStationWaitingTimes)
	mux.HandleFunc("GET /infoDestinos/todos", app.handleAllDestinations)

	log.Println("Listening on port " + app.cfg.port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.port, mux))
}
