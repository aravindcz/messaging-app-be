# Messaging API Architecture

## Objective
The objective of this project is to build an API that enables users to:
1. Send text messages asynchronously using a message queue (RabbitMQ).
2. Retrieve conversation history between two users.
3. Mark messages as read.

## High-Level Design

### Components
- **API Service (Gin Framework in Go)**: Handles incoming HTTP requests.
- **Message Queue (RabbitMQ)**: Manages asynchronous message processing.
- **Worker Service (Go)**: Consumes messages from RabbitMQ and stores them in PostgreSQL.
- **Database (PostgreSQL)**: Stores message history.
- **Docker Compose**: Manages multi-container deployment.

### Architecture Diagram
```
[Client] --> [API Service] --> [RabbitMQ] --> [Worker Service] --> [PostgreSQL]
                         |--> [Retrieve from PostgreSQL]
```

## Low-Level Design

### Database Schema (PostgreSQL)
#### `messages` Table
| Column       | Type          | Constraints          |
|-------------|--------------|----------------------|
| `id`        | SERIAL INTEGER (PK)     | Primary Key         |
| `sender_id` | INTEGER       | Not Null            |
| `receiver_id` | INTEGER     | Not Null            |
| `content`   | TEXT          | Not Null            |
| `timestamp` | TIMESTAMP     | Default: Now()      |
| `read`      | BOOLEAN       | Default: False      |

### API Endpoints
#### 1️⃣ Send a Message (Asynchronous)
**POST** `/messages`
```json
{
  "sender_id": 100,
  "receiver_id": 250,
  "content": "Hello, how are you?"
}
```
_Response:_
```json
{
  "message": "Message sent to queue"
}
```

#### 2️⃣ Retrieve Conversation History
**GET** `/messages?user1=100&user2=250`

_Response:_
```json
[
  {
    "id": 1,
    "sender_id": 100,
    "receiver_id": 250,
    "content": "Hello, how are you?",
    "timestamp": "2024-03-13T10:00:00Z",
    "read": false
  }
]
```

#### 3️⃣ Mark a Message as Read
**PATCH** `/messages/{message_id}/read` 

_Response:_
```json
{
  "status": "read"
}
```

## Prerequisites

### Docker and Docker Compose
```sh
sudo apt update
sudo apt install docker.io -y
sudo apt install docker-compose -y
```

## Setup Instructions

### 1️⃣ Run Services with Docker Compose
Clone the repository and ensure `docker-compose.yml` is correctly set up and then run:
```sh
docker-compose up --build
```

### 2️⃣ Verify Running Containers
Check if all services are up:
```sh
docker ps
```

### 3️⃣ Log into PostgreSQL Container
```sh
docker exec -it postgres psql -U postgres -d messaging
```
Create message tables:
```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read BOOLEAN DEFAULT FALSE
);
```
```sql

CREATE INDEX idx_sender_timestamp ON messages (sender_id, timestamp); -- Index on sender_id and timestamp

CREATE INDEX idx_receiver_timestamp ON messages (receiver_id, timestamp); -- Index on receiver_id and timestamp
```

### 4️⃣ Create RabbitMQ Queue
Inside the RabbitMQ container:
```sh
docker exec -it rabbitmq bash
rabbitmqadmin -u admin -p adminpassword declare queue name=message_queue durable=true
```

### 5️⃣ Test API Using Postman
- Import the provided API request examples into Postman.
- Send a message, retrieve conversation history, and mark messages as read.

Sent Messages Curl Request
------------------------------
curl --location 'http://ip-address:8080/messages' \
--header 'Content-Type: application/json' \
--data '{
           "sender_id": 100,
           "reciever_id": 250,
           "content": "Hello, how are you?"
         }'

Get Conversation History Curl Request
----------------------------------------
curl --location 'http://ip-address:8080/messages?user1=100&user2=250' \
--header 'Content-Type: application/json'

Change Message Status Curl Request
-----------------------------------
curl --location --request PATCH 'http://ip-address:8080/messages/1/read' \
--header 'Content-Type: application/json'

---
This document provides a comprehensive overview of the messaging API, ensuring a smooth setup and testing process. 🚀
