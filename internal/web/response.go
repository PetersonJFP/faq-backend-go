package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ResponseError struct {
	Error string `json:"error"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Esta função ensina o validador a procurar a tag "label" na struct.
	// Se não encontrar "label", ele tenta o "json". Se não tiver, usa o nome do campo.
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
		// e.Field() agora retornará o valor da tag "label" ou "json"
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

	if len(messages) > 0 {
		return messages[0]
	}
	return "Erro de validação"
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ResponseError{Error: message})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func ReadJSON(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("JSON inválido: verifique a sintaxe")
	}

	err := validate.Struct(data)
	if err != nil {
		return fmt.Errorf("%s", formatValidationError(err))
	}

	return nil
}
