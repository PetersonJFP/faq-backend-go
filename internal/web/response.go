package web

import (
	"encoding/json"
	"net/http"
)

// ResponseError define a estrutura padrão para todos os erros da API
type ResponseError struct {
	Error string `json:"error"`
}

// JSON envia uma resposta em formato JSON com o status code especificado
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error envia uma resposta de erro padronizada
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ResponseError{Error: message})
}

// NoContent envia uma resposta de sucesso sem corpo (comum em DELETE)
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
