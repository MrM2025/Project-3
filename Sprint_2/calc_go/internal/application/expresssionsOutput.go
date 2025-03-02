package application

import (
	"net/http"
	"encoding/json"
	"sync"
)

var EmptyExpression = Expression{
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
	
	if exprStore == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(exprStore)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exprStore)
}
