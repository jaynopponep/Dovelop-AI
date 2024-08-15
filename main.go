package main

import (
	"bufio"
	"fmt"
	"github.com/alexeyco/simpletable"
	"log"
	"os"
	"strings"
)

type Task struct {
	Num    int
	Task   string
	Status string
}

func buildTable() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Task #"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Status"},
		},
	}
	var taskData = getData()
	for _, row := range taskData {
		taskNum := int(row[0].(float64))
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", taskNum)},
			{Text: row[1].(string)},
			{Text: row[2].(string)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignCenter, Text: "ChatGPT suggests breaking down these tasks into smaller ones! Accept? (Y/N)"},
			{},
		},
	}
	table.SetStyle(simpletable.StyleUnicode)
	fmt.Println(table.String())
}

func processChoice(choice int) {
	switch choice {
	case 1:
		fmt.Println("You chose to create a new task")
		var tasks = getData()
		var num float64
		if len(tasks) <= 0 {
			num = 1
		} else {
			var last = len(tasks) - 1
			num = float64(int(tasks[last][0].(float64))) + 1
		}
		var task string
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your task: ")
		task, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		var status = "Incomplete"
		task = strings.TrimSpace(task)
		// remove newline from input aka give the string a haircut from leading and trailing white space
		t := Task{int(num), task, status}
		createTask(t)
		promptUser()
	case 2:
		fmt.Println("You chose to mark a task done")
		fmt.Print("Enter the task number you want to mark done: ")
		var choice int
		_, inputErr := fmt.Scanln(&choice)
		if inputErr != nil {
			fmt.Println("Invalid input, try again")
			return
		} else {
			markDone(choice)
		}
		buildTable()
		promptUser()
	case 3:
		fmt.Println("You chose to delete a task")
		fmt.Print("Enter the task number you want to delete: ")
		var choice int
		_, inputErr := fmt.Scanln(&choice)
		if inputErr != nil {
			fmt.Println("Invalid input, try again")
			return
		} else {
			deleteTask(choice)
		}
		buildTable()
		promptUser()
	case 4:
		fmt.Println("You chose to edit a task")
		fmt.Print("Enter the task number you want to edit: ")
		var choice int
		_, inputErr := fmt.Scanln(&choice)
		if inputErr != nil {
			fmt.Println("Invalid input, try again")
			return
		} else {
			editTask(choice)
		}
	case 5:
		fmt.Println("You chose to breakdown a task into smaller ones!")
		fmt.Print("Enter the task number you want to breakdown: ")
		var choice int
		_, inputErr := fmt.Scanln(&choice)
		if inputErr != nil {
			fmt.Println("Invalid input, try again")
			return
		} else {
			breakDownTask(choice)
		}
	case 6:
		break
	}
}

func promptUser() {
	// begin main menu process
	fmt.Println("Choose an option" + "\n" + "1. Create a new task" + "\n" + "2. Mark a task done" + "\n" + "3. Delete a task" + "\n" + "4. Edit a task" + "\n" + "5. Breakdown a task with AI" + "\n" + "6. Exit")
	var choice int
	_, inputErr := fmt.Scanln(&choice)
	if inputErr != nil {
		fmt.Println("Invalid input, try again")
		return
	} else {
		processChoice(choice)
	}
}

func main() {
	buildTable()
	promptUser()
}
