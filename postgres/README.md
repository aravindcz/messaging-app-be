Once the intial infrastructure has been setup based on the docker compose we just need to create a messages table that can be created either using ui or from the infra automation scripts itself using 

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