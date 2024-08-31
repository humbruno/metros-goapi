package main

import (
	"fmt"
	"net/http"
)

func (app *app) handleAllWaitingTimes(w http.ResponseWriter, req *http.Request) {
	path := "tempoEspera/Estacao/todos"
	app.handleMetroApiResponse(path, w)
}

func (app *app) handleStationInfo(w http.ResponseWriter, req *http.Request) {
	station := req.PathValue("estacao")
	path := fmt.Sprintf("infoEstacao/%s", station)
	app.handleMetroApiResponse(path, w)
}

// TODO: cache this response (directly in memory?)
func (app *app) handleAllStations(w http.ResponseWriter, req *http.Request) {
	path := "infoEstacao/todos"
	app.handleMetroApiResponse(path, w)
}

func (app *app) handleAllLinesInfo(w http.ResponseWriter, req *http.Request) {
	path := "estadoLinha/todos"
	app.handleMetroApiResponse(path, w)
}

func (app *app) handleLineInfo(w http.ResponseWriter, req *http.Request) {
	line := req.PathValue("linha")
	path := fmt.Sprintf("estadoLinha/%s", line)
	app.handleMetroApiResponse(path, w)
}

func (app *app) handleLineWaitingTimes(w http.ResponseWriter, req *http.Request) {
	line := req.PathValue("linha")
	path := fmt.Sprintf("tempoEspera/Linha/%s", line)
	app.handleMetroApiResponse(path, w)
}

func (app *app) handleStationWaitingTimes(w http.ResponseWriter, req *http.Request) {
	station := req.PathValue("estacao")
	path := fmt.Sprintf("tempoEspera/Estacao/%s", station)
	app.handleMetroApiResponse(path, w)
}

// TODO: cache this response (directly in memory?)
func (app *app) handleAllDestinations(w http.ResponseWriter, r *http.Request) {
	path := "infoDestinos/todos"
	app.handleMetroApiResponse(path, w)
}
