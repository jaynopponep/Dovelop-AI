package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"io"
	"log"
	"os"
	"strings"
)

type Task struct {
	Num  int
	Task string
}

func buildTable() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Task #"},
			{Align: simpletable.AlignCenter, Text: "Task"},
		},
	}
	var taskData = getData()
	for _, row := range taskData {
		taskNum := int(row[0].(float64))
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", taskNum)},
			{Text: row[1].(string)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{Align: simpletable.AlignCenter, Text: "ChatGPT suggests breaking down these tasks into smaller ones! Accept? (Y/N)"},
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
		task = strings.TrimSpace(task)
		// remove newline from input aka give the string a haircut from leading and trailing white space
		t := Task{int(num), task}
		createTask(t)

	case 2:
		fmt.Println("You chose to mark a task done")
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
	}
}

func createTask(t Task) {
	var tasks = getData()

	// create new task from t Task
	newTask := []interface{}{t.Num, t.Task}
	// append t Task into tasks
	tasks = append(tasks, newTask)

	// marshal tasks into json
	newData, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal("Error marshalling new data:", err)
	}

	if err := os.WriteFile("data.json", newData, 0644); err != nil {
		log.Fatal("Error writing new data to file:", err)
	}

	log.Println("Data added")
	buildTable()
}

func getData() [][]interface{} {
	// open file, checking if exists (read write, create permissions)
	file, err := os.OpenFile("data.json", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Error closing file:", err)
		}
	}(file)

	// retrieve tasks in data.json
	var tasks [][]interface{}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) > 0 {
		// unmarshal if content is available (data -> tasks interface)
		if err := json.Unmarshal(data, &tasks); err != nil {
			log.Fatal("Error unmarshalling data:", err)
		}
	}
	return tasks
}

func deleteTask(num int) {
	tasks := getData()
	var del int
	found := false
	for i, task := range tasks {
		if int(task[0].(float64)) == num {
			del = i
			found = true
			break
		}
	}
	if found {
		tasks = append(tasks[:del], tasks[del+1:]...)

		updatedData, err := json.Marshal(tasks)
		if err != nil {
			log.Fatal("Error marshalling data:", err)
		}
		if err := os.WriteFile("data.json", updatedData, 0644); err != nil {
			log.Fatal("Error writing new data to file:", err)
		}
		fmt.Printf("Deleted task #%d\n", num)
	} else {
		fmt.Printf("Task #%d not found\n", num)
	}
}

func main() {
	buildTable()
	// begin main menu process
	fmt.Println("Choose an option" + "\n" + "1. Create a new task" + "\n" + "2. Mark a task done" + "\n" + "3. Delete a task" + "\n" + "4. Edit a task" + "\n" + "5. Exit")
	var choice int
	_, inputErr := fmt.Scanln(&choice)
	if inputErr != nil {
		fmt.Println("Invalid input, try again")
		return
	} else {
		processChoice(choice)
	}
}
