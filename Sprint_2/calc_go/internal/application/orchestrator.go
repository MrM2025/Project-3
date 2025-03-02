package application

import (
	"encoding/json"
	"errors"
	//"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	//"time"

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

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		//taskQueue: make([]Task, 0),
		Config: ConfigFromEnv(),
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
	ID     string  `json:"id,omitempty"`
	Expr   string  `json:"expression,omitempty"`
	Status string  `json:"status,omitempty"`
	Result float64 `json:"result,omitempty"`
}

type Task struct {
	ID             string  `json:"id,omitempty"`
	ExprID         string  `json:"expression,omitempty"`
	Arg1           float64 `json:"arg1,omitempty"`
	Arg2           float64 `json:"arg2,omitempty"`
	Operation      string  `json:"operation,omitempty"`
	Operation_time int     `json:"operation_time,omitempty"`
}

type Orchestrator struct {
	mu            sync.Mutex
	//exprStore     map[string]Expression
	exprID        int
	//taskStore     []map[string]Task
	currentTaskID int
	//taskQueue     []Task
	Config *Config
}

var (
	calc TCalc
	exprStore = make(map[string]Expression)
	taskStore = make([]map[string]Task, 0)
)

func (o *Orchestrator) CalcHandler(w http.ResponseWriter, r *http.Request) { //Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
	var (
		emsg string
		expr Expression
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

	if !ok && err != nil { // Присваиваем ошибке статус-код, выводим их
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

		case errors.Is(err, errorStore.DvsByZeroErr):
			emsg = errorStore.DvsByZeroErr.Error()
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(OrchResJSON{Error: emsg})
		return
	}

	o.exprID++
	ID := strconv.Itoa(o.exprID)
	expr = Expression{
		ID:     strconv.Itoa(o.exprID),
		Expr:   request.Expression,
		Status: "pending",
	}

	tasks, err := calc.ExprtolightExprs(request.Expression, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskStore = tasks
	//log.Println(request.Expression, o.taskStore, tasks)

	exprStore[ID] = expr
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ID)

}

func (o *Orchestrator) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `Mehod error, expected: "GET"`, http.StatusMethodNotAllowed)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	if len(taskStore) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	if o.currentTaskID >= len(taskStore) {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	task := taskStore[o.currentTaskID]
	o.currentTaskID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func (o *Orchestrator) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Mehod error, expected: "POST"`, http.StatusMethodNotAllowed)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	var result struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	expr, exists := exprStore[result.ID]
	if !exists {
	    http.Error(w, "Expression not found", http.StatusNotFound)
	    return
	}
	
	expr.Result = result.Result
	expr.Status = "done"
	exprStore[result.ID] = expr
	
	w.WriteHeader(http.StatusOK)}

func (o *Orchestrator) RunOrchestrator() error {
	a := NewAgent() // Инициализация агента
	if a == nil {
		return errors.New("failed to initialize agent")
	}

	mux := http.NewServeMux()
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
		} else {
			http.Error(w, `Wrong method, expected: "GET" or "POST"`, http.StatusMethodNotAllowed)
			return
		}
	})
	http.Handle("/", mux)
	return http.ListenAndServe(":"+o.Config.Addr, nil)

}
