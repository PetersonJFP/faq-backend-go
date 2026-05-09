package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ResponseError define o contrato único de erro da API
type ResponseError struct {
	Error string `json:"error"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		}
		if name == "-" {
			return ""
		}
		return name
	})
}

func formatValidationError(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Dados inválidos"
	}

	var messages []string
	for _, e := range errs {
		field := e.Field()

		switch e.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("O campo '%s' é obrigatório", field))
		case "email":
			messages = append(messages, fmt.Sprintf("O campo '%s' deve ser um e-mail válido", field))
		case "min":
			messages = append(messages, fmt.Sprintf("O campo '%s' deve ter no mínimo %s caracteres", field, e.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("O campo '%s' deve ter no máximo %s caracteres", field, e.Param()))
		case "numeric":
			messages = append(messages, fmt.Sprintf("O campo '%s' deve conter apenas números", field))
		default:
			messages = append(messages, fmt.Sprintf("O campo '%s' é inválido", field))
		}
	}

	// Retornamos apenas o primeiro erro para manter a resposta Flat e limpa
	if len(messages) > 0 {
		return messages[0]
	}
	return "Erro de validação"
}

// JSON envia uma resposta de sucesso
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

// NoContent envia uma resposta vazia (comum em Deletes)
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ReadJSON decodifica o corpo da requisição e valida as regras da struct
func ReadJSON(r *http.Request, data interface{}) error {
	// 1. Tenta decodificar o JSON
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("JSON inválido: verifique a sintaxe")
	}

	// 2. Executa a validação baseada nas tags da struct
	err := validate.Struct(data)
	if err != nil {
		return fmt.Errorf("%s", formatValidationError(err))
	}

	return nil
}
