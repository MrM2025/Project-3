package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	//"io"

	"github.com/MrM2025/rpforcalc/tree/master/calc_go/pkg/errorStore"
)

type AgentTask struct {
	ID             string  `json:"id,omitempty"`
	ExprID         string  `json:"expression,omitempty"`
	Arg1           float64 `json:"arg1,omitempty"`
	Arg2           float64 `json:"arg2,omitempty"`
	Operation      string  `json:"operation,omitempty"`
	Operation_time int     `json:"operation_time,omitempty"`
	Result         float64 `json:"result,omitempty"`
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
		log.Println(3)
		req, err := http.NewRequest("GET", a.OrchestratorURL+"/internal/task", nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		res, err := a.client.Do(req)
		if err != nil {
			log.Printf("Error doing request: %v. Retrying in 2 seconds...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println(4)

		request := new(AgentTask)
		dec := json.NewDecoder(res.Body) //Достаем подвыражение из res
		dec.DisallowUnknownFields()
		derr := dec.Decode(&request)
		if derr != nil {
			log.Printf("Error decoding task: %v.", derr)
			time.Sleep(3 * time.Second)
			continue
		}
		res.Body.Close()

		//log.Println("a", request.Arg1, request.Arg2)

		calcresult, cerr := calculator(request.Operation, request.Arg1, request.Arg2, request.Operation_time) //Производим вычисления
		if cerr != nil {
			log.Printf("Calculator error: %v.", cerr)
			return
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
		log.Println(4)

		a.SendResult(request, body)

		time.Sleep(3 * time.Second)
	}
}

func calculator(operator string, arg1, arg2 float64, operation_time int) (float64, error) {
	time.Sleep(time.Duration(operation_time) * time.Millisecond)
	var result float64

	log.Println("c", operator, arg1, arg2)

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

	computingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	for i := 0; i < computingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker()
	}
}
