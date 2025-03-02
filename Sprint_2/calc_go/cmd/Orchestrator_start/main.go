package main

import (
	"github.com/MrM2025/rpforcalc/tree/master/calc_go/internal/application"
)

func main() {
	app := application.NewOrchestrator()
	app.RunOrchestrator()
}
