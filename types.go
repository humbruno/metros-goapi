package main

import (
	"io"
	"net/http"
	"time"
)

type loggingRoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
}

type app struct {
	client *http.Client
	cfg    *cfg
}

type cfg struct {
	port     string
	apiKey   string
	tokenUrl string
	apiUrl   string
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type authRoundTripper struct {
	next       http.RoundTripper
	token      string
	maxRetries int
	retryDelay time.Duration
	cfg        *cfg
}

type apiResponse struct {
	// resposta is an array of empty interfaces
	Resposta []interface{} `json:"resposta"`
}
