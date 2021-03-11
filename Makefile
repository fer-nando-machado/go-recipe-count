PROJECT_NAME = recipe-count
MODULE_NAME = cmd
DB_NAME = data
ARGS = -file=$(file) -postcode=$(postcode) -time=$(time) -recipes=$(recipes)

.SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help:
	$(info recipe-count commands:)
	$(info . setup                      installs dependencies)
	$(info . format                     formats go files)
	$(info . build                      compiles binary)
	$(info . test                       runs available tests)
	$(info . run                        starts application, accepts the following args:)
	$(info .    file=data/demo.json     fixtures data file path (required))
	$(info .    postcode=99999          postcode to search for)
	$(info .    time=12AM-12PM          delivery time to search for)
	$(info .    recipes=apple,cake      recipe(s) name(s) to search for, separated by commas)
	$(info . docker-build               builds application @ docker)
	$(info . docker-test                runs available tests @ docker)
	$(info . docker-run                 starts application @ docker (accepts the same args from 'run'))

.PHONY: setup
setup:
	go get -d -v -t ./...
	go install -v ./...
	go mod tidy -v

.PHONY: format
format:
	go fmt ./...

.PHONY: build
build:
	go build -v -o $(MODULE_NAME) ./$(MODULE_NAME)

.PHONY: test
test:
	go test ./... -v -covermode=count

.PHONY: run
run:
	go run ./$(MODULE_NAME) $(ARGS)

.PHONY: docker-build
docker-build:
	docker build --build-arg root_dir=./$(MODULE_NAME) --build-arg db_dir=./$(DB_NAME) -t $(PROJECT_NAME) .

.PHONY: docker-run
docker-run: docker-build
	docker run -a stdout -a stderr -t --name $(PROJECT_NAME) --rm $(PROJECT_NAME) go run . $(ARGS)

.PHONY: docker-test
docker-test: docker-build
	docker run --name $(PROJECT_NAME) --rm $(PROJECT_NAME) go test ./... -v -covermode=count
