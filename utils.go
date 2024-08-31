package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (l loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.logger, "[%s] %s %s\n", time.Now().Format(time.ANSIC), r.Method, r.URL.String())
	return l.next.RoundTrip(r)
}

func (a *authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	clonedReq := r.Clone(r.Context())
	clonedReq.Header.Set("Authorization", "Bearer "+a.token)
	clonedReq.Header.Set("Accept", "application/json")
	clonedReq.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	var err error
	retries := 0

	for retries < a.maxRetries {
		// Perform the request
		resp, err = a.next.RoundTrip(clonedReq)
		if err != nil {
			return nil, err
		}

		// If the response is not 401, break the loop and return the response
		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}

		// If the response is 401 and retries are available, refresh the token and retry
		if retries < a.maxRetries {
			a.refreshToken()
			clonedReq.Header.Set("Authorization", "Bearer "+a.token)

			time.Sleep(a.retryDelay)

			retries++
		} else {
			// Max retries reached, return the last response
			return resp, err
		}
	}

	return resp, err
}

func (a *authRoundTripper) refreshToken() {
	req, err := http.NewRequest(http.MethodPost, a.cfg.tokenUrl, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		log.Fatalf("Error creating request to refresh token: %v", err)
		return
	}

	req.Header = http.Header{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {fmt.Sprintf("Basic %s", a.cfg.apiKey)},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error making request to refresh token: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body in refresh token: %v", err)
		return
	}

	var tokenRes refreshTokenResponse

	err = json.Unmarshal(body, &tokenRes)
	if err != nil {
		log.Fatalf("Error unmarshalling response body for refresh token: %v", err)
		return
	}

	a.token = tokenRes.AccessToken
}

func (app *app) handleMetroApiResponse(path string, w http.ResponseWriter) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", app.cfg.apiUrl, path), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := app.client.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var jsonBody apiResponse
	if err := readJSON(res, &jsonBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, jsonBody.Resposta)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func readJSON(res *http.Response, v interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
