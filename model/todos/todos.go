package todos

import (
	"fmt"
	"sort"
	"time"
)

type Todo struct {
	Id        int    `json:"id"`
	Task      string `json:"task"`
	DateAdded string `json:"date_added"`
	Completed string `json:"status"`
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

func (td *Todos) Add(id int, task string) Todo {
	ct := time.Now()
	t := Todo{
		Id:        id,
		Task:      task,
		DateAdded: ct.Format("2006-01-02"),
		Completed: "Incomplete",
		Priority:  0,
	}

	return t
}

func (todos *Todos) List(sortBy TodoSort, showCompleted bool) {
	switch sortBy {
	case SortPriority:
		sort.Stable(sort.Reverse(TodosByPriority(*todos)))
	default:
		sort.Stable(*todos)
	}

	fmt.Printf("Index\tStatus\t\tDate Added\tPriority\tTask\n")
	for idx, todo := range *todos {
		if !showCompleted && todo.Completed == "Completed" {
			continue
		}
		fmt.Printf("%v\t%s\t%s\t%v\t\t%s\n", idx+1, todo.Completed, todo.DateAdded, todo.Priority, todo.Task)
	}
}

func (todos *Todos) Update(idx int, task string) error {
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

func (todos *Todos) Delete(idx int) error {
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

func (todos *Todos) Complete(idx int) error {
	ai := idx - 1

	if ai < 0 || ai >= len(*todos) {
		return fmt.Errorf("The given index is out of bounds.")
	}

	for _, todo := range *todos {
		if todo.Id == (*todos)[ai].Id {
			(*todos)[ai].Completed = "Completed"
		}
	}

	return nil
}

func (todos *Todos) SetPriority(idx int, priority int) error {
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

func (todos *Todos) SortList(sortParam string, showCompleted bool) error {
	if sortParam == "id" {
		todos.List(SortId, showCompleted)
	} else if sortParam == "priority" {
		todos.List(SortPriority, showCompleted)
	} else {
		return fmt.Errorf("Invalid Sort Column")
	}

	return nil
}
