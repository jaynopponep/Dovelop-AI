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
}

func sendToGPT(prompt string) {
	var token = ""
	client := openai.NewClient(token)
	var content = "The following is a task that you will break down into an ARRAY of THREE tasks only. Do not respond with anything else except the ARRAY. Here's the task: " + prompt
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
	var resTasks []string
	err = json.Unmarshal([]byte(res), &resTasks)
	if err != nil {
		log.Fatalf("Error parsing response %v\n", err)
	}
	fmt.Println("Tasks:")
	for i, task := range resTasks {
		fmt.Printf("Task %d: %s\n", i+1, task)
	}
}
