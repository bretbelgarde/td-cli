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
	//addSort := addCmd.String("sort", "id", "Sort list by <column name>")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	//listSort := listCmd.String("sort", "id", "Sort list by <column name>")
	//listComplete := listCmd.Bool("completed", false, "toggle display of Completed")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	// updateSort := updateCmd.String("sort", "id", "Sort list by <column name>")

	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// delSort := delCmd.String("sort", "id", "Sort list by <column name>")

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	// completeSort := completeCmd.String("sort", "id", "Sort list by <column name>")

	priorityCmd := flag.NewFlagSet("priority", flag.ExitOnError)
	// prioritySort := priorityCmd.String("sort", "id", "Sort list by <column name>")

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
		todoList, err := tdb.List(0)

		if err != nil {
			fmt.Println("ERR: ", err)
			os.Exit(1)
		}

		if len(todoList) < 1 {
			fmt.Println("No Todos in DB")
			os.Exit(0)
		}

		formatOutput(todoList)

	case "update":
		updateCmd.Parse(os.Args[2:])
		id, err := strconv.Atoi(updateCmd.Args()[0])
		if err != nil {
			fmt.Println("Index parse error: ", err)
			os.Exit(1)
		}

		if _, err = tdb.Update(int64(id), "task", updateCmd.Args()[1]); err != nil {
			fmt.Println("Error updating todo: ", err)
			os.Exit(1)
		}

		fmt.Println("Todo updated.")

	case "delete":
		delCmd.Parse(os.Args[2:])
		id, err := strconv.Atoi(delCmd.Args()[0])

		if err != nil {
			fmt.Println("Index parse error: ", err)
			os.Exit(1)
		}

		if _, err = tdb.Delete(int64(id)); err != nil {
			fmt.Println("Error deleting todo: ", err)
		}

		fmt.Println("Todo deleted.")

	case "complete":
		completeCmd.Parse(os.Args[2:])
		id, err := parseValue(completeCmd.Arg(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = tdb.Complete(id); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Todo completed.")

	case "priority":
		priorityCmd.Parse(os.Args[2:])

		id, err := parseValue(priorityCmd.Arg(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pri, err := parseValue(priorityCmd.Arg(1))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := tdb.SetPriority(id, pri); err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

	os.Exit(0)
}

func formatOutput(todoList []td.Todo) {
	for _, todo := range todoList {
		fmt.Println(todo)
	}
}

func parseValue(val string) (int64, error) {
	conv, err := strconv.Atoi(val)

	if err != nil {
		return 0, fmt.Errorf("Index parse error: %v", err)
	}

	return int64(conv), nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
