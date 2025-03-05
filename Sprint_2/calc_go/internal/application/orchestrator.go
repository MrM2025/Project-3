package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MrM2025/rpforcalc/tree/master/calc_go/pkg/errorStore"
)

type Config struct {
	Addr                string
	TimeAddition        int
	TimeSubtraction     int
	TimeMultiplications int
	TimeDivisions       int
}

func ConfigFromEnv() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ta, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if ta == 0 {
		ta = 100
	}
	ts, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if ts == 0 {
		ts = 100
	}
	tm, _ := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if tm == 0 {
		tm = 1000
	}
	td, _ := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if td == 0 {
		td = 1000
	}

	return &Config{
		Addr:                port,
		TimeAddition:        ta,
		TimeSubtraction:     ts,
		TimeMultiplications: tm,
		TimeDivisions:       td,
	}
}

type Orchestrator struct {
	Config      *Config
	taskStore   map[string]*Task
	taskQueue   []*Task
	mu          sync.Mutex
	exprCounter int
	taskCounter int
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		Config:    ConfigFromEnv(),
		taskStore: make(map[string]*Task),
		taskQueue: make([]*Task, 0),
	}
}

type OrchReqJSON struct {
	Expression string `json:"expression"`
}

type OrchResJSON struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type Expression struct {
	ID     string   `json:"id,omitempty"`
	Expr   string   `json:"expression,omitempty"`
	Status string   `json:"status,omitempty"`
	Result float64  `json:"result,omitempty"`
	AST    *ASTNode `json:"-"`
}

type Task struct {
	ID             string   `json:"id,omitempty"`
	ExprID         string   `json:"expression,omitempty"`
	Arg1           float64  `json:"arg1,omitempty"`
	Arg2           float64  `json:"arg2,omitempty"`
	Operation      string   `json:"operation,omitempty"`
	Operation_time int      `json:"operation_time,omitempty"`
	Node           *ASTNode `json:"-"`
}

var (
	exprStore = make(map[string]*Expression)
	calc      TCalc
)

func (o *Orchestrator) Tasks(expr *Expression) {
	var traverse func(node *ASTNode)
	traverse = func(node *ASTNode) {

		if node == nil || node.IsLeaf {
			return
		}

		traverse(node.Left)
		traverse(node.Right)
		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			if !node.TaskScheduled {
				o.taskCounter++
				taskID := strconv.Itoa(o.taskCounter)
				var opTime int
				switch node.Operator {
				case "+":
					opTime = o.Config.TimeAddition
				case "-":
					opTime = o.Config.TimeSubtraction
				case "*":
					opTime = o.Config.TimeMultiplications
				case "/":
					opTime = o.Config.TimeDivisions
				default:
					opTime = 100
				}

				task := &Task{
					ID:             taskID,
					ExprID:         expr.ID,
					Arg1:           node.Left.Value,
					Arg2:           node.Right.Value,
					Operation:      node.Operator,
					Operation_time: opTime,
					Node:           node,
				}
				node.TaskScheduled = true
				o.taskStore[taskID] = task
				o.taskQueue = append(o.taskQueue, task)
			}
		}
	}
	traverse(expr.AST)
}

var divbyzeroeerr error

func (o *Orchestrator) CalcHandler(w http.ResponseWriter, r *http.Request) { //Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
	var (
		emsg string
	)
	o.mu.Lock()
	defer o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	request := new(OrchReqJSON)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body) //Достаем выражение
	dec.DisallowUnknownFields()
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := calc.IsCorrectExpression(request.Expression) // Проверяем выражение на наличие ошибок

	if !ok && err != nil || divbyzeroeerr != nil { // Присваиваем ошибке статус-код, выводим их
		switch {
		case errors.Is(err, errorStore.EmptyExpressionErr):
			emsg = errorStore.EmptyExpressionErr.Error()

		case errors.Is(err, errorStore.IncorrectExpressionErr):
			emsg = errorStore.IncorrectExpressionErr.Error()

		case errors.Is(err, errorStore.NumToPopMErr): // numtopop > nums' slise length
			emsg = errorStore.NumToPopMErr.Error()

		case errors.Is(err, errorStore.NumToPopZeroErr): // numtopop <= 0
			emsg = errorStore.NumToPopZeroErr.Error()

		case errors.Is(err, errorStore.NthToPopErr): // no operator to pop
			emsg = errorStore.NthToPopErr.Error()

		case errors.Is(divbyzeroeerr, errorStore.DvsByZeroErr):
			emsg = errorStore.DvsByZeroErr.Error()
			divbyzeroeerr = nil
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(OrchResJSON{Error: emsg})
		return
	}

	o.exprCounter++
	exprID := strconv.Itoa(o.exprCounter)

	ast, err := ParseAST(request.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	expr := &Expression{
		ID:     exprID,
		Expr:   request.Expression,
		Status: "pending",
		AST:    ast,
	}

	exprStore[exprID] = expr
	o.Tasks(expr)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})

}

func (o *Orchestrator) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Mehod error, expected: "GET"`, http.StatusMethodNotAllowed)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	if len(o.taskQueue) == 0 {
		http.Error(w, `{"Error":"No task available"}`, http.StatusNotFound)
		return
	}

	task := o.taskQueue[0]
	o.taskQueue = o.taskQueue[1:]

	if expr, exists := exprStore[task.ExprID]; exists {
		expr.Status = "in_progress"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
	defer r.Body.Close()

}

func (o *Orchestrator) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	task, ok := o.taskStore[req.ID]

	if !ok {
		o.mu.Unlock()
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}

	task.Node.IsLeaf = true
	task.Node.Value = req.Result
	delete(o.taskStore, req.ID)

	if expr, exists := exprStore[task.ExprID]; exists {
		o.Tasks(expr)
		if expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = expr.AST.Value
		}
	}

	o.mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"result accepted"}`))
}

func makeAnAtomicExpr(Operation string, Arg1, Arg2 float64) (string, error) {
	arg1 := strconv.FormatFloat(Arg1, 'g', 8, 32)
	arg2 := strconv.FormatFloat(Arg2, 'g', 8, 32)

	var result string

	switch {
	case Operation == "+":
		result = arg1 + "+" + arg2
	case Operation == "-":
		result = arg1 + "-" + arg2
	case Operation == "*":
		result = arg1 + "*" + arg2
	case Operation == "/":
		if arg2 == "0" {
			return "0", errorStore.DvsByZeroErr //DvsByZeroErr
		}
		result = arg1 + "/" + arg2
	}
	return result, nil
}

func (o *Orchestrator) RunOrchestrator() {
	a := NewAgent()

	go func() {
		for i := 0; i < a.ComputingPower; i++ {
			a.worker()
		}
	}()

	mux := http.NewServeMux()
	http.Handle("/", mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //можно открыть README.md
		http.ServeFile(w, r, "..\\README.md")
	})
	mux.HandleFunc("/api/v1/calculate", o.CalcHandler)
	mux.HandleFunc("/api/v1/expressions", ExpressionsOutput)
	mux.HandleFunc("/api/v1/expression/id", ExpressionByID)
	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.GetTaskHandler(w, r)

		} else if r.Method == http.MethodPost {
			o.PostTaskHandler(w, r)
		}
	})

	go func() {
		for {
			time.Sleep(2 * time.Second)
			o.mu.Lock()
			if len(o.taskQueue) > 0 {
				log.Printf("Pending tasks in queue: %d", len(o.taskQueue))
			}
			o.mu.Unlock()
		}
	}()

	http.ListenAndServe(":"+o.Config.Addr, nil)

}
