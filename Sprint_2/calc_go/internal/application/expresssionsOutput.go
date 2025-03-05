package application

import (
	"net/http"
	"encoding/json"
	"sync"
)

var EmptyExpression = &Expression{
	Status: "",
}

func ExpressionsOutput(w http.ResponseWriter, r *http.Request) { //Сервер, который выводит все переданные серверу выражения
	var(
		mu sync.Mutex
		//o *Orchestrator
	)
	mu.Lock()
	defer mu.Unlock()


	w.Header().Set("Content-Type", "application/json")
	
	exprs := make([]*Expression, 0, len(exprStore))

	for _, expr := range exprStore {
		if expr.AST != nil && expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = expr.AST.Value
		}
		exprs = append(exprs, expr)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}
