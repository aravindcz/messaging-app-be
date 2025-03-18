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
| `id`        | UUID (PK)     | Primary Key         |
| `sender_id` | VARCHAR       | Not Null            |
| `receiver_id` | VARCHAR     | Not Null            |
| `content`   | TEXT          | Not Null            |
| `timestamp` | TIMESTAMP     | Default: Now()      |
| `read`      | BOOLEAN       | Default: False      |

### API Endpoints
#### 1️⃣ Send a Message (Asynchronous)
**POST** `/messages`
```json
{
  "sender_id": "user123",
  "receiver_id": "user456",
  "content": "Hello, how are you?"
}
```
_Response:_
```json
{
  "status": "Message queued"
}
```

#### 2️⃣ Retrieve Conversation History
**GET** `/messages?user1=user123&user2=user456`
_Response:_
```json
[
  {
    "message_id": "msg001",
    "sender_id": "user123",
    "receiver_id": "user456",
    "content": "Hey!",
    "timestamp": "2024-03-13T10:00:00Z",
    "read": true
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

## Setup Instructions

### 1️⃣ Run Services with Docker Compose
Ensure `docker-compose.yml` is correctly set up and then run:
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
List all tables:
```sql
\dt
```

### 4️⃣ Create RabbitMQ Queue
Inside the RabbitMQ container:
```sh
docker exec -it rabbitmq bash
rabbitmqadmin declare queue name=message_queue durable=true
```

### 5️⃣ Test API Using Postman
- Import the provided API request examples into Postman.
- Send a message, retrieve conversation history, and mark messages as read.

---
This document provides a comprehensive overview of the messaging API, ensuring a smooth setup and testing process. 🚀
