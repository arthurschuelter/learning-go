package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Habit struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Entries     []HabitEntry `json:"entries"`
	Streak      int          `json:"streak"`
}

type HabitEntry struct {
	EntryDate string `json:"entry_date"`
	HabitId   string `json:"habit_id"`
}

type EntryRequest struct {
	HabitId   string `json:"habit_id"`
	EntryDate string `json:"entry_date"`
}

func main() {
	habits := GetHabits()
	for {

		var command string
		PrintMenu()

		_, err := fmt.Scan(&command)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		if command == "1" {
			fmt.Println("Fetching habits...")
			fmt.Println("")
			habits := GetHabits()
			printFetchHabits(habits)
		}

		if command == "2" {
			var habitChoice int
			clear()
			fmt.Println("=======================")
			PrintHabitsOptions(habits)
			fmt.Print("Enter command: ")
			_, err := fmt.Scan(&habitChoice)
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			// fmt.Printf("You chose habit: %s\n", habits[habitChoice-1].Name)
			currentTime := time.Now()
			var requestBody = EntryRequest{
				HabitId:   habits[habitChoice-1].Id,
				EntryDate: currentTime.Format("2006-01-02"),
			}
			PostHabitEntry(requestBody)
		}

		if command == "3" {
			fmt.Println("Bye bye")
			return
		}

	}
}
func PrintMenu() {
	fmt.Println("=======================")
	fmt.Println("(1) Fetch habits")
	fmt.Println("(2) Add habit entry")
	fmt.Println("(3) Exit")
	fmt.Print("Enter command: ")
}

func GetHabits() []Habit {
	base_url := "http://localhost:8000/api/habits"
	response, err := http.Get(base_url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var habits []Habit
	err = json.Unmarshal(body, &habits)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	return habits
}

func PostHabitEntry(entryRequest EntryRequest) {
	base_url := "http://localhost:8000/api/habit_entries"
	jsonData, err := json.Marshal(entryRequest)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	response, err := http.Post(base_url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()
	fmt.Println("Habit successfully logged!")
}

func PrintHabitsOptions(habits []Habit) {
	for i, habit := range habits {
		fmt.Printf("(%d) %s\n", i+1, habit.Name)
	}
}

func printFetchHabits(habits []Habit) {
	clear()
	for _, habit := range habits {
		fmt.Printf("%s (%s)\n", habit.Name, habit.Id)
		fmt.Printf("Description: %s\n", habit.Description)
		fmt.Printf("Streak: %d\n", habit.Streak)
		fmt.Println("Entries:")
		for _, entry := range habit.Entries {
			fmt.Printf("  %s\n", entry.EntryDate)
		}
		fmt.Println("")
	}
}
func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
