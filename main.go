package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alexeyco/simpletable"
	"log"
	"os"
	"strings"
)

var (
	sampleData = [][]interface{}{
		{1, "Finish FAFSA Application"},
		{2, "Reply to [person]'s email"},
		{3, "Review code written by [person]"},
		{4, "Close issue #420"},
		{5, "Walk Dog"},
	}
)

type Task struct {
	Num  int
	Task string
}

func processChoice(choice int) {
	switch choice {
	case 1:
		fmt.Println("You chose to create a new task")
		var num int = 5
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
		t := Task{Num: num, Task: task}
		createTask(t)
	case 2:
		fmt.Println("You chose option 2")
	case 3:
		fmt.Println("You chose option 3")
	}
}

func buildTable() {
	// beginning of building table
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Task #"},
			{Align: simpletable.AlignCenter, Text: "Task"},
		},
	}

	// mapping sampleData below, will use a reliable database instead!
	for _, row := range sampleData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", row[0].(int))}, // <- Task Number of type (int)
			{Text: row[1].(string)}, // <- Task of type (string)
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

func createTask(t Task) {
	// Marshal retrieved data into bytes
	tBytes, err := json.Marshal(t)
	log.Print(string(tBytes))
	if err != nil {
		log.Println("Error marshalling JSON item", err)
		log.Print(err)
		return
	}

	// check if file already exists; if not, automatically create. otherwise, append because already exists
	file, err := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file) // schedule a file close

	input, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// create new line if data already not empty
	if input.Size() > 0 {
		if _, err := file.WriteString(",\n"); err != nil {
			log.Fatal(err)
		}
	}

	// simply add new task
	_, err = file.Write(tBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data added")
}

/*
func getData() [][]interface{} {
	file, err := os.Open("data.json")
	if err != nil { log.Fatal(err) }
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err!= nil { log.Fatal(err) }


}
*/

func main() {
	buildTable()
	// begin main menu process
	fmt.Println("Choose an option")
	var choice int
	_, inputErr := fmt.Scanln(&choice)
	if inputErr != nil {
		fmt.Println("Invalid input, try again")
		return
	} else {
		processChoice(choice)
	}
}
