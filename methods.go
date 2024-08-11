package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

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

func createTask(t Task) {
	var tasks = getData()

	// create new task from t Task
	newTask := []interface{}{t.Num, t.Task, t.Status}
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

// Edit task
func markDone(num int) {
	tasks := getData()
	for i, task := range tasks {
		if int(task[0].(float64)) == num {
			tasks[i][2] = "‎✅  ✅  ✅‎‎"
		}
	}
	updatedData, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal("Error marshalling data:", err)
	}
	if err := os.WriteFile("data.json", updatedData, 0644); err != nil {
		log.Fatal("Error writing new data to file:", err)
	}
	fmt.Printf("Marked task #%d done\n", num)
}
