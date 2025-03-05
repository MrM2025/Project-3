# Распределённый вычислитель арифметических выражений

### Дисклеймер
К сожалению, я слишком поздно сообразил, что без абстрактного синтаксического дереваб разбиение на атомарные операции плохореализуемо, поэтому был вынужден взять чужой код построения дерева отсюда: https://github.com/Killered672/Module2calc/blob/main/internal/orchestrator/astnode.go и встроить его  у себя (

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
go run cmd/Orchestrator_start/main.go
```

### Агент запускать не нужно(он запускается автоматически). 

# Для отправки curl используйте Postman

Выражение для вычисления должно передаваться в JSON-формате, в единственном поле "expression", если поле отсутствует - сервер вернет ошибку 422, "Empty expression"; если в запросе будут поля, отличные от "expression" - сервер вернет ошибку 400, "Bad request" также как и при отсуствии JSON'а в теле запроса;

Должны быть установлены Go и Git.

## Пример запроса с использованием curl(Рекомендую использовать Postman)



Для git bash:
``` bash
curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{ "expression": "-1+1*2.54+41+((3/3+10)/2-(-2.5-1+(-1))*10)-1" }'
```
(пример корректного запроса, код:200)

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

### * Важно: при отображении readme в HTLM, экранирующие слэши не отображаются, поэтому копировать команды лучше из raw-формата, либо самостоятельно экранировать ковычки в json'е слэшом слева, иначе получите ошибку!

Верно заданный запрос, Status: 200
```
curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"20-(9+1)\"}" http://localhost:8080/api/v1/calculate
```
Ответ:
{
    "expressions": [
        {
            "id": "1",
            "expression": "20-(9+1)",
            "status": "completed",
            "result": 10
        }
    ]
}


Запрос с пустым выражением, Status: 422, Error: empty expression
```
curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"\"}" http://localhost:8080/api/v1/calculate
```
Ответ:
{
    "error": "empty expression"
}

Запрос с делением на 0, Status: 422, Error: division by zero
```
curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"1/0\"}" http://localhost:8080/api/v1/calculate
```
Запрос неверным выражением, Status : 422, Error: invalid expression
```
curl -i -X POST -H "Content-Type:application/json" -d "{\"expression\": \"1++*2\"}" http://localhost:8080/api/v1/calculate
```
## Тесты
Для тестирования перейдите в файл agent_calc_test.go и используйте команду go test или(для вывода дополнительной информации) go test -v

Для запусков всех тестов разом воспользуйтесь - go test ./...

