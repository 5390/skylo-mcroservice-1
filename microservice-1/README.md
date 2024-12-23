### Microservices Project Readme

This project involves two microservices, microservice-1 and microservice-2, orchestrated using Docker Compose. The system also includes a PostgreSQL database and Kafka for message streaming.

Prerequisites

Ensure the following are installed on your system:

Docker

Docker Compose

cURL (for testing endpoints)

Project Structure

```
project/
├── microservice-1/
│   ├── main.go
│   ├── queue/         # Queue consumer logic
│   ├── retry/         # Retry mechanism
│   ├── db/            # Database interactions
│   └── config/        # Configurations and environment variables
├── microservice-2/
│   ├── main.go
│   ├── server/           # REST Server implementation
│   ├── db/            # Database interactions
│   └── config/        # Configurations and environment variables
├── README.md          # Setup, build, and run instructions
└── docker-compose.yml # Optional: For containerized setup
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

```FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o microservice-1 .
CMD ["./microservice-1"]
```

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

```FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o microservice-2 .
CMD ["./microservice-2"]
```

3. PostgreSQL Database

Purpose: Stores data received by microservice-2.

Initial Setup: The init.sql file creates the necessary table.

CREATE TABLE IF NOT EXISTS received_messages (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

Docker Configuration:
```
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
```

4. Kafka

Purpose: Enables message streaming between services.

Setup in docker-compose.yml:
```
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
```

Running the Services

Build and Start Containers:

```docker-compose up --build```

**Verify Services**:

Using Kafka CLI (Command-Line Interface)
Prerequisite:
```docker ps```

Kafka CLI tools must be available. These are included in the Kafka Docker container.
Steps:

Access the Kafka Container: Run the following command to get a shell inside the Kafka Docker container:

```docker exec -it kafka bash```

Produce a Message: Use the kafka-console-producer command to send a message to the topic my-topic:

```kafka-console-producer --broker-list localhost:9092 --topic my-topic```

Enter the Message: After running the above command, the terminal will wait for input. Type a message and press Enter:

```{"data": "test-message"}```

Each line you type will be sent as a message to the topic.
Verify the Message: Use the kafka-console-consumer to verify the message:

```kafka-console-consumer --bootstrap-server localhost:9092 --topic my-topic --from-beginning```

Check logs to ensure all containers are running.

Verify ```Kafka```, ```PostgreSQL```, and the two microservices are active.

Testing the Endpoints:

Microservice-2

Test data submission:
```
curl --location --request POST 'http://localhost:8081/api/data' \
--header 'Content-Type: application/json' \
--data-raw '{"data": "some data"}'
```

**Verify Kafka Messages**:
Use a Kafka client to confirm the messages are published to the specified topic.
Connect Kafka  : 
```docker exec -it kafka kafka-console-producer --broker-list kafka:9092 --topic my-topic```
Add Message to Queue:

**Detailed Flow**
**Message Production:**

Messages are added by Commands (or another producer system) and added to the Kafka topic ```my-topic```.

**Microservice-1 Consumes Messages**:

**Microservice-1** connects to Kafka and continuously polls messages from the topic my-topic.
Microservice-1 Sends Messages to **Microservice-2**:

Microservice-1 reads the message payload and sends it as a ```POST``` request to Microservice-2's API ```(/api/data)```.
If the API call fails (e.g., due to a network error or Microservice-2 being down), Microservice-1 retries the operation every 10 seconds until the message is successfully delivered.
Microservice-2 Processes and Stores Messages:

Microservice-2 receives the message from Microservice-1 through its POST API.
The message is validated and then inserted into the PostgreSQL table received_messages.

**Notes**

Ensure proper network configuration in Docker Compose for service-to-service communication.

Use environment variables to adjust configurations based on deployment environments (e.g., development, staging, production).

Logs are available for debugging; ensure log levels are appropriately set for production.

Troubleshooting

**Database Connection Errors**: Ensure ```DB_HOST, DB_USER, and DB_PASSWORD ``` are correctly set.

**Kafka Issues**: ```Verify KAFKA_BROKER``` and ```KAFKA_ZOOKEEPER_CONNECT``` settings.

**CORS Issues**: Ensure microservice-2 has the correct CORS settings if accessed from a different origin.

