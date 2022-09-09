package htp

import (
	"net/http"

	"github.com/go-chi/render"
)

func Render400(w http.ResponseWriter, r *http.Request, errors []Error, example interface{}) {
	render.Status(r, http.StatusBadRequest)

	render.JSON(w, r, ValidationErrors{
		Errors:  errors,
		Example: example,
	})
}
