# gometa

Gometa: Генератор кода для создания крудов (все 3 слоя: контроллер, сервис, репа)

## Как использовать? 

Создаём файлик {name}.schema.json с примером объекта. См client.schema.json 

Запускаем `go run ./cmd/cli/ {name}.schema.json ../rzd-app-svc/`. `../rzd-app-svc/` поменять на путь к проекту

Сгенерировалась модель в models, сервис в services, репа repository и контроллер в api/v1/controllers. 

Остаётся только подключить контроллер в app/initializers и всё готово