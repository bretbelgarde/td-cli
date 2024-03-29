package todos

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const create string = `
	CREATE TABLE IF NOT EXISTS todos (
	id INTEGER NOT NULL PRIMARY KEY,
	task TEXT NOT NULL,
	date_added DATETIME NOT NULL,
	date_due DATETIME,
	date_completed DATETIME,
	completed INTEGER NOT NULL DEFAULT 0,
	priority  INTEGER NOT NULL DEFAULT 0
	);`

const (
	SortDefault       = "id ASC"
	SortDueDate       = "date_due DESC"
	SortDateCompleted = "date_completed DESC"
	SortPriority      = "priority DESC"
)

type Todo struct {
	Id            int
	Task          string
	DateAdded     string
	DateDue       sql.NullString
	DateCompleted sql.NullString
	Completed     int
	Priority      int
}

type Todos struct {
	db *sql.DB
}

func NewTodos(dbpath string) (*Todos, error) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &Todos{
		db: db,
	}, nil
}

func (td *Todos) Insert(todo Todo) (int, error) {
	var id int64
	res, err := td.db.Exec(
		"INSERT INTO todos VALUES(null, ?, ?, ?, ?, ?, ?);",
		&todo.Task,
		&todo.DateAdded,
		&todo.DateDue,
		&todo.DateCompleted,
		&todo.Completed,
		&todo.Priority)

	if err != nil {
		return 0, err
	}

	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (td *Todos) Retrieve(id int) (Todo, error) {
	row := td.db.QueryRow("SELECT * FROM todos WHERE id=?", id)

	var err error
	todo := Todo{}

	if err = row.Scan(
		&todo.Id,
		&todo.Task,
		&todo.DateAdded,
		&todo.DateDue,
		&todo.DateCompleted,
		&todo.Completed,
		&todo.Priority); err == sql.ErrNoRows {
		return Todo{}, err
	}
	return todo, err
}

func (td *Todos) List(offset int, sortBy string) ([]Todo, error) {

	rows, err := td.db.Query("SELECT * FROM todos WHERE ID > ? AND completed <> 1 ORDER BY "+sortBy+" LIMIT 100;", offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		todo := Todo{}

		if err = rows.Scan(
			&todo.Id,
			&todo.Task,
			&todo.DateAdded,
			&todo.DateDue,
			&todo.DateCompleted,
			&todo.Completed,
			&todo.Priority); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (td *Todos) ListCompleted(offset int) ([]Todo, error) {
	rows, err := td.db.Query("SELECT * FROM todos WHERE completed = 1 ORDER BY date_completed DESC;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		todo := Todo{}

		if err = rows.Scan(
			&todo.Id,
			&todo.Task,
			&todo.DateAdded,
			&todo.DateDue,
			&todo.DateCompleted,
			&todo.Completed,
			&todo.Priority,
		); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (td *Todos) Update(id int64, field string, value string) (int, error) {
	var ra int64
	sql := "UPDATE todos SET " + field + "=? WHERE id=?;"

	res, err := td.db.Exec(sql, value, id)

	if err != nil {
		return 0, err
	}

	if ra, err = res.RowsAffected(); err != nil {
		return 0, err
	}

	return int(ra), err
}

func (td *Todos) Delete(id int64) (int, error) {
	res, err := td.db.Exec("DELETE FROM todos WHERE id = ?", id)

	if err != nil {
		return 0, err
	}

	if id, err = res.RowsAffected(); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (td *Todos) Complete(id int64) error {
	if _, err := td.Update(id, "completed", "1"); err != nil {
		return err
	}

	if _, err := td.Update(id, "date_completed", time.Now().Format(time.DateOnly)); err != nil {
		return err
	}

	return nil
}

func (td *Todos) SetPriority(id int64, priority int64) error {
	_, err := td.Update(id, "priority", strconv.FormatInt(priority, 10))

	if err != nil {
		return err
	}

	return nil
}

func (td *Todos) SetDueDate(id int64, dueDate string) error {
	date, err := time.Parse("01-02-2006", dueDate)

	if err != nil {
		return err
	}

	if _, err := td.Update(id, "date_due", date.Format(time.DateOnly)); err != nil {
		return err
	}

	return nil
}
