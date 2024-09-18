package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSON[T any](body io.ReadCloser) (*T, error) {

	model := new(T)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	defer body.Close()

	err := decoder.Decode(model)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func JSON[T any](w http.ResponseWriter, status int, data T) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
