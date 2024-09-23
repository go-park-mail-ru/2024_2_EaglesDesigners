# Устанавливаем имя исполняемого файла

BINARY_NAME=app


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


# Команда для установки зависимостей

deps:

	go mod tidy


# Основная команда по умолчанию

.PHONY: build clean run deps