package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Id        int    `json:"id"`
	Task      string `json:"task"`
	DateAdded string `json:"date_added"`
	Completed bool   `json:"status"`
	Priority  int    `json:"priority"`
}

type TodoSort byte

const (
	SortId TodoSort = iota
	SortPriority
)

type Todos []Todo

func (t Todos) Len() int {
	return len(t)
}

func (t Todos) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Todos) Less(i, j int) bool {
	return t[i].Id < t[j].Id
}

type TodosByPriority Todos

func (t TodosByPriority) Len() int {
	return len(t)
}

func (t TodosByPriority) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TodosByPriority) Less(i, j int) bool {
	return t[i].Priority < t[j].Priority
}

func main() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSort := addCmd.String("sort", "id", "Sort list by <column name>")
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listSort := listCmd.String("sort", "id", "Sort list by <column name>")
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateSort := updateCmd.String("sort", "id", "Sort list by <column name>")
	delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	delSort := delCmd.String("sort", "id", "Sort list by <column name>")
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	completeSort := completeCmd.String("sort", "id", "Sort list by <column name>")
	priorityCmd := flag.NewFlagSet("priority", flag.ExitOnError)
	prioritySort := priorityCmd.String("sort", "id", "Sort list by <column name>")

	var todos Todos

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
		todos = append(todos, addTodo(id, task))
		err := save(&todos, appPath)

		if err != nil {
			fmt.Printf("There was an error saving the file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nTask Added!\n\n")
		if err := sortList(&todos, *addSort); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Printf("\nTodo List:\n\n")
		if err := sortList(&todos, *listSort); err != nil {
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

		if err := sortList(&todos, *updateSort); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

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

		if err := sortList(&todos, *delSort); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

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
		if err := sortList(&todos, *completeSort); err != nil {
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

		err = setPriority(&todos, idx, priority)

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
		if err := sortList(&todos, *prioritySort); err != nil {
			fmt.Printf("Sorting Error: %s\n", err)
		}

	default:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	}

}

func addTodo(id int, task string) Todo {
	ct := time.Now()
	t := Todo{
		Id:        id,
		Task:      task,
		DateAdded: ct.Format("2006-01-02"),
		Completed: false,
		Priority:  0,
	}

	return t
}

func listTodos(todos *Todos, sortBy TodoSort) {
	switch sortBy {
	case SortPriority:
		sort.Stable(sort.Reverse(TodosByPriority(*todos)))
	default:
		sort.Stable(*todos)
	}

	fmt.Printf("Index\tStatus\t\tDate Added\tPriority\tTask\n")
	for idx, todo := range *todos {
		var status string
		if todo.Completed {
			status = "Completed"
		} else {
			status = "Incomplete"
		}

		fmt.Printf("%v\t%s\t%s\t%v\t\t%s\n", idx+1, status, todo.DateAdded, todo.Priority, todo.Task)
	}
}

func updateTodo(todos *Todos, idx int, task string) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	for _, todo := range *todos {
		if todo.Id == (*todos)[ai].Id {
			(*todos)[ai].Task = task
		}
	}

	return nil
}

func deleteTodo(todos *Todos, idx int) error {
	ai := idx - 1

	var tmpTodos Todos

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	for _, todo := range *todos {
		if todo.Id != (*todos)[ai].Id {
			tmpTodos = append(tmpTodos, todo)
		}
	}

	(*todos) = tmpTodos

	return nil
}

func completeTodo(todos *Todos, idx int) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	for _, todo := range *todos {
		if todo.Id == (*todos)[ai].Id {
			(*todos)[ai].Completed = true
		}
	}

	return nil
}

func setPriority(todos *Todos, idx int, priority int) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	for _, todo := range *todos {
		if todo.Id == (*todos)[ai].Id {
			(*todos)[ai].Priority = priority
		}
	}

	return nil
}

func save(todos *Todos, appPath string) error {
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

func load(todos *Todos, filePath string) error {
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

func sortList(todos *Todos, sortParam string) error {
	if sortParam == "id" {
		listTodos(todos, SortId)
	} else if sortParam == "priority" {
		listTodos(todos, SortPriority)
	} else {
		return fmt.Errorf("Invalid Sort Column")
	}

	return nil
}
