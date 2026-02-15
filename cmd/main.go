package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	App()
}

func App() {
	if err := GetAllTasks(); err != nil {
		log.Fatal(err)
	}
}

type Due struct {
	Date        string  `json:"date"`
	Timezone    string `json:"timezone"`
	String      string  `json:"string"`
	Lang        string  `json:"lang"`
	IsRecurring bool    `json:"is_recurring"`
}

type Task struct {
	TaskID    string  `json:"id"`
	ProjectID string  `json:"project_id"`
	ParentID  string `json:"parent_id"`
	Content   string  `json:"content"`
	Priority  int     `json:"priority"`
	Checked   bool    `json:"checked"`

	Due *Due `json:"due"`
}

type TasksResponse struct {
	Results []Task `json:"results"`
}

func GetAllTasks() error {
	token := os.Getenv("TODOIST_TOKEN")
	log.Println("Token upload: ", token)
	// curl https://api.todoist.com/api/v1/tasks \                                              at  13:32:35
	// -H "Authorization: Bearer b3411fd99b9e90a17fab2f29246dbbc94dfb9958"

	url := "https://api.todoist.com/api/v1/tasks"
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new request failed: %w", err)
	}
	startReq := time.Now()

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf(
		"method=%s url=%s status=%d duration=%v",
		req.Method,
		req.URL.String(),
		resp.StatusCode,
		time.Since(startReq),
	)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("todoist bad status %d", resp.StatusCode)
	}

	var tasks TasksResponse

	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}

	for _, t := range tasks.Results {
		if t.Due == nil {
			continue
		}
		b, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println(string(b))
	}
	return nil
}
