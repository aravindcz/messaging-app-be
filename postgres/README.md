Once the intial infrastructure has been setup based on the docker compose we just need to create a messages table that can be created either using ui or from the infra automation scripts itself using 

"
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id VARCHAR(255) NOT NULL,
    receiver_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read BOOLEAN DEFAULT FALSE
);
"

