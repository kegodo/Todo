// File: todo/cmd/api/handlers.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"todo.kegodo.net/internal/data"
	"todo.kegodo.net/internal/validator"
)

func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	//Our target decode destination
	var input struct {
		Title       string `json:"title"`
		Descritpion string `json:"description"`
		Done        bool   `json:"status"`
	}

	//Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//coping the valeus from the input struct to the new todo struct
	todo := &data.Todo{
		Title:       input.Title,
		Descritpion: input.Descritpion,
		Done:        input.Done,
	}

	//Initialize a new Validator Instance
	v := validator.New()

	//check the map to determine if ther were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Creating a todo element
	err = app.models.Todos.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//Create a location header for the newly created resource/School
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/toto/%d", todo.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The showentry handler will display an individual todo element
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	//getting the request data from param function
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//Fetching the specific todo element
	todo, err := app.models.Todos.Get(id)

	//Handling errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Writing the data from the returned get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Facilitates an update action to the todo element in the database
func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//Fetch the original record from the database
	todo, err := app.models.Todos.Get(id)

	//Handling the errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Creating an input struct to hold data read in from the client
	//Updating the input struct to use pointers because pointers have a default value of nil
	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Done        *bool   `json:"status"`
	}

	//Initilizing a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//checking for any updates
	if input.Title != nil {
		todo.Title = *input.Title
	}
	if input.Description != nil {
		todo.Descritpion = *input.Description
	}
	if input.Done != nil {
		todo.Done = *input.Done
	}

	//Initilize a new Validator Instance
	v := validator.New()

	//Checking the map to determin if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Passing the updated todo element to the update() method
	err = app.models.Todos.Update(todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Writing the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// To facilitate deletion of a todo element
func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	err = app.models.Todos.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Returning 200 status ok to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "todo element sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listTodo handler allows the client to see a listing of a schools based on a set of criteria
func (app *application) listTododHandler(w http.ResponseWriter, r *http.Request) {
	//creating an input struct to hold our query parameters
	var input struct {
		Title       string
		Description string
		Done        bool
		data.Filters
	}

	//Initializing a validator
	v := validator.New()

	//getting the URL values map
	qs := r.URL.Query()

	//Using the helper method to extract the values
	input.Title = app.readString(qs, "title", "")
	input.Description = app.readString(qs, "decription", "")
	input.Done = app.readBool(qs, "status", false, v)

	//filering now
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortList = []string{"id", "title", "status", "-id", "-description", "-status"}

	//checking for validation errors
	if data.ValidateFilter(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Geting a listing of all todo elements
	todos, metadata, err := app.models.Todos.GetAll(input.Title, input.Description, input.Done, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//sending JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{"todos": todos, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
