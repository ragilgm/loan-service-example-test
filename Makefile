PROJ=$(shell pwd)

test:
	go test -v ./...

test-coverage:
	go test -failfast -coverprofile=test_result/coverage.out -covermode=count ./...  && go tool cover -html=test_result/coverage.out -o test_result/coverage.html

run:
	docker-compose -f ./deploy/pg.yaml up --build -d
	docker-compose -f ./deploy/kafka.yaml up --build -d

stop:
	docker-compose -f ./deploy/pg.yaml down
	docker-compose -f ./deploy/kafka.yaml down


migrate-up:
	migrate -database "postgresql://dbuser:dbpass@:5432/dbname?sslmode=disable" -path ./database/pg/migration up

migrate-down:
	migrate -database "postgresql://dbuser:dbpass@:5432/dbname?sslmode=disable" -path ./database/pg/migration down

generate:
	@go install github.com/golang/mock/mockgen@v1.6.0
	@PROJ=${PROJ} go generate ./...
