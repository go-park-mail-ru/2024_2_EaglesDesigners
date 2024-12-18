# Устанавливаем имя исполняемого файла
BINARY_NAME=backend_app

# Путь к исходному коду
SRC_DIR=src

CMD_DIR=cmd/app

up:
	sudo docker-compose down && sudo docker-compose up --build

down:
	sudo docker-compose down

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
	swag init --parseDependency --parseInternal -g main_app/cmd/app/main.go

# Команда для запуска тестов
test:
	go test ./...

# Команда для отчета покрытия тестами
cover:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
	cat coverprofile_.tmp | grep -v mocks.go | grep -v .pb.go | grep -v _grpc.go | grep -v _mock.go | grep -v main.go | grep -v docs.go > coverprofile.tmp ; \
	rm coverprofile_.tmp ; \
	go tool cover -html coverprofile.tmp ; \
	go tool cover -func coverprofile.tmp

# Команда для установки зависимостей
deps:
	go mod tidy

# Команда запуска линтера
lint:
	golangci-lint run

# Команда автофикса линтера
lintfix:
	golangci-lint run --fix

# Основная команда по умолчанию
.PHONY: build lint clean run deps