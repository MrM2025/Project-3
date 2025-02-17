package main

import (
	"github.com/MrM2025/rpforcalc/tree/master/calc_go/internal/application"
)

func main() {
	app := application.New()
	//app.Run() // Используется для проверки работы калькулятора без сервера: тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране, exit - останавливает приложение
	app.RunServer()
}
