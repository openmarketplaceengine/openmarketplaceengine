package htp

import (
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/validate"
	"net/http"

	"github.com/go-chi/render"
)

func Render500(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ValidationErrors{
		Errors: []validate.Error{
			{
				Message: "something went wrong",
				Details: fmt.Sprintf("%s", err),
			},
		}})
}

func Render500ve(w http.ResponseWriter, r *http.Request, e validate.Error) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ValidationErrors{
		Errors: []validate.Error{
			e,
		}})
}
