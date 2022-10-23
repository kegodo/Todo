// File: todo/cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	//security routes
	router.NotFound = http.HandlerFunc(app.notFoundReponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/todo", app.createTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.showTodoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/todo/:id", app.updateTodoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todo/:id", app.deleteTodoHandler)

	return router
}
