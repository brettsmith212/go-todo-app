package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"

	"github.com/spf13/cobra"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   string
	IsCompleted bool
}

func readTasksFromCSV(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	var tasks []Task
	for _, record := range records[1:] { // Skip header
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error converting ID to int: %w", err)
		}
		task := Task{
			ID:          id,
			Description: record[1],
			CreatedAt:   record[2],
			IsCompleted: record[3] == "true",
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func formatTimeAgo(timeStr string) string {
	createdTime, err := time.Parse("2006-01-02T15:04:05-07:00", timeStr)
	if err != nil {
		return timeStr
	}
	diff := timediff.TimeDiff(createdTime)
	return diff
}

func printTasks(tasks []Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)

	fmt.Fprintln(w, "ID\tDescription\tCreated At\tStatus")
	for _, task := range tasks {
		status := "Pending"
		if task.IsCompleted {
			status = "Completed"
		}
		timeAgo := formatTimeAgo(task.CreatedAt)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", task.ID, task.Description, timeAgo, status)
	}
	w.Flush()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks in the todo list`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := readTasksFromCSV("datastore.csv")
		if err != nil {
			fmt.Println("Error reading tasks:", err)
			return
		}
		printTasks(tasks)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
