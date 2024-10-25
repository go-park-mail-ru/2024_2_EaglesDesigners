# Устанавливаем имя исполняемого файла

BINARY_NAME=backend_app


# Путь к исходному коду

SRC_DIR=src

CMD_DIR=cmd/app


# Устанавливаем команду для сборки

build:

	go build -o $(BINARY_NAME) $(CMD_DIR)/main.go


# Команда для очистки скомпилированных бинарников

clean:

	rm -f $(BINARY_NAME)


# Команда для запуска приложения

run: build

	./$(BINARY_NAME)


# Команда для генерации сваггера

swagger:

	swag init -g cmd/app/main.go


# Команда для отчета покрытия тестами

test:

	go test -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out


# Команда для установки зависимостей

deps:

	go mod tidy


# Основная команда по умолчанию

.PHONY: build clean run deps