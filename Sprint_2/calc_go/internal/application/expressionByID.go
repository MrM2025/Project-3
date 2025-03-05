package application

import (
	"encoding/json"
	"math"
	"net/http"
	"sync"
	//"strconv"
)

type IDForExpression struct {
	ID string `json:"id,omitempty"`
}

func ExpressionByID(w http.ResponseWriter, r *http.Request) {
	var (
		mu sync.Mutex
		//o *Orchestrator
	)

	mu.Lock()
	defer mu.Unlock()

	request := new(IDForExpression)
	json.NewDecoder(r.Body).Decode(&request)

	expr, ok := exprStore[request.ID]

	if !ok {
		http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
		return
	}

	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = math.Round(expr.AST.Value*100) / 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}
