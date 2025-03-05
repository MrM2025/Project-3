package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MrM2025/rpforcalc/tree/master/calc_go/pkg/errorStore"
)

type AgentTask struct {
	ID             string  `json:"id,omitempty"`
	ExprID         string  `json:"expression,omitempty"`
	Arg1           float64 `json:"arg1,omitempty"`
	Arg2           float64 `json:"arg2,omitempty"`
	Operation      string  `json:"operation,omitempty"`
	Operation_time int     `json:"operation_time,omitempty"`
	Result         string  `json:"result,omitempty"`
}

type AgentResJSON struct {
	ID     string  `json:"ID,omitempty"`
	Result float64 `json:"result,omitempty"`
}

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
	client          *http.Client
}

func NewAgent() *Agent {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cp < 1 {
		cp = 1
	}

	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")

	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}
	return &Agent{
		ComputingPower:  cp,
		OrchestratorURL: orchestratorURL,
	}
}

func (a *Agent) worker() {
	for {
		resp, err := http.Get(a.OrchestratorURL + "/internal/task")
		if err != nil {
			log.Printf("Worker %d: error getting task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}
		var taskResp struct {
			Task struct {
				ID            string  `json:"id"`
				Arg1          float64 `json:"arg1"`
				Arg2          float64 `json:"arg2"`
				Operation     string  `json:"operation"`
				OperationTime int     `json:"operation_time"`
			} `json:"task"`
		}
		err = json.NewDecoder(resp.Body).Decode(&taskResp)
		resp.Body.Close()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		task := taskResp.Task
		log.Printf("Worker: received task %s: %f %s %f, simulating %d ms", task.ID, task.Arg1, task.Operation, task.Arg2, task.OperationTime)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		divbyzeroeerr = nil
		result, diverr := calculator(task.Operation, task.Arg1, task.Arg2)

		if errors.Is(diverr, errorStore.DvsByZeroErr) {
			divbyzeroeerr = errorStore.DvsByZeroErr
		}

		resultPayload := map[string]interface{}{
			"id":     task.ID,
			"result": result,
		}

		payloadBytes, _ := json.Marshal(resultPayload)
		respPost, err := http.Post(a.OrchestratorURL+"/internal/task", "application/json", bytes.NewReader(payloadBytes))
		respPost.Body.Close()

	}
}

func calculator(operator string, arg1, arg2 float64) (float64, error) {
	var result float64

	switch {
	case operator == "+":
		result = arg1 + arg2
	case operator == "-":
		result = arg1 - arg2
	case operator == "*":
		result = arg1 * arg2
	case operator == "/":
		if arg2 == 0 {
			return 0, errorStore.DvsByZeroErr
		}
		result = arg1 / arg2
	default:
		return 0, fmt.Errorf("Error")
	}

	return result, nil
}

func (a *Agent) SendResult(request *AgentTask, result []byte) {

	_, err := strconv.Atoi(request.ID)
	if err != nil {
		log.Printf("Error of type conversion %v", err)
	}

	req, err := http.NewRequest("POST", a.OrchestratorURL+"/internal/task", bytes.NewReader(result))
	if err != nil {
		log.Printf("Error fetching task: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)

	}

	_, err = a.client.Do(req)
	if err != nil {
		log.Printf("Error doing request: %v. Retrying in 2 seconds...", err)

	}
	/*
		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			log.Printf("Worker : error response posting result for task %v: %s", taskStore[ID-1], string(body))
		} else {
			log.Printf("Worker : successfully completed task %v with result %s", taskStore[ID-1], result)
		}
	*/

	defer req.Body.Close()

}

func (a *Agent) RunAgent() {
	for i := 0; i < a.ComputingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker()
	}
	select {}
}
