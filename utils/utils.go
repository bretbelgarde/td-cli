package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	td "bretbelgarde.com/td-cli/model/todos"
)

func GetTodoList(tdb td.Todos, sortBy string, completed bool) []td.Todo {
	var todoList []td.Todo
	var err error

	if !completed {
		todoList, err = tdb.List(0, sortBy)
	} else {
		todoList, err = tdb.ListCompleted(0)
	}

	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}

	if len(todoList) < 1 {
		fmt.Println("No todos in todo list")
		os.Exit(0)
	}

	return todoList
}

func FormatOutput(todoList []td.Todo) {
	fmt.Printf("ID\tDue\t\tPri\tTask\n")
	for _, todo := range todoList {
		var dateString string
		if todo.DateDue.String == "" {
			dateString = "-         "
		} else {
			due, err := time.Parse("2006-01-02T00:00:00Z", todo.DateDue.String)
			if err != nil {
				fmt.Printf("Time parse error: %v\n", err)
				os.Exit(1)
			}

			dateString = due.Format("01-02-2006")
		}

		fmt.Printf("%v\t%s\t%v\t%s\n", todo.Id, dateString, todo.Priority, todo.Task)
	}
}

func FormatCompleted(todoList []td.Todo) {
	fmt.Printf("ID\tCompleted\tTask\n")
	for _, todo := range todoList {
		completed, err := time.Parse("2006-01-02T00:00:00Z", todo.DateCompleted.String)
		if err != nil {
			fmt.Printf("Time parse error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%v\t%s\t%s\n", todo.Id, completed.Format("01-02-2006"), todo.Task)
	}
}

func ParseValue(val string) int64 {
	conv, err := strconv.Atoi(val)

	if err != nil {
		fmt.Printf("Index parse error: %v\n", err)
		os.Exit(1)
	}

	return int64(conv)
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
