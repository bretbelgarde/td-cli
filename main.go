package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	td "bretbelgarde.com/td-cli/model/todos"
	"bretbelgarde.com/td-cli/utils"
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

	if !utils.PathExists(appDir) {
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

		fmt.Printf("\nTodo added.\n\n")

		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	case "list":
		listCmd.Parse(os.Args[2:])
		var sortBy string

		switch *listSort {
		case "due":
			sortBy = td.SortDueDate
		case "priority":
			sortBy = td.SortPriority
		default:
			sortBy = td.SortDefault
		}

		if *listComplete {
			utils.FormatCompleted(utils.GetTodoList(*tdb, sortBy, true))
		} else {
			utils.FormatOutput(utils.GetTodoList(*tdb, sortBy, false))
		}

	case "update":
		updateCmd.Parse(os.Args[2:])
		id := utils.ParseValue(updateCmd.Arg(0))

		if _, err = tdb.Update(id, "task", updateCmd.Arg(1)); err != nil {
			fmt.Println("Error updating todo: ", err)
			os.Exit(1)
		}

		fmt.Printf("\nTodo updated.\n\n")

		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	case "delete":
		delCmd.Parse(os.Args[2:])
		id := utils.ParseValue(delCmd.Arg(0))

		if _, err = tdb.Delete(id); err != nil {
			fmt.Println("Error deleting todo: ", err)
			os.Exit(1)
		}

		fmt.Println("Todo deleted.")

		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	case "complete":
		completeCmd.Parse(os.Args[2:])
		id := utils.ParseValue(completeCmd.Arg(0))

		if err = tdb.Complete(id); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nTodo completed.\n\n")

		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	case "priority":
		priorityCmd.Parse(os.Args[2:])
		id := utils.ParseValue(priorityCmd.Arg(0))
		tp := utils.ParseValue(priorityCmd.Arg(1))

		if err := tdb.SetPriority(id, tp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nTodo priority updated.\n\n")
		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	case "due":
		dueCmd.Parse(os.Args[2:])
		id := utils.ParseValue(updateCmd.Arg(0))
		date := dueCmd.Arg(1)

		if err := tdb.SetDueDate(id, date); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nTodo due date updated.\n\n")
		utils.FormatOutput(utils.GetTodoList(*tdb, td.SortDefault, false))

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

	os.Exit(0)
}
