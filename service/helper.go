package service

import (
	"encoding/json"
	"github.com/c12s/discovery/model"
	"log"
	"net/http"
	"strings"
)

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}

func sendErrorMessage(w http.ResponseWriter, msg string, status int) {
	body, err := json.Marshal(map[string]string{"message": msg})
	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}

func resp(service, address string) (string, error) {
	data, err := json.Marshal(
		&model.Resp{
			Service: service,
			Address: address,
		},
	)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func read(body []byte) (*model.Resp, error) {
	data := &model.Resp{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}
	return data, nil
}

// key /heartbeat/service|address
func form(data *model.Resp) string {
	prefix := strings.Join([]string{"/heartbeat", data.Service}, "/")
	return strings.Join([]string{prefix, data.Address}, "|")
}
