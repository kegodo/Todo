// File: todo/internal/data/todo.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todo.kegodo.net/internal/validator"
)

// Todo struct supports the infromation for the todo todo
type Todo struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Descritpion string    `json:"description"`
	Done        bool      `json:"status`
	Version     int32     `json:"version"`
}

func Validatetodo(v *validator.Validator, todo *Todo) {
	//using check() method to check our validation checks
	v.Check(todo.Title != "", "title", "must be provided")
	v.Check(len(todo.Title) <= 250, "title", "must not be more than 250 bytes long")

	v.Check(todo.Descritpion != "", "description", "must be provided")
	v.Check(len(todo.Descritpion) <= 250, "description", "must no be more than 250 bytes long")

	v.Check(todo.Done != false, "Done", "default task as false")
}

type TodoModel struct {
	DB *sql.DB
}

// Insert() allows us to create a new todo
func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (title, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`

	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	//collect the date field into a slice
	args := []interface{}{todo.Title, todo.Descritpion}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

// Get() allows us to retrieve a specific task
func (m TodoModel) Get(id int64) (*Todo, error) {
	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//Construct our query with the given id
	query := `
		SELECT id, create_at, title, description, status, version
		FROM todos
		WHERE id = $1
	`

	//Declaring the Todo varaible to hold the returned data
	var todo Todo

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Title,
		&todo.Descritpion,
		&todo.Done,
		&todo.Version,
	)

	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//Succes
	return &todo, nil
}

// Update() allows us to edit/alter a specific todo task
// Optimistic locking (version number)
func (m TodoModel) Update(todo *Todo) error {
	//create a query
	query := `
		UPDATE todos
		SET title = $1, description = $2, status = $3, version = version + 1
		WHERE id = $4
		AND version = $5
		RETURNING version
	`
	args := []interface{}{todo.Title, todo.Descritpion, todo.Done, todo.ID, todo.Version}

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	//Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (m TodoModel) Delete(id int64) error {
	//Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	//creating the delete query
	query := `
		DELETE FROM todos
		WHERE id = $1
	`

	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//clearing up to prevent memory leaks
	defer cancel()

	//Executing the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//checking how many rows were affected by the delete operation. we call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m TodoModel) GetAll(title string, description string, status bool, filters Filters) ([]*Todo, Metadata, error) {
	//constructing the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
		id, created_at, title, description, status, version
		FROM todos
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (status = FALSE OR status = TRUE)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
	`, filters.sortColumn(), filters.sortOrder())

	//creating the 3 second time out context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Execute the query
	args := []interface{}{title, description, status, filters.limit(), filters.offSet()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	//Closing the result set
	defer rows.Close()
	totalRecords := 0

	//Initialize an empty slice to hold the task data
	tasks := []*Todo{}

	//Iterate over the rows in the result set
	for rows.Next() {
		var todo Todo

		//Scanning the valus from the row into the todo struct
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Descritpion,
			&todo.Done,
			&todo.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//Add the todo to our slice
		tasks = append(tasks, &todo)
	}
	//checking for errors after looping through the result set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	//returning the slice of todos
	return tasks, metadata, nil
}
