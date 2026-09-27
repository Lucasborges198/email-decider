package handler

import (
	"bytes"
	dto "email-decider/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}


func FetchLayaService(w http.ResponseWriter, r *http.Request) {
	var data dto.ClientMessageRequest
	
	loadEnvVar(w)
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	if err != nil {
		http.Error(w, "Error while fetching payload json object", http.StatusInternalServerError)
		return
	}
	
	payload, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error while building json object", http.StatusInternalServerError)
		return
	}
	
	var LayaURL = os.Getenv("LAYA_URL")
	if LayaURL == "" {
		http.Error(
			w,
			"LAYA_URL not found",
			http.StatusInternalServerError,
		)
		return
	}

	url := fmt.Sprintf("%s/laya/decision", LayaURL)

	request, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)

	request.Header.Set("Content-Type", "application/json")


	response, err := client.Do(request)

	if err != nil {
		if response == nil {
			http.Error(
				w,
				"Laya Service didnt answered",
				http.StatusBadGateway,
			)
			return
		}

		http.Error(
			w,
			"Error while calling Laya Service",
			http.StatusBadGateway,
		)
		return
	}



	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		http.Error(
			w,
			"Invalid response",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(body)
}

func loadEnvVar(w http.ResponseWriter) {
	err := godotenv.Load()
	if err != nil {
		http.Error(
			w,
			"Error while starting env var LAYA_URL",
			http.StatusInternalServerError,
		)
	}
}