# Calculator API Service

Этот проект был создан как финальная задача для первого спринта в Яндекс Лицее, он представляет из себя веб-сервис для вычисления арифметических выражений. Сервис принимает арифметическое выражение через HTTP-запрос и возвращает результат вычисления.

## Установка и запуск

### Предварительные требования

- Установленный [Go](https://golang.org/dl/) версии 1.20 или выше.
- Git (для клонирования репозитория).

### Клонирование репозитория

Сначала клонируйте репозиторий:
```bash
git clone https://github.com/username/calc_service.git
cd calc_service
```

### Установка зависимостей

Убедитесь, что все зависимости установлены:

```bash
go mod tidy
```

### Запуск сервиса

Запустите сервис с помощью следующей команды:
```bash
go run ./cmd/calc_service/...
```
Сервер будет запущен на порту `8080`.

## Примеры использования

### Успешный запрос

Для выполнения арифметического выражения отправьте POST-запрос на `/api/v1/calculate` с телом запроса в формате JSON:

#### Powershell
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/v1/calculate" `
     -Method POST `
     -ContentType "application/json" `
     -Body '{ "expression": "2+2*2" }'
```

#### Windows CMD
```cmd
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{ \"expression\": \"2+2*2\" }"
```

#### Bash
```bash
curl --location "http://localhost:8080/api/v1/calculate" \
--header "Content-Type: application/json" \
--data '{
  "expression": "2+2*2"
}'
```

**Ожидаемый ответ:**
```json
{
"result": 6
}
```

### Ошибки

В случае возникновения ошибок сервис вернет соответствующий код состояния и сообщение об ошибке. Ниже приведена таблица с примерами ошибок:

| Код ошибки | Описание ошибки             | Пример запроса                                 | Ожидаемый ответ                              |
|------------|------------------------------|------------------------------------------------|----------------------------------------------|
| 422        | Неверное выражение          | `{"expression": "2+2*2a"}`                    | `{"error": "Expression is not valid"}`     |
| 400        | Деление на ноль             | `{"expression": "10/0"}`                       | `{"error": "Division by zero"}`             |
| 422        | Несоответствие скобок       | `{"expression": "(2+3"`                        | `{"error": "Expression is not valid"}`     |
| 422        | Два оператора подряд        | `{"expression": "2++2"}`                       | `{"error": "Expression is not valid"}`     |
