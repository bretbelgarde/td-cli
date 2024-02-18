package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Id        int64  `json:"id"`
	Task      string `json:"task"`
	DateAdded string `json:"date_added"`
	Completed bool   `json:"status"`
	Priority  int    `json:"priority"`
}

func main() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)

	var todos []Todo

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
		task := strings.Join(addCmd.Args(), "")
		todos = append(todos, addTodo(task))
		err := save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Added!\n\n")
		listTodos(&todos)

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Printf("\nTodo List:\n\n")
		listTodos(&todos)

	case "update":
		updateCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(updateCmd.Args()[0])
		task := updateCmd.Args()[1]

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = updateTodo(&todos, arg, task)

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
		listTodos(&todos)

	case "delete":
		delCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(delCmd.Args()[0])

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = deleteTodo(&todos, arg)

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
		listTodos(&todos)

	case "complete":
		completeCmd.Parse(os.Args[2:])
		arg, err := strconv.Atoi(completeCmd.Args()[0])

		if err != nil {
			fmt.Printf("Unable to parse parameter:%s\n", err)
			os.Exit(1)
		}

		err = completeTodo(&todos, arg)

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
		listTodos(&todos)

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

}

func addTodo(task string) Todo {
	ct := time.Now()
	t := Todo{
		Id:        ct.Unix(),
		Task:      task,
		DateAdded: ct.Format("2006-01-02"),
		Completed: false,
		Priority:  0,
	}

	return t
}

func listTodos(todos *[]Todo) {
	fmt.Println("ID\tStatus\t\tDate Added\tTask")
	for i, todo := range *todos {
		var status string
		if todo.Completed {
			status = "Completed"
		} else {
			status = "Incomplete"
		}
		fmt.Printf("%v\t%s\t%s\t%s\n", i+1, status, todo.DateAdded, todo.Task)
	}
}

func updateTodo(todos *[]Todo, idx int, task string) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	(*todos)[ai].Task = task

	return nil
}

func deleteTodo(todos *[]Todo, idx int) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	(*todos) = append((*todos)[:ai], (*todos)[ai+1:]...)

	return nil
}

func completeTodo(todos *[]Todo, idx int) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	(*todos)[ai].Completed = true

	return nil
}

func save(todos *[]Todo, appPath string) error {
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

func load(todos *[]Todo, filePath string) error {
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
