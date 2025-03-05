# Распределённый вычислитель арифметических выражений
### Описание
Этот проект реализует веб-сервис, который вычисляет арифметические выражения, переданные пользователем через HTTP-запрос.



## Запуск 
#### Через git Bash
С помощью
``` bash
 git clone github.com/MrM2025/Project-3/tree/main/Sprint_2/calc_go
 ```
сделайте клон проекта.

Запустите Оркестратор:
#### Важно - проверьте, что вы находитесь в папке calc_go.

``` bash
export TIME_ADDITION_MS=100
export TIME_SUBTRACTION_MS=100
export TIME_MULTIPLICATIONS_MS=1000
export TIME_DIVISIONS_MS=1000

go run cmd/Orchestrator_start/main.go
```

#### В новом окне Git Bash
Запустите Агента:
#### Важно - проверьте, что вы находитесь в папке calc_go.

``` bash
export COMPUTING_POWER=4
export ORCHESTRATOR_URL=http://localhost:8080

 go run cmd/agent.start/main.go
``` 

# Для отправки curl используйте Postman

#### При обращении к http://localhost:8080 будет возвращен README-файл

Выражение для вычисления должно передаваться в JSON-формате, в единственном поле "expression", если поле отсутствует - сервер вернет ошибку 422, "Empty expression"; если в запросе будут поля, отличные от "expression" - сервер вернет ошибку 400, "Bad request" также как и при отсуствии JSON'а в теле запроса;

Должны быть установлены Go и Git.

## Пример запроса с использованием curl(Рекомендую использовать постман)
Для cmd windows:  

 curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"-1+1*2.54+41+((3/3+10)/2-(-2.5-1+(-1))*10)-1\" }" http://localhost:8080/api/v1/calculate (пример корректного запроса, код:200)

Для git bash:

curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{ "expression": "-1+1*2.54+41+((3/3+10)/2-(-2.5-1+(-1))*10)-1" }'
#

Postman:

https://identity.getpostman.com/signup?deviceId=c30fc039-7460-4f58-8cb9-b74256c4186c  

^

|

Регистрация

https://www.postman.com/downloads/

^

|

Ссылка на скачивание приложения.

#
Мануал №1 - https://timeweb.com/ru/community/articles/kak-polzovatsya-postman

Мануал №2 - https://blog.skillfactory.ru/glossary/postman/

Мануал №3 - https://gb.ru/blog/kak-testirovat-api-postman/

## Примеры использования (cmd Windows)

Верно заданный запрос, Status: 200

curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"20-(9+1)\"}" http://localhost:8080/api/v1/calculate

Запрос с пустым выражением, Status: 422, Error: empty expression

curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"\"}" http://localhost:8080/api/v1/calculate

Запрос с делением на 0, Status: 422, Error: division by zero

curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"1/0\"}" http://localhost:8080/api/v1/calculate

Запрос неверным выражением, Status : 422, Error: invalid expression

curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"1++*2\"}" http://localhost:8080/api/v1/calculate

## Тесты
Для тестирования перейдите в файл agent_calc_test.go и используйте команду go test или(для вывода дополнительной информации) go test -v

Для запусков всех тестов разом воспользуйтесь - go test ./...

