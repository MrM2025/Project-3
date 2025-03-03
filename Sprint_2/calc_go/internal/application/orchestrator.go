package application

import (
	"encoding/json"
	"errors"
	"log"

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

type abc struct {
	Atomic Task // Arg1 = A.Result, arg2 = B.result, Operator = Operator
	A Task // Atomic
	B Task // Atomic и так до простейсшего действия
	Operator string
}

type Orchestrator struct {
	mu sync.Mutex
	exprID int
	currentTaskID int
	Config        *Config
}

var (
	calc      TCalc
	exprStore = make(map[string]Expression)
	taskStore = make([]Task, 0)
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
		ID:     ID,
		Expr:   request.Expression,
		Status: "pending",
	}
//
//
// 1
	tasks, err := calc.ExprtolightExprs(request.Expression, ID, "None")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskStore = tasks

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

	expr := exprStore[task.ExprID]
	expr.Status = "processing"
	exprStore[task.ExprID] = expr

	o.currentTaskID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
	defer r.Body.Close()

}

func (o *Orchestrator) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `Mehod error, expected: "POST"`, http.StatusMethodNotAllowed)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		ID     string  `json:"ID,omitempty"`
		Result float64 `json:"result,omitempty"`
	}
	var result Result

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	expr, exists := exprStore[result.ID]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	
	ID, err := strconv.Atoi(result.ID)
	if err != nil {
		log.Printf("Error of type conversion %v", err)
	}

	expression := exprStore[taskStore[ID].ExprID].Expr
	arg1 := taskStore[ID].Arg1
	arg2 := taskStore[ID].Arg2
	op := taskStore[ID].Operation

	atomicExpr, err := makeAnAtomicExpr(op, arg1, arg2)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	taskStore, err = calc.ExprtolightExprs(expression, taskStore[ID].ExprID, atomicExpr)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	expr.Result = result.Result
	expr.Status = "done"
	exprStore[taskStore[ID].ExprID] = expr

	log.Println(exprStore)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(result)
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

func (o *Orchestrator) RunOrchestrator() error {
	a := NewAgent() // Инициализация агента
	if a == nil {
		return errors.New("failed to initialize agent")
	}

	computingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if computingPower == 0 {
		computingPower = 1
	}
	for i := 0; i < computingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker()
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
