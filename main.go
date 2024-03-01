package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	td "bretbelgarde.com/td-cli/model/todos"
)

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listSort := listCmd.String("sort", "id", "Sort list by <column name>")
	listComplete := listCmd.Bool("completed", false, "List completed todos")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)

	priorityCmd := flag.NewFlagSet("priority", flag.ExitOnError)

	dueCmd := flag.NewFlagSet("due", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected one of the following: 'add', 'list', 'delete', 'update', or 'complete'")
		os.Exit(1)
	}

	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	appDir := userDir + "/.td-cli"
	appData := "todos.db"
	appPath := appDir + "/" + appData

	if !pathExists(appDir) {
		if err = os.Mkdir(appDir, 0755); err != nil {
			fmt.Println("Directiory creation error: ", err)
			os.Exit(1)
		}
	}

	tdb, err := td.NewTodos(appPath)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		task := strings.Join(addCmd.Args(), " ")

		todo := td.Todo{
			Task:      task,
			DateAdded: time.Now().Format("2006-01-02"),
			Completed: 0,
			Priority:  0,
		}

		if _, err := tdb.Insert(todo); err != nil {
			fmt.Println("Err: ", err)
			os.Exit(1)
		}

		fmt.Println("Todo added")

	case "list":
		listCmd.Parse(os.Args[2:])
		var sortBy string
		var err error
		var todoList []td.Todo

		switch *listSort {
		case "due":
			sortBy = td.SortDueDate
		case "priority":
			sortBy = td.SortPriority
		default:
			sortBy = td.SortDefault
		}

		if *listComplete {
			todoList, err = tdb.ListCompleted(0)
		} else {
			todoList, err = tdb.List(0, sortBy)
		}

		if err != nil {
			fmt.Println("ERR: ", err)
			os.Exit(1)
		}

		if len(todoList) < 1 {
			fmt.Println("No todos in todo list")
			os.Exit(0)
		}

		if *listComplete {
			formatCompleted(todoList)
		} else {
			formatOutput(todoList)
		}

	case "update":
		updateCmd.Parse(os.Args[2:])
		id := parseValue(updateCmd.Arg(0))

		if _, err = tdb.Update(id, "task", updateCmd.Arg(1)); err != nil {
			fmt.Println("Error updating todo: ", err)
			os.Exit(1)
		}

		fmt.Println("Todo updated.")

	case "delete":
		delCmd.Parse(os.Args[2:])
		id := parseValue(delCmd.Arg(0))

		if _, err = tdb.Delete(id); err != nil {
			fmt.Println("Error deleting todo: ", err)
			os.Exit(1)
		}

		fmt.Println("Todo deleted.")

	case "complete":
		completeCmd.Parse(os.Args[2:])
		id := parseValue(completeCmd.Arg(0))

		if err = tdb.Complete(id); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Todo completed.")

	case "priority":
		priorityCmd.Parse(os.Args[2:])
		id := parseValue(priorityCmd.Arg(0))
		tp := parseValue(priorityCmd.Arg(1))

		if err := tdb.SetPriority(id, tp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "due":
		dueCmd.Parse(os.Args[2:])
		id := parseValue(dueCmd.Arg(0))
		date := dueCmd.Arg(1)

		if err := tdb.SetDueDate(id, date); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

	os.Exit(0)
}

func formatOutput(todoList []td.Todo) {
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

func formatCompleted(todoList []td.Todo) {
	fmt.Printf("ID\tCompleted\tTask\n")
	for _, todo := range todoList {
		completed, err := time.Parse("2006-01-02T00:00:00Z", todo.DateCompleted.String)
		if err != nil {
			fmt.Printf("Time parse error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%v\t%s\t%s", todo.Id, completed.Format("01-02-2006"), todo.Task)
	}
}

func parseValue(val string) int64 {
	conv, err := strconv.Atoi(val)

	if err != nil {
		fmt.Printf("Index parse error: %v\n", err)
		os.Exit(1)
	}

	return int64(conv)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
