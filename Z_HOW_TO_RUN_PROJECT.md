# Step-by-Step Guide to Run the Application Using Docker and Docker Compose

This guide provides the steps to set up and run your application with Docker and Docker Compose.

### Prerequisites

Before starting, ensure you have the following tools installed:
1. **Docker** - [Download here](https://www.docker.com/get-started).
2. **Docker Compose** - [Installation guide here](https://docs.docker.com/compose/install/).
3. **Go** - [Installation guide here](https://go.dev/doc/install).

### A Steps to Run the Application (DOCKER)

#### Step 0: Download all dependency 
First, ensure all Go dependencies are installed. Run the following command:
```bash
go mod tidy
```

#### Step 1: Create a Docker Network

First, you need to create a Docker network for your services. This network will enable communication between the PostgreSQL database, Kafka, and your loan-service.

Run the following command:

```bash
docker network create test-network
```

#### Step 2: Set Up the PostgreSQL Database
Next, start the PostgreSQL container using Docker Compose. This step will initialize the database container in the background.
Run the following command:
```bash
docker-compose -f ./deploy/pg.yaml up --build -d
```

#### Step 3: Set Up Kafka and Zookeeper
Now, you need to start the Kafka and Zookeeper containers, which are essential for handling messages in the application.

Run the following command:
```bash
docker-compose -f ./deploy/kafka.yaml up --build -d
```


#### Step 4: Build the Application
Once the services are set up, you need to build the application container with Docker Compose.
Run the following command to build the loan-service container:

```bash
docker-compose -f ./docker-compose.yaml up --build -d
```

#### Step 6: Run Database Migrations

Finally, apply the necessary database migrations to ensure the schema is up to date with the application.
Run the following command:
```bash
migrate -database "postgresql://dbuser:dbpass@:5432/dbname?sslmode=disable" -path ./database/pg/migration up
```


### B Steps to Run the Application (LOCAL)


#### Step 0: Download all dependency
First, ensure all Go dependencies are installed. Run the following command:
```bash
go mod tidy
```

#### Step 1: Setup Configuration File .env
Ensure the .env file is configured correctly with your database, Kafka, and email credentials. Example .env file:

```bash
APP_ADDRESS=:9090
APP_DEBUG=true
APP_READ_TIMEOUT=5s
APP_WRITE_TIMEOUT=10s

#db
PG_CONN_MAX_LIFETIME=30m
PG_DBNAME=dbname
PG_DBPASS=dbpass
PG_DBUSER=dbuser
PG_HOST=pg
PG_MAX_IDLE_CONNS=6
PG_MAX_OPEN_CONNS=30
PG_PORT=5432

#kafka
KAFKA_BROKER_ADDRESS=kafka:9092
KAFKA_TOPIC=loan-topic
KAFKA_CONSUMER_GROUP=loan-consumer-group
PRODUCER_RETRIES=3
KAFKA_TIMEOUT=30s


#smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=xxxx
SMTP_PASSWORD=xxxx
```

#### Step 2: Build the Application
To build the Go application, run the following command:
```bash
go build -o loan-service ./cmd
```

#### Step 3: Run the Application
Run the application locally by executing the following command:
```bash
./loan-service
```


#### Step 4: Set Up Database
If you need to run migrations on your local database, run the following command:
```bash
migrate -database "postgresql://dbuser:dbpass@localhost:5432/dbname?sslmode=disable" -path ./database/pg/migration up
```



