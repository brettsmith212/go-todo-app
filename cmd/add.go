package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05-07:00")
}

func writeToCSV(filename string, task Task) error {
	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{fmt.Sprint(task.ID), task.Description, task.CreatedAt, fmt.Sprint(task.IsCompleted)}
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("error writing record: %w", err)
	}

	fmt.Printf("Successfully added task: %+v\n", task)
	return nil
}

func readMostRecentTaskIdFromCSV(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return -1, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) == 0 {
		return -1, fmt.Errorf("no records found")
	}

	record := records[len(records)-1]

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return -1, fmt.Errorf("error converting stoi: %w", err)
	}

	return id, nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long:  "Add a new task to the todo list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Must pass 1 arg -> Task Description")
			return
		}

		id, err := readMostRecentTaskIdFromCSV(Database)
		if err != nil {
			fmt.Println("Error getting most recent task id: ", err)
			return
		}

		id += 1
		description := args[0]
		createdAt := getCurrentTimestamp()

		task := Task{id, description, createdAt, false}

		err = writeToCSV(Database, task)
		if err != nil {
			fmt.Println("Error writing task: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
