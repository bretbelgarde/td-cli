package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	td "bretbelgarde.com/td-cli/model/todos"
)

func main() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSort := addCmd.String("sort", "id", "Sort list by <column name>")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listSort := listCmd.String("sort", "id", "Sort list by <column name>")
	listComplete := listCmd.Bool("completed", false, "toggle display of Completed")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateSort := updateCmd.String("sort", "id", "Sort list by <column name>")

	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	delSort := delCmd.String("sort", "id", "Sort list by <column name>")

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeSort := completeCmd.String("sort", "id", "Sort list by <column name>")

	priorityCmd := flag.NewFlagSet("priority", flag.ExitOnError)
	prioritySort := priorityCmd.String("sort", "id", "Sort list by <column name>")

	var todos td.Todos

	if len(os.Args) < 2 {
		fmt.Println("Expected one of the following: 'add', 'list', 'delete', 'update', or 'complete'")
		os.Exit(1)
	}

	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	appDir := userDir + "/td-cli"
	appData := "todo.json"
	appPath := appDir + "/" + appData

	if !pathExists(appDir) {
		if err = os.Mkdir(appDir, 0755); err != nil {
			fmt.Println("Directiory creation error: ", err)
			os.Exit(1)
		}
	}

	if !pathExists(appPath) {
		// If the file doesn't exist save an empty todo file
		if err = save(&todos, appPath); err != nil {
			fmt.Println("File creation error: ", err)
		}
	}

	if err = load(&todos, appPath); err != nil {
		fmt.Println("There was an error loading the file: ", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		id := todos[todos.Len()-1].Id + 1
		task := strings.Join(addCmd.Args(), " ")
		todos = append(todos, todos.Add(id, task))
		err := save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Added!\n\n")
		if err := todos.SortList(*addSort, false); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Printf("\nTodo List:\n\n")
		if err := todos.SortList(*listSort, *listComplete); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "update":
		updateCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(updateCmd.Args()[0])
		task := updateCmd.Args()[1]

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = todos.Update(arg, task)

		if err != nil {
			fmt.Printf("Error while updating todo at index: %v. Error: %s\n", arg, err)
			os.Exit(1)
		}

		err = save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Updated!\n\n")

		if err := todos.SortList(*updateSort, false); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "delete":
		delCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(delCmd.Args()[0])

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = todos.Delete(arg)

		if err != nil {
			fmt.Printf("Error while deleteing todo at index: %v. Error: %s\n", arg, err)
			os.Exit(1)
		}

		err = save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Deleted!\n\n")

		if err := todos.SortList(*delSort, false); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "complete":
		completeCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(completeCmd.Args()[0])

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = todos.Complete(arg)

		if err != nil {
			fmt.Printf("Error while completing todo at index: %v. Error: %s\n", arg, err)
			os.Exit(1)
		}

		err = save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Completed!\n\n")
		if err := todos.SortList(*completeSort, true); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "priority":
		priorityCmd.Parse(os.Args[2:])
		idx, err := strconv.Atoi(priorityCmd.Args()[0])

		if err != nil {
			fmt.Printf("Unable to parse index parameter: %s\n", err)
			os.Exit(1)
		}

		priority, err := strconv.Atoi(priorityCmd.Args()[1])

		if err != nil {
			fmt.Printf("Unable to parse priority paramter: %s\n", err)
			os.Exit(1)
		}

		err = todos.SetPriority(idx, priority)

		if err != nil {
			fmt.Printf("Error while setting the priority (value: %v) of the todo at index: %v. Error: %s\n", priority, idx, err)
			os.Exit(1)
		}

		err = save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Priority Set!\n")
		if err := todos.SortList(*prioritySort, false); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

}

func save(todos *td.Todos, appPath string) error {
	todoJson, err := json.Marshal(todos)

	if err != nil {
		return err
	}

	err = os.WriteFile(appPath, todoJson, 0644)

	if err != nil {
		return err
	}

	return nil
}

func load(todos *td.Todos, filePath string) error {
	file, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(file), todos)

	if err != nil {
		return err
	}

	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
