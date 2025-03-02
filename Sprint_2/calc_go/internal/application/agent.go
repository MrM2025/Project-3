package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MrM2025/rpforcalc/tree/master/calc_go/pkg/errorStore"
)

type AgentResJSON struct {
	ID     string  `json:"ID,omitempty"`
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type Agent struct {
	ComputingPower  int
	OrchestratorURL string
	client          *http.Client
}

func NewAgent() *Agent {

	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}

	return &Agent{
		ComputingPower:  1,
		OrchestratorURL: orchestratorURL,
		client: &http.Client{ // Инициализация клиента
			Timeout: 10 * time.Second,
		},
	}
}

func (a *Agent) worker() {
	if a.client == nil {
		log.Fatal("HTTP client is not initialized")
	}
	for {
		req, err := http.NewRequest("GET", a.OrchestratorURL+"/internal/task", nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		res, err := a.client.Do(req)
		if err != nil {
			log.Printf("Error fetching task: %v. Retrying in 2 seconds...", err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer res.Body.Close()  

		if res.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code: %d. Retrying...", res.StatusCode)
			time.Sleep(2 * time.Second)
			continue
		}

		request := new(Task)
		dec := json.NewDecoder(res.Body) //Достаем подвыражение из res
		dec.DisallowUnknownFields()
		derr := dec.Decode(&request)
		if derr != nil {
			log.Printf("Error fetching task: %v. Retrying in 2 seconds...", derr)
			time.Sleep(2 * time.Second)
			continue
		}

		calcresult, cerr := calculator(request.Operation, request.Arg1, request.Arg2, request.Operation_time) //Производим вычисления
		if cerr != nil {
			log.Printf("Error fetching task: %v. Retrying in 2 seconds...", cerr)
			time.Sleep(2 * time.Second)
			continue
		}

		result := AgentResJSON{
			ID:     request.ID,
			Result: calcresult,
		}

		body, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshaling result: %v.", err)
			return
		}

		req, err = http.NewRequest("POST", a.OrchestratorURL+"/internal/task", bytes.NewReader(body))
		req.Body.Close()

		time.Sleep(1 * time.Second)
	}
}

func calculator(operator string, arg1, arg2 float64, operation_time int) (float64, error) {
	time.Sleep(time.Duration(operation_time) * time.Millisecond)
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
	}

	return result, nil
}

func (a *Agent) RunAgent() {

	computingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	for i := 0; i < computingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker()
	}
}
