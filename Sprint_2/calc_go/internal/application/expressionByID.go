package application

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type IDForExpression struct {
	ID string `json:"id,omitempty"`
}

func ExpressionByID(w http.ResponseWriter, r *http.Request) {
	var(
		mu sync.Mutex
		//o *Orchestrator
	)
	
	mu.Lock()
	w.Header().Set("Content-Type", "application/json")
	request := new(IDForExpression)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body) //Достаем выражение
	dec.DisallowUnknownFields()
	err := dec.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	ID, err := strconv.Atoi(request.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	sID := strconv.Itoa(ID) // ID in string

	if exprStore[sID] == EmptyExpression {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("There is no such expression")
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(exprStore[sID])
	}
	mu.Unlock()
}