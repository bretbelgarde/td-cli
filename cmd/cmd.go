package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	td "bretbelgarde.com/td-cli/model/todos"
	"bretbelgarde.com/td-cli/utils"
)

var (
	tdb          *td.Todos
	addCmd       *flag.FlagSet = flag.NewFlagSet("add", flag.ExitOnError)
	listCmd      *flag.FlagSet = flag.NewFlagSet("list", flag.ExitOnError)
	listSort     *string       = listCmd.String("sort", "id", "Sort list by <column name>")
	listComplete *bool         = listCmd.Bool("completed", false, "List completed todos")
	updateCmd    *flag.FlagSet = flag.NewFlagSet("update", flag.ExitOnError)
	delCmd       *flag.FlagSet = flag.NewFlagSet("delete", flag.ExitOnError)
	completeCmd  *flag.FlagSet = flag.NewFlagSet("complete", flag.ExitOnError)
	priorityCmd  *flag.FlagSet = flag.NewFlagSet("priority", flag.ExitOnError)
	dueCmd       *flag.FlagSet = flag.NewFlagSet("due", flag.ExitOnError)
)

type Cmder interface {
	Cmd(args []string) ([]td.Todo, error)
}

type AddCmd struct {
	command *flag.FlagSet
}

func (a AddCmd) Cmd(args []string) ([]td.Todo, error) {
	a.command.Parse(args)
	task := strings.Join(a.command.Args(), " ")

	todo := td.Todo{
		Task:      task,
		DateAdded: time.Now().Format("2006-01-02"),
		Completed: 0,
		Priority:  0,
	}

	if _, err := tdb.Insert(todo); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo added.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

type ListCmd struct {
	command  *flag.FlagSet
	sort     *string
	complete *bool
}

func (l ListCmd) Cmd(args []string) ([]td.Todo, error) {
	l.command.Parse(args)
	var sortBy string

	switch *l.sort {
	case "due":
		sortBy = td.SortDueDate
	case "priority":
		sortBy = td.SortPriority
	default:
		sortBy = td.SortDefault
	}

	if *l.complete {
		return utils.GetTodoList(*tdb, sortBy, true), nil
	} else {
		return utils.GetTodoList(*tdb, sortBy, false), nil
	}
}

type UpdateCmd struct {
	command *flag.FlagSet
}

func (u UpdateCmd) Cmd(args []string) ([]td.Todo, error) {
	u.command.Parse(args)
	id := utils.ParseValue(u.command.Arg(0))
	value := u.command.Arg(1)

	if _, err := tdb.Update(id, "task", value); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo updated.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

type DelCmd struct {
	command *flag.FlagSet
}

func (d DelCmd) Cmd(args []string) ([]td.Todo, error) {
	d.command.Parse(args)
	id := utils.ParseValue(d.command.Arg(0))

	if _, err := tdb.Delete(id); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo deleted.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

type CompleteCmd struct {
	command *flag.FlagSet
}

func (c CompleteCmd) Cmd(args []string) ([]td.Todo, error) {
	c.command.Parse(args)
	id := utils.ParseValue(c.command.Arg(0))

	if err := tdb.Complete(id); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo completed.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

type PriorityCmd struct {
	command *flag.FlagSet
}

func (p PriorityCmd) Cmd(args []string) ([]td.Todo, error) {
	p.command.Parse(args)
	id := utils.ParseValue(p.command.Arg(0))
	tp := utils.ParseValue(p.command.Arg(1))

	if err := tdb.SetPriority(id, tp); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo priority updated.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

type DueCmd struct {
	command *flag.FlagSet
}

func (d DueCmd) Cmd(args []string) ([]td.Todo, error) {
	d.command.Parse(args)
	id := utils.ParseValue(d.command.Arg(0))
	date := d.command.Arg(1)

	if err := tdb.SetDueDate(id, date); err != nil {
		return nil, err
	}

	fmt.Printf("\nTodo due date updated.\n\n")

	return utils.GetTodoList(*tdb, td.SortDefault, false), nil
}

func Init() {
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

	tdb, err = td.NewTodos(appPath)
	if err != nil {
		fmt.Println("Err: ", err)
	}
}

type CmdMap map[string]Cmder

func todoCommands(c Cmder) {

}
func Execute(args []string) {
	Init()

	cmdMap := CmdMap{
		"add":      AddCmd{command: addCmd},
		"list":     ListCmd{command: listCmd, sort: listSort, complete: listComplete},
		"update":   UpdateCmd{command: updateCmd},
		"delete":   DelCmd{command: delCmd},
		"complete": CompleteCmd{command: completeCmd},
		"priority": PriorityCmd{command: priorityCmd},
		"due":      DueCmd{command: dueCmd},
	}

	cmdName := args[1]
	cmd := cmdMap[cmdName]
	switch cmd {
	case nil:
		fmt.Println("expected one of the following 'add', 'list', 'delete', 'update', 'complete'")
		os.Exit(1)
	default:
		todoList, err := cmd.Cmd(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		utils.FormatOutput(todoList)
	}
}
