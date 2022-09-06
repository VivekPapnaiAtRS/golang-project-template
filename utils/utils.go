package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DecodeToJson(data interface{}) ([]byte, error) {
	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonData, jsonErr
	}
	return jsonData, nil
}

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
}
