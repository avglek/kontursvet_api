.PHONY: swag build run clean

# Переменная с вашей длинной командой swag
SWAG_CMD = swag init -g cmd/server/main.go -o ./docs --parseInternal --parseDependency

# Короткая команда для генерации сваггера
swag:
	$(SWAG_CMD)

# Команда для сборки проекта
build:
	go build -o server.exe ./cmd/server/main.go

# Команда для запуска проекта (сначала обновит сваггер, потом запустит)
run:
	go run ./cmd/server/main.go

# Очистка сгенерированных файлов
clean:
	rm -rf docs
	rm -f server.exe
