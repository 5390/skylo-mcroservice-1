### Microservices Project Readme

This project involves two microservices, microservice-1 and microservice-2, orchestrated using Docker Compose. The system also includes a PostgreSQL database and Kafka for message streaming.

Prerequisites

Ensure the following are installed on your system:

Docker

Docker Compose

cURL (for testing endpoints)

Project Structure

```
├── docker-compose.yml
├── microservice-1
│   ├── main.go
│   ├── Dockerfile
│   └── ...
├── microservice-2
│   ├── server.go
│   ├── Dockerfile
│   └── ...
├── db
│   ├── init.sql
├── kafka
│   └── ...
└── README.md
```

Services Overview

1. Microservice-1

Purpose: Sends data to microservice-2 and integrates with Kafka.

Key Endpoints:

POST /api/send: Sends a JSON payload to microservice-2.

Kafka Producer: Publishes messages to a Kafka topic.

Environment Variables:

MICROSERVICE_2_URL: URL for microservice-2.

KAFKA_BROKER: Address of the Kafka broker.

Dockerfile:

FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o microservice-1 .
CMD ["./microservice-1"]

2. Microservice-2

Purpose: Receives and processes data from microservice-1 and stores it in a PostgreSQL database.

Key Endpoints:

POST /api/data: Accepts JSON data and saves it to the database.

Environment Variables:

DB_HOST: PostgreSQL database host.

DB_USER: PostgreSQL user.

DB_PASSWORD: PostgreSQL password.

DB_NAME: PostgreSQL database name.

Dockerfile:

FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o microservice-2 .
CMD ["./microservice-2"]

3. PostgreSQL Database

Purpose: Stores data received by microservice-2.

Initial Setup: The init.sql file creates the necessary table.

CREATE TABLE IF NOT EXISTS received_messages (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Docker Configuration:

postgres:
  image: postgres:15
  container_name: postgres-db
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: mysecretpassword
    POSTGRES_DB: postgres
  volumes:
    - ./db/init.sql:/docker-entrypoint-initdb.d/init.sql
  ports:
    - "5432:5432"

4. Kafka

Purpose: Enables message streaming between services.

Setup in docker-compose.yml:

kafka:
  image: confluentinc/cp-kafka:latest
  container_name: kafka-broker
  environment:
    KAFKA_BROKER_ID: 1
    KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
  ports:
    - "9092:9092"
zookeeper:
  image: confluentinc/cp-zookeeper:latest
  container_name: zookeeper
  environment:
    ZOOKEEPER_CLIENT_PORT: 2181
  ports:
    - "2181:2181"

Running the Services

Build and Start Containers:

docker-compose up --build

Verify Services:

Check logs to ensure all containers are running.

Verify Kafka, PostgreSQL, and the two microservices are active.

Testing the Endpoints:

Microservice-2

Test data submission:

curl --location --request POST 'http://localhost:8081/api/data' \
--header 'Content-Type: application/json' \
--data-raw '{"data": "some data"}'

Microservice-1

Test forwarding requests to microservice-2:

curl --location --request POST 'http://localhost:8080/api/send' \
--header 'Content-Type: application/json' \
--data-raw '{"data": "message to service 2"}'

Verify Kafka Messages:
Use a Kafka client to confirm the messages are published to the specified topic.

Notes

Ensure proper network configuration in Docker Compose for service-to-service communication.

Use environment variables to adjust configurations based on deployment environments (e.g., development, staging, production).

Logs are available for debugging; ensure log levels are appropriately set for production.

Troubleshooting

Database Connection Errors: Ensure DB_HOST, DB_USER, and DB_PASSWORD are correctly set.

Kafka Issues: Verify KAFKA_BROKER and KAFKA_ZOOKEEPER_CONNECT settings.

CORS Issues: Ensure microservice-2 has the correct CORS settings if accessed from a different origin.

