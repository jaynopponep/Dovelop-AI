package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"os"
	"strings"
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

func editTask(num int) {
	var newTask string
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your new task for task %d: ", num)
	newTask, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}
	newTask = strings.TrimSpace(newTask)
	tasks := getData()
	for i, task := range tasks {
		if int(task[0].(float64)) == num {
			tasks[i][1] = newTask
			break
		}
	}
	updatedData, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal("Error marshalling data:", err)
	}
	if err := os.WriteFile("data.json", updatedData, 0644); err != nil {
		log.Fatal("Error writing new data to file:", err)
	}
	fmt.Printf("Edited task #%d\n", num)
	buildTable()
	promptUser()
}

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

func breakDownTask(num int) {
	tasks := getData()
	var taskName interface{}
	for i, task := range tasks {
		if int(task[0].(float64)) == num {
			taskName = tasks[i][1]
			break
		}
	}
	var taskString = taskName.(string)
	sendToGPT(taskString)
	newTasks, err := txtToArray("tasks.txt")
	if err != nil { log.Fatalf("Error reading tasks from file: %v\n", err) }
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("ChatGPT has suggested your task becomes the three below. Accept? (Y/N)")
	userChoice, err := reader.ReadString('\n')
	if err != nil { log.Fatalf("Error reading user input: %v\n", err) }
	if (strings.TrimSpace(userChoice) == "Y") {
		replaceTasks(newTasks, num)	
	} else {
		os.Exit(0)
	}
}

func sendToGPT(prompt string) {
	var token = ""
	client := openai.NewClient(token)
	var content = "The following is a task that you will break down into THREE tasks only. These three tasks you will store into a .txt file where each line is each task. Only write the tasks and do not numerically label or dot label them at all. Do not respond with anything else except the .txt. Here's the task: " + prompt
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	var res = resp.Choices[0].Message.Content
	
	// writing txt
	err = os.WriteFile("tasks.txt", []byte(res), 0644)
	if err != nil { log.Fatalf("Failed writing task to file: %v\n", err) }
}

func txtToArray(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil { return nil, err }
	defer file.Close()
	var tasks []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tasks = append(tasks, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil { return nil, err }
	return tasks, nil
}

func replaceTasks(newTasks []string, num int) {
// void function that just replaces the task chosen with the new tasks
	var tasks = getData()
	var convertedTasks [][]interface{}
	for i, task := range newTasks {
		convertedTask := []interface{}{num+i, task, "Incomplete"}
		convertedTasks = append(convertedTasks, convertedTask)
	}
	// will obtain task num, and use append logic that concatenates old task less than 'num' : new tasks : old task greater than 'num'

	for i := range tasks[num:] {
		currentIndex := num + i
		adjustedTaskNumber := (tasks[currentIndex][0].(float64) - float64(num)) + float64(len(newTasks)) + 1
		tasks[currentIndex][0] = int(adjustedTaskNumber)
	}
	tasks = append(tasks[:(num-1)], append(convertedTasks, tasks[num:]...)...)
	updatedData, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal("Error marshalling data:", err)
	}
	if err := os.WriteFile("data.json", updatedData, 0644); err != nil {
		log.Fatal("Error writing new data to file:", err)
	}
	fmt.Printf("Replaced task #%d with three new smaller tasks! Good luck :)\n", num)
	buildTable()
}

