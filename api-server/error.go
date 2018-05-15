package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pborman/uuid"
)

type ErrorSource struct {
	Pointer   string `json:"pointer"`
	Parameter string `json:"parameter"`
}

type ApiError struct {
	Id     string      `json:"id"`
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Source ErrorSource `json:"source"`
}

func NewApiError(status string, code int, title string, err error, r *http.Request) ApiError {
	apiError := ApiError{
		Id:     uuid.New(),
		Status: status,
		Code:   code,
		Title:  title,
	}

	if err != nil {
		apiError.Detail = err.Error()
	}

	if r != nil {
		apiError.Source = ErrorSource{
			Pointer:   fmt.Sprintf("(%s) %s", r.Method, r.URL.String()),
			Parameter: fmt.Sprintf("%+v", r.Form),
		}
	}

	return apiError
}

func (a *ApiError) toString() string {
	bin, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a)
	}
	return string(bin)
}
