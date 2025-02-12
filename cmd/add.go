package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long:  "Add a new task to the todo list",
	Run: func(cmd *cobra.Command, args []string) {
		id := 100
		description := "Did this work?"
		createdAt := getCurrentTimestamp()

		task := Task{id, description, createdAt, false}

		err := writeToCSV("datastore.csv", task)
		if err != nil {
			fmt.Println("Error writing task: ", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
